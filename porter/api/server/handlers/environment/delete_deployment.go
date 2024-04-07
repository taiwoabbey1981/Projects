package environment

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v41/github"
	"github.com/porter-dev/porter/api/server/authz"
	"github.com/porter-dev/porter/api/server/handlers"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/apierrors"
	"github.com/porter-dev/porter/api/server/shared/commonutils"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/server/shared/requestutils"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/models"
	"gorm.io/gorm"
)

type DeleteDeploymentHandler struct {
	handlers.PorterHandlerReadWriter
	authz.KubernetesAgentGetter
}

func NewDeleteDeploymentHandler(
	config *config.Config,
	decoderValidator shared.RequestDecoderValidator,
	writer shared.ResultWriter,
) *DeleteDeploymentHandler {
	return &DeleteDeploymentHandler{
		PorterHandlerReadWriter: handlers.NewDefaultPorterHandler(config, decoderValidator, writer),
		KubernetesAgentGetter:   authz.NewOutOfClusterAgentGetter(config),
	}
}

func (c *DeleteDeploymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	project, _ := r.Context().Value(types.ProjectScope).(*models.Project)
	cluster, _ := r.Context().Value(types.ClusterScope).(*models.Cluster)

	deplID, reqErr := requestutils.GetURLParamUint(r, "deployment_id")

	if reqErr != nil {
		c.HandleAPIError(w, r, apierrors.NewErrInternal(reqErr))
		return
	}

	// read the deployment
	depl, err := c.Repo().Environment().ReadDeploymentByID(project.ID, cluster.ID, deplID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.HandleAPIError(w, r, apierrors.NewErrNotFound(errDeploymentNotFound))
			return
		}

		c.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// check that the environment belongs to the project and cluster IDs
	env, err := c.Repo().Environment().ReadEnvironmentByID(project.ID, cluster.ID, depl.EnvironmentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.HandleAPIError(w, r, apierrors.NewErrNotFound(errEnvironmentNotFound))
			return
		}

		c.HandleAPIError(w, r, apierrors.NewErrInternal(reqErr))
		return
	}

	// try to cancel any existing github workflow for this deployment
	client, err := getGithubClientFromEnvironment(c.Config(), env)
	if err == nil {
		workflowRun, err := commonutils.GetLatestWorkflowRun(client, depl.RepoOwner, depl.RepoName,
			fmt.Sprintf("porter_%s_env.yml", env.Name), depl.PRBranchFrom)
		if err == nil {
			if workflowRun.GetStatus() == "in_progress" || workflowRun.GetStatus() == "queued" ||
				workflowRun.GetStatus() == "waiting" || workflowRun.GetStatus() == "requested" ||
				workflowRun.GetStatus() == "pending" {
				client.Actions.CancelWorkflowRunByID(
					context.Background(), depl.RepoOwner, depl.RepoName, workflowRun.GetID(),
				)
			}
		}
	}

	// delete corresponding namespace
	agent, err := c.GetAgent(r, cluster, "")
	if err != nil {
		c.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// make sure we do not delete any kubernetes "system" namespaces
	if !isSystemNamespace(depl.Namespace) {
		agent.DeleteNamespace(depl.Namespace)
	}

	_, err = c.Repo().Environment().DeleteDeployment(depl)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.HandleAPIError(w, r, apierrors.NewErrNotFound(errDeploymentNotFound))
			return
		}

		c.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	originalBranches := strings.Split(env.GitDeployBranches, ",")
	newBranches := []string{}

	for _, branch := range originalBranches {
		if branch != depl.PRBranchFrom {
			newBranches = append(newBranches, branch)
		}
	}

	env.GitDeployBranches = strings.Join(newBranches, ",")

	_, err = c.Repo().Environment().UpdateEnvironment(env)

	if err != nil {
		c.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	if depl.GHDeploymentID != 0 {
		// set the GitHub deployment status to be inactive
		_, _, err := client.Repositories.CreateDeploymentStatus(
			context.Background(),
			env.GitRepoOwner,
			env.GitRepoName,
			depl.GHDeploymentID,
			&github.DeploymentStatusRequest{
				State: github.String("inactive"),
			},
		)
		if err != nil {
			c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
				fmt.Errorf("%v: %w", errGithubAPI, err), http.StatusConflict,
			))
			return
		}
	}

	c.WriteResult(w, r, depl.ToDeploymentType())
}
