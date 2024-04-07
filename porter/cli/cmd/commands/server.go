package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/porter-dev/porter/cli/cmd/config"
	"github.com/porter-dev/porter/cli/cmd/docker"
	"github.com/porter-dev/porter/cli/cmd/github"
	"github.com/porter-dev/porter/cli/cmd/utils"

	"github.com/spf13/cobra"
)

type startOps struct {
	imageTag string `form:"required"`
	db       string `form:"oneof=sqlite postgres"`
	driver   string `form:"required"`
	port     *int   `form:"required"`
}

var opts = &startOps{}

func registerCommand_Server(cliConf config.CLIConfig) *cobra.Command {
	serverCmd := &cobra.Command{
		Use:     "server",
		Aliases: []string{"svr"},
		Short:   "Commands to control a local Porter server",
	}

	// startCmd represents the start command
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a Porter server instance on the host",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			if cliConf.Driver == "docker" {
				_ = cliConf.SetDriver("docker")

				err := startDocker(
					ctx,
					cliConf,
					opts.imageTag,
					opts.db,
					*opts.port,
				)
				if err != nil {
					red := color.New(color.FgRed)
					_, _ = red.Println("Error running start:", err.Error())
					_, _ = red.Println("Shutting down...")

					err = stopDocker(ctx)

					if err != nil {
						_, _ = red.Println("Shutdown unsuccessful:", err.Error())
					}

					os.Exit(1)
				}
			} else {
				_ = cliConf.SetDriver("local")
				err := startLocal(
					ctx,
					cliConf,
					opts.db,
					*opts.port,
				)
				if err != nil {
					red := color.New(color.FgRed)
					_, _ = red.Println("Error running start:", err.Error())
					os.Exit(1)
				}
			}
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops a Porter instance running on the Docker engine",
		Run: func(cmd *cobra.Command, args []string) {
			if cliConf.Driver == "docker" {
				if err := stopDocker(cmd.Context()); err != nil {
					_, _ = color.New(color.FgRed).Println("Shutdown unsuccessful:", err.Error())
					os.Exit(1)
				}
			}
		},
	}

	serverCmd.AddCommand(startCmd)
	serverCmd.AddCommand(stopCmd)

	serverCmd.PersistentFlags().AddFlagSet(utils.DriverFlagSet)

	startCmd.PersistentFlags().StringVar(
		&opts.db,
		"db",
		"sqlite",
		"the db to use, one of sqlite or postgres",
	)

	startCmd.PersistentFlags().StringVar(
		&opts.imageTag,
		"image-tag",
		"latest",
		"the Porter image tag to use (if using docker driver)",
	)

	opts.port = startCmd.PersistentFlags().IntP(
		"port",
		"p",
		8080,
		"the host port to run the server on",
	)
	return serverCmd
}

func startDocker(
	ctx context.Context,
	cliConf config.CLIConfig,
	imageTag string,
	db string,
	port int,
) error {
	env := []string{
		"NODE_ENV=production",
		"FULLSTORY_ORG_ID=VXNSS",
	}

	var porterDB docker.PorterDB

	switch db {
	case "postgres":
		porterDB = docker.Postgres
	case "sqlite":
		porterDB = docker.SQLite
	}

	startOpts := &docker.PorterStartOpts{
		ProcessID:      "main",
		ServerImageTag: imageTag,
		ServerPort:     port,
		DB:             porterDB,
		Env:            env,
	}

	_, _, err := docker.StartPorter(ctx, startOpts)
	if err != nil {
		return err
	}

	green := color.New(color.FgGreen)

	green.Printf("Server ready: listening on localhost:%d\n", port)

	return cliConf.SetHost(fmt.Sprintf("http://localhost:%d", port))
}

