//go:build ee

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joeshaw/envdecode"
	"github.com/porter-dev/porter/api/server/shared/config/env"
	"github.com/porter-dev/porter/internal/adapter"
	"github.com/porter-dev/porter/internal/opa"
	"github.com/porter-dev/porter/internal/repository"
	"github.com/porter-dev/porter/internal/worker"
	"github.com/porter-dev/porter/workers/jobs"
	"gorm.io/gorm"

	"github.com/porter-dev/porter/ee/integrations/vault"
	rcreds "github.com/porter-dev/porter/internal/repository/credentials"
	pgorm "github.com/porter-dev/porter/internal/repository/gorm"
)

var (
	jobQueue    chan worker.Job
	envDecoder  = EnvConf{}
	dbConn      *gorm.DB
	repo        repository.Repository
	opaPolicies *opa.KubernetesPolicies
)

// EnvConf holds the environment variables for this binary
type EnvConf struct {
	// ServerURL is the URL of the Porter server
	ServerURL string `env:"SERVER_URL,default=http://localhost:8080"`

	// Porter instance's database configuration
	DBConf env.DBConf

	// DigitalOcean OAuth2 credentials
	DOClientID     string `env:"DO_CLIENT_ID"`
	DOClientSecret string `env:"DO_CLIENT_SECRET"`

	// Worker pool configuration
	MaxWorkers uint `env:"MAX_WORKERS,default=10"`
	MaxQueue   uint `env:"MAX_QUEUE,default=100"`
	Port       uint `env:"PORT,default=3000"`

	/**
	 * Job-specific configuration
	 */

	// "helm-revisions-count-tracker"
	AWSAccessKeyID     string `env:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
	AWSRegion          string `env:"AWS_REGION"`
	S3BucketName       string `env:"S3_BUCKET_NAME"`
	EncryptionKey      string `env:"S3_ENCRYPTION_KEY"`
	RevisionsCount     int    `env:"REVISIONS_COUNT,default=20"`

	// "recommender"
	OPAConfigFileDir string `env:"OPA_CONFIG_FILE_DIR,default=./internal/opa"`
	LegacyProjectIDs []uint `env:"LEGACY_PROJECT_IDS"`

	// "preview-deployments-ttl-deleter"
	PreviewDeploymentsTTL string `env:"PREVIEW_DEPLOYMENTS_TTL"`
}

func main() {
	ctx := context.Background()

	if err := envdecode.StrictDecode(&envDecoder); err != nil {
		log.Fatalf("Failed to decode server conf: %v", err)
	}

	log.Printf("setting max worker count to: %d\n", envDecoder.MaxWorkers)
	log.Printf("setting max job queue count to: %d\n", envDecoder.MaxQueue)

	log.Printf("legacy project ids are: %v", envDecoder.LegacyProjectIDs)

	db, err := adapter.New(&envDecoder.DBConf)
	if err != nil {
		log.Fatalln(err)
	}

	dbConn = db

	var credBackend rcreds.CredentialStorage

	if envDecoder.DBConf.VaultAPIKey != "" && envDecoder.DBConf.VaultServerURL != "" && envDecoder.DBConf.VaultPrefix != "" {
		credBackend = vault.NewClient(
			envDecoder.DBConf.VaultServerURL,
			envDecoder.DBConf.VaultAPIKey,
			envDecoder.DBConf.VaultPrefix,
		)
	}

	var key [32]byte

	for i, b := range []byte(envDecoder.DBConf.EncryptionKey) {
		key[i] = b
	}

	repo = pgorm.NewRepository(db, &key, credBackend)

	opaPolicies, err = opa.LoadPolicies(envDecoder.OPAConfigFileDir)

	if err != nil {
		log.Fatalln(err)
	}

	jobQueue = make(chan worker.Job, envDecoder.MaxQueue)
	d := worker.NewDispatcher(int(envDecoder.MaxWorkers))

	log.Println("starting worker dispatcher")

	err = d.Run(ctx, jobQueue)

	if err != nil {
		log.Fatalln(err)
	}

	server := &http.Server{Addr: fmt.Sprintf(":%d", envDecoder.Port), Handler: httpService(ctx)}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		log.Println("shutting down server")

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCtxCancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err = server.Shutdown(shutdownCtx)

		if err != nil {
			log.Fatalln(err)
		}

		log.Println("server shutdown completed")

		serverStopCtx()
	}()

	log.Println("starting HTTP server at :3000")

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("error starting HTTP server: %v", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	d.Exit()
}

func httpService(ctx context.Context) http.Handler {
	log.Println("setting up HTTP router and adding middleware")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	// r.Use(middleware.AllowContentType("application/json"))

	r.Mount("/debug", middleware.Profiler())

	log.Println("setting up HTTP POST endpoint to enqueue jobs")

	r.Post("/enqueue/{id}", func(w http.ResponseWriter, r *http.Request) {
		req := make(map[string]interface{})

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("error converting body to json: %v", err)
			return
		}

		job := getJob(ctx, chi.URLParam(r, "id"), req)

		if job == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		jobQueue <- job
		w.WriteHeader(http.StatusCreated)
	})

	return r
}

func getJob(ctx context.Context, id string, input map[string]interface{}) worker.Job {
	if id == "helm-revisions-count-tracker" {
		newJob, err := jobs.NewHelmRevisionsCountTracker(ctx, dbConn, time.Now().UTC(), &jobs.HelmRevisionsCountTrackerOpts{
			DBConf:             &envDecoder.DBConf,
			DOClientID:         envDecoder.DOClientID,
			DOClientSecret:     envDecoder.DOClientSecret,
			DOScopes:           []string{"read", "write"},
			ServerURL:          envDecoder.ServerURL,
			AWSAccessKeyID:     envDecoder.AWSAccessKeyID,
			AWSSecretAccessKey: envDecoder.AWSSecretAccessKey,
			AWSRegion:          envDecoder.AWSRegion,
			S3BucketName:       envDecoder.S3BucketName,
			EncryptionKey:      envDecoder.EncryptionKey,
			RevisionsCount:     envDecoder.RevisionsCount,
		})
		if err != nil {
			log.Printf("error creating job with ID: helm-revisions-count-tracker. Error: %v", err)
			return nil
		}

		return newJob
	} else if id == "recommender" {
		newJob, err := jobs.NewRecommender(dbConn, time.Now().UTC(), &jobs.RecommenderOpts{
			DBConf:           &envDecoder.DBConf,
			DOClientID:       envDecoder.DOClientID,
			DOClientSecret:   envDecoder.DOClientSecret,
			DOScopes:         []string{"read", "write"},
			ServerURL:        envDecoder.ServerURL,
			Input:            input,
			LegacyProjectIDs: envDecoder.LegacyProjectIDs,
		}, opaPolicies)
		if err != nil {
			log.Printf("error creating job with ID: recommender. Error: %v", err)
			return nil
		}

		return newJob
	} else if id == "preview-deployments-ttl-deleter" {
		newJob, err := jobs.NewPreviewDeploymentsTTLDeleter(dbConn, time.Now().UTC(), &jobs.PreviewDeploymentsTTLDeleterOpts{
			DBConf:                &envDecoder.DBConf,
			ServerURL:             envDecoder.ServerURL,
			DOClientID:            envDecoder.DOClientID,
			DOClientSecret:        envDecoder.DOClientSecret,
			DOScopes:              []string{"read", "write"},
			PreviewDeploymentsTTL: envDecoder.PreviewDeploymentsTTL,
		})
		if err != nil {
			log.Printf("error creating job with ID: preview-deployments-ttl-deleter. Error: %v", err)
			return nil
		}

		return newJob
	}

	return nil
}
