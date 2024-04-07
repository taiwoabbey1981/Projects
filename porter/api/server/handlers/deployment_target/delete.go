package deployment_target

import (
	"net/http"

	"connectrpc.com/connect"
	porterv1 "github.com/porter-dev/api-contracts/generated/go/porter/v1"
	"github.com/porter-dev/porter/api/server/authz"
	"github.com/porter-dev/porter/api/server/handlers"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/apierrors"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/server/shared/requestutils"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/models"
	"github.com/porter-dev/porter/internal/telemetry"
)

// DeleteDeploymentTargetHandler is the handler for DELETE /api/projects/{project_id}/clusters/{cluster_id}/deployment-targets/{deployment_target_id}
type DeleteDeploymentTargetHandler struct {
	handlers.PorterHandlerReadWriter
	authz.KubernetesAgentGetter
}

// NewDeleteDeploymentTargetHandler creates a new DeleteDeploymentTargetHandler
func NewDeleteDeploymentTargetHandler(
	config *config.Config,
	decoderValidator shared.RequestDecoderValidator,
	writer shared.ResultWriter,
) *DeleteDeploymentTargetHandler {
	return &DeleteDeploymentTargetHandler{
		PorterHandlerReadWriter: handlers.NewDefaultPorterHandler(config, decoderValidator, writer),
		KubernetesAgentGetter:   authz.NewOutOfClusterAgentGetter(config),
	}
}

// ServeHTTP deletes the deployment target from the cluster
func (c *DeleteDeploymentTargetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.NewSpan(r.Context(), "server-delete-deployment-target-by-id")
	defer span.End()

	project, _ := ctx.Value(types.ProjectScope).(*models.Project)

	deploymentTargetID, reqErr := requestutils.GetURLParamString(r, types.URLParamDeploymentTargetID)
	if reqErr != nil {
		err := telemetry.Error(ctx, span, reqErr, "error parsing deployment target id")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
		return
	}
	if deploymentTargetID == "" {
		err := telemetry.Error(ctx, span, nil, "deployment target id cannot be empty")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
		return
	}

	deleteReq := connect.NewRequest(&porterv1.DeleteDeploymentTargetRequest{
		ProjectId:          int64(project.ID),
		DeploymentTargetId: deploymentTargetID,
	})

	_, err := c.Config().ClusterControlPlaneClient.DeleteDeploymentTarget(ctx, deleteReq)
	if err != nil {
		err := telemetry.Error(ctx, span, err, "error deleting deployment target")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
		return
	}

	c.WriteResult(w, r, nil)
}