func startLocal(
	ctx context.Context,
	cliConf config.CLIConfig,
	db string,
	port int,
) error {
	if db == "postgres" {
		return fmt.Errorf("postgres not available for local driver, run \"porter server start --db postgres --driver docker\"")
	}

	cliConf.SetHost(fmt.Sprintf("http://localhost:%d", port))

	porterDir := filepath.Join(home, ".porter")
	cmdPath := filepath.Join(home, ".porter", "portersvr")
	sqlLitePath := filepath.Join(home, ".porter", "porter.db")
	staticFilePath := filepath.Join(home, ".porter", "static")

	if _, err := os.Stat(cmdPath); os.IsNotExist(err) {
		err := downloadMatchingRelease(ctx, porterDir)
		if err != nil {
			color.New(color.FgRed).Println("Failed to download server binary:", err.Error())
			os.Exit(1)
		}
	}

	// otherwise, check the version flag of the binary
	cmdVersionPorter := exec.Command(cmdPath, "--version")
	writer := &config.VersionWriter{}
	cmdVersionPorter.Stdout = writer

	err := cmdVersionPorter.Run()

	if err != nil || writer.Version != config.Version {
		err := downloadMatchingRelease(ctx, porterDir)
		if err != nil {
			color.New(color.FgRed).Println("Failed to download server binary:", err.Error())
			os.Exit(1)
		}
	}

	cmdPorter := exec.Command(cmdPath)
	cmdPorter.Env = os.Environ()
	cmdPorter.Env = append(cmdPorter.Env, []string{
		"IS_LOCAL=true",
		"SQL_LITE=true",
		"SQL_LITE_PATH=" + sqlLitePath,
		"STATIC_FILE_PATH=" + staticFilePath,
		fmt.Sprintf("SERVER_PORT=%d", port),
		"REDIS_ENABLED=false",
	}...)

	if _, found := os.LookupEnv("GITHUB_ENABLED"); !found {
		cmdPorter.Env = append(cmdPorter.Env, "GITHUB_ENABLED=false")
	}

	if _, found := os.LookupEnv("PROVISIONER_ENABLED"); !found {
		cmdPorter.Env = append(cmdPorter.Env, "PROVISIONER_ENABLED=false")
	}

	cmdPorter.Stdout = os.Stdout
	cmdPorter.Stderr = os.Stderr

	err = cmdPorter.Run()

	if err != nil {
		color.New(color.FgRed).Println("Failed:", err.Error())
		os.Exit(1)
	}

	return nil
}

func stopDocker(ctx context.Context) error {
	agent, err := docker.NewAgentFromEnv(ctx)
	if err != nil {
		return err
	}

	err = agent.StopPorterContainersWithProcessID(ctx, "main", false)

	if err != nil {
		return err
	}

	green := color.New(color.FgGreen)

	green.Println("Successfully stopped the Porter server.")

	return nil
}

func downloadMatchingRelease(ctx context.Context, porterDir string) error {
	z := &github.ZIPReleaseGetter{
		AssetName:           "portersvr",
		AssetFolderDest:     porterDir,
		ZipFolderDest:       porterDir,
		ZipName:             "portersvr_latest.zip",
		EntityID:            "porter-dev",
		RepoName:            "porter",
		IsPlatformDependent: true,
		Downloader: &github.ZIPDownloader{
			ZipFolderDest:   porterDir,
			AssetFolderDest: porterDir,
			ZipName:         "portersvr_latest.zip",
		},
	}

	err := z.GetRelease(ctx, config.Version)
	if err != nil {
		return err
	}

	zStatic := &github.ZIPReleaseGetter{
		AssetName:           "static",
		AssetFolderDest:     filepath.Join(porterDir, "static"),
		ZipFolderDest:       porterDir,
		ZipName:             "static_latest.zip",
		EntityID:            "porter-dev",
		RepoName:            "porter",
		IsPlatformDependent: false,
		Downloader: &github.ZIPDownloader{
			ZipFolderDest:   porterDir,
			AssetFolderDest: filepath.Join(porterDir, "static"),
			ZipName:         "static_latest.zip",
		},
	}

	return zStatic.GetRelease(ctx, config.Version)
}
