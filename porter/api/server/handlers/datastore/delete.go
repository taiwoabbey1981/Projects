package datastore

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	porterv1 "github.com/porter-dev/api-contracts/generated/go/porter/v1"
	"github.com/porter-dev/porter/api/server/authz"
	"github.com/porter-dev/porter/api/server/handlers"
	"github.com/porter-dev/porter/api/server/handlers/release"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/apierrors"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/server/shared/requestutils"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/models"
	"github.com/porter-dev/porter/internal/telemetry"
)

// DeleteDatastoreHandler is a struct for handling datastore deletion requests
type DeleteDatastoreHandler struct {
	handlers.PorterHandlerReadWriter
	authz.KubernetesAgentGetter
}

// NewDeleteDatastoreHandler constructs a datastore DeleteDatastoreHandler
func NewDeleteDatastoreHandler(
	config *config.Config,
	decoderValidator shared.RequestDecoderValidator,
	writer shared.ResultWriter,
) *DeleteDatastoreHandler {
	return &DeleteDatastoreHandler{
		PorterHandlerReadWriter: handlers.NewDefaultPorterHandler(config, decoderValidator, writer),
		KubernetesAgentGetter:   authz.NewOutOfClusterAgentGetter(config),
	}
}

func (h *DeleteDatastoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.NewSpan(r.Context(), "serve-delete-datastore")
	defer span.End()
	project, _ := ctx.Value(types.ProjectScope).(*models.Project)

	datastoreName, reqErr := requestutils.GetURLParamString(r, types.URLParamDatastoreName)
	if reqErr != nil {
		err := telemetry.Error(ctx, span, nil, "error parsing datastore name")
		h.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
		return
	}
	telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "datastore-name", Value: datastoreName})

	datastoreRecord, err := h.Repo().Datastore().GetByProjectIDAndName(ctx, project.ID, datastoreName)
	if err != nil {
		err = telemetry.Error(ctx, span, err, "datastore record not found")
		h.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	if datastoreRecord == nil || datastoreRecord.ID == uuid.Nil {
		err = telemetry.Error(ctx, span, nil, "datastore record does not exist")
		h.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusNotFound))
		return
	}

	_, err = h.Repo().Datastore().UpdateStatus(ctx, datastoreRecord, models.DatastoreStatus_AwaitingDeletion)
	if err != nil {
		err = telemetry.Error(ctx, span, err, "error updating datastore status")
		h.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	updateReq := connect.NewRequest(&porterv1.UpdateDatastoreRequest{
		ProjectId:   int64(project.ID),
		DatastoreId: datastoreRecord.ID.String(),
	})

	_, err = h.Config().ClusterControlPlaneClient.UpdateDatastore(ctx, updateReq)
	if err != nil {
		err := telemetry.Error(ctx, span, err, "error calling ccp update datastore")
		h.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	// if the release was deleted by helm without error, mark it as accepted
	w.WriteHeader(http.StatusAccepted)
}

// UninstallDatastoreInput is the input type for UninstallDatastore
type UninstallDatastoreInput struct {
	ProjectID                         uint
	Name                              string
	CloudProvider                     string
	CloudProviderCredentialIdentifier string
	Request                           *http.Request
}

// UninstallDatastore uninstalls a datastore from a cluster
func (h *DeleteDatastoreHandler) UninstallDatastore(ctx context.Context, inp UninstallDatastoreInput) error {
	ctx, span := telemetry.NewSpan(ctx, "uninstall-datastore")
	defer span.End()

	telemetry.WithAttributes(span,
		telemetry.AttributeKV{Key: "project-id", Value: inp.ProjectID},
		telemetry.AttributeKV{Key: "name", Value: inp.Name},
		telemetry.AttributeKV{Key: "cloud-provider", Value: inp.CloudProvider},
		telemetry.AttributeKV{Key: "cloud-provider-credential-identifier", Value: inp.CloudProviderCredentialIdentifier},
	)

	var datastoreCluster *models.Cluster
	clusters, err := h.Repo().Cluster().ListClustersByProjectID(inp.ProjectID)
	if err != nil {
		return telemetry.Error(ctx, span, err, "unable to get project clusters")
	}

	for _, cluster := range clusters {
		if cluster.CloudProvider == inp.CloudProvider && cluster.CloudProviderCredentialIdentifier == inp.CloudProviderCredentialIdentifier {
			datastoreCluster = cluster
		}
	}

	if datastoreCluster == nil {
		return telemetry.Error(ctx, span, nil, "unable to find datastore cluster")
	}

	telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "cluster-id", Value: datastoreCluster.ID})

	helmAgent, err := h.GetHelmAgent(ctx, inp.Request, datastoreCluster, release.Namespace_ACKSystem)
	if err != nil {
		return telemetry.Error(ctx, span, err, "unable to get helm client for cluster")
	}

	_, err = helmAgent.UninstallChart(ctx, inp.Name)
	if err != nil {
		return telemetry.Error(ctx, span, err, "unable to uninstall chart")
	}

	return nil
}
