package environment

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/go-github/v41/github"
	"github.com/porter-dev/porter/api/server/handlers"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/apierrors"
	"github.com/porter-dev/porter/api/server/shared/commonutils"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/models"
	"github.com/porter-dev/porter/internal/telemetry"
)

type ListDeploymentsByClusterHandler struct {
	handlers.PorterHandlerReadWriter
}

func NewListDeploymentsByClusterHandler(
	config *config.Config,
	decoderValidator shared.RequestDecoderValidator,
	writer shared.ResultWriter,
) *ListDeploymentsByClusterHandler {
	return &ListDeploymentsByClusterHandler{
		PorterHandlerReadWriter: handlers.NewDefaultPorterHandler(config, decoderValidator, writer),
	}
}

func (c *ListDeploymentsByClusterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.NewSpan(r.Context(), "serve-list-cluster-deployments")
	defer span.End()

	project, _ := ctx.Value(types.ProjectScope).(*models.Project)
	cluster, _ := ctx.Value(types.ClusterScope).(*models.Cluster)

	req := &types.ListDeploymentRequest{}
	telemetry.WithAttributes(span,
		telemetry.AttributeKV{Key: "cluster-id", Value: cluster.ID},
		telemetry.AttributeKV{Key: "project-id", Value: project.ID},
		telemetry.AttributeKV{Key: "environment-id", Value: req.EnvironmentID},
	)

	if ok := c.DecodeAndValidate(w, r, req); !ok {
		return
	}

	var deployments []*types.Deployment
	var pullRequests []*types.PullRequest

	if req.EnvironmentID == 0 {
		depls, err := c.Repo().Environment().ListDeploymentsByCluster(project.ID, cluster.ID)
		if err != nil {
			err = telemetry.Error(ctx, span, err, "failed to list deployments from cluster")
			c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
			return
		}

		deplInfoMap := make(map[string]bool)

		for _, depl := range depls {
			deployment := depl.ToDeploymentType()
			deplInfoMap[fmt.Sprintf(
				"%s-%s-%d", deployment.RepoOwner, deployment.RepoName, deployment.PullRequestID,
			)] = true

			env, err := c.Repo().Environment().ReadEnvironmentByID(project.ID, cluster.ID, deployment.EnvironmentID)
			if err != nil {
				err = telemetry.Error(ctx, span, err, "failed to get environment from deployment")
				c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
				return
			}

			deployment.InstallationID = env.GitInstallationID

			deployments = append(deployments, deployment)
		}

		envToGithubClientMap := make(map[uint]*github.Client)

		var wg sync.WaitGroup
		wg.Add(len(deployments))

		for _, deployment := range deployments {
			env, err := c.Repo().Environment().ReadEnvironmentByID(project.ID, cluster.ID, deployment.EnvironmentID)
			if err != nil {
				err = telemetry.Error(ctx, span, err, "failed to get environment from deployment")
				c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
				return
			}

			if _, ok := envToGithubClientMap[env.ID]; !ok {
				client, err := getGithubClientFromEnvironment(c.Config(), env)
				if err != nil {
					err = telemetry.Error(ctx, span, err, "error getting github client from environment")
					c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
					return
				}

				envToGithubClientMap[env.ID] = client
			}

			go func(depl *types.Deployment) {
				defer wg.Done()

				updateDeploymentWithGithubWorkflowRunStatus(c.Config(), envToGithubClientMap[env.ID], env, depl)
			}(deployment)
		}

		wg.Wait()

		envList, err := c.Repo().Environment().ListEnvironments(project.ID, cluster.ID)
		if err != nil {
			c.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}

		for _, env := range envList {
			if _, ok := envToGithubClientMap[env.ID]; !ok {
				client, err := getGithubClientFromEnvironment(c.Config(), env)
				if err != nil {
					err = telemetry.Error(ctx, span, err, "error getting github client from environment")
					c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
					return
				}

				envToGithubClientMap[env.ID] = client
			}

			prs, err := fetchOpenPullRequests(ctx, c.Config(), envToGithubClientMap[env.ID], env, deplInfoMap)
			if err != nil {
				err = telemetry.Error(ctx, span, err, "error fetching pull requests")
				c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
				return
			}

			pullRequests = append(pullRequests, prs...)
		}
	} else {
		env, err := c.Repo().Environment().ReadEnvironmentByID(project.ID, cluster.ID, req.EnvironmentID)
		if err != nil {
			err = telemetry.Error(ctx, span, err, "error fetching environment")
			c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
			return
		}

		depls, err := c.Repo().Environment().ListDeployments(env.ID)
		if err != nil {
			err = telemetry.Error(ctx, span, err, "error listing deployments")
			c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
			return
		}

		deplInfoMap := make(map[string]bool)

		client, err := getGithubClientFromEnvironment(c.Config(), env)
		if err != nil {
			err = telemetry.Error(ctx, span, err, "error getting github client from environment")
			c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
			return
		}

		for _, depl := range depls {
			deployment := depl.ToDeploymentType()
			deplInfoMap[fmt.Sprintf(
				"%s-%s-%d", deployment.RepoOwner, deployment.RepoName, deployment.PullRequestID,
			)] = true

			deployment.InstallationID = env.GitInstallationID

			deployments = append(deployments, deployment)
		}

		var wg sync.WaitGroup
		wg.Add(len(deployments))

		for _, deployment := range deployments {
			go func(depl *types.Deployment) {
				defer wg.Done()

				updateDeploymentWithGithubWorkflowRunStatus(c.Config(), client, env, depl)
			}(deployment)
		}

		wg.Wait()

		prs, err := fetchOpenPullRequests(ctx, c.Config(), client, env, deplInfoMap)
		if err != nil {
			err = telemetry.Error(ctx, span, err, "error fetching pull requests")
			c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
			return
		}

		pullRequests = append(pullRequests, prs...)
	}

	c.WriteResult(w, r, map[string]interface{}{
		"pull_requests": pullRequests,
		"deployments":   deployments,
	})
}

func updateDeploymentWithGithubWorkflowRunStatus(
	config *config.Config,
	client *github.Client,
	env *models.Environment,
	deployment *types.Deployment,
) {
	if deployment.Status == types.DeploymentStatusInactive {
		return
	}

	latestWorkflowRun, err := commonutils.GetLatestWorkflowRun(client, env.GitRepoOwner, env.GitRepoName,
		fmt.Sprintf("porter_%s_env.yml", env.Name), deployment.PRBranchFrom)

	if err == nil {
		deployment.LastWorkflowRunURL = latestWorkflowRun.GetHTMLURL()

		if (latestWorkflowRun.GetStatus() == "in_progress" ||
			latestWorkflowRun.GetStatus() == "queued") &&
			deployment.Status != types.DeploymentStatusCreating {
			deployment.Status = types.DeploymentStatusUpdating
		} else if latestWorkflowRun.GetStatus() == "completed" {
			if latestWorkflowRun.GetConclusion() == "failure" {
				deployment.Status = types.DeploymentStatusFailed
			} else if latestWorkflowRun.GetConclusion() == "timed_out" {
				deployment.Status = types.DeploymentStatusTimedOut
			} else if latestWorkflowRun.GetConclusion() == "success" {
				deployment.Status = types.DeploymentStatusCreated
			}
		}
	}
}

func fetchOpenPullRequests(
	ctx context.Context,
	config *config.Config,
	client *github.Client,
	env *models.Environment,
	deplInfoMap map[string]bool,
) ([]*types.PullRequest, error) {
	branchesMap := make(map[string]bool)

	for _, br := range env.ToEnvironmentType().GitRepoBranches {
		branchesMap[br] = true
	}

	openPRs, resp, err := client.PullRequests.List(ctx, env.GitRepoOwner, env.GitRepoName,
		&github.PullRequestListOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		},
	)

	var prs []*types.PullRequest

	if resp != nil && resp.StatusCode == 404 {
		return prs, nil
	}

	if err != nil {
		return nil, err
	}

	var ghPRs []*github.PullRequest

	for resp.NextPage != 0 && err == nil {
		ghPRs, resp, err = client.PullRequests.List(ctx, env.GitRepoOwner, env.GitRepoName,
			&github.PullRequestListOptions{
				ListOptions: github.ListOptions{
					PerPage: 100,
					Page:    resp.NextPage,
				},
			},
		)

		openPRs = append(openPRs, ghPRs...)
	}

	for _, pr := range openPRs {
		if len(branchesMap) > 0 {
			if _, ok := branchesMap[pr.GetBase().GetRef()]; !ok {
				continue
			}
		}

		if isDeployBranch(pr.GetHead().GetRef(), env) {
			continue
		}

		if _, ok := deplInfoMap[fmt.Sprintf("%s-%s-%d", env.GitRepoOwner, env.GitRepoName, pr.GetNumber())]; !ok {
			prs = append(prs, &types.PullRequest{
				Title:      pr.GetTitle(),
				Number:     uint(pr.GetNumber()),
				RepoOwner:  env.GitRepoOwner,
				RepoName:   env.GitRepoName,
				BranchFrom: pr.GetHead().GetRef(),
				BranchInto: pr.GetBase().GetRef(),
				CreatedAt:  pr.GetCreatedAt(),
				UpdatedAt:  pr.GetUpdatedAt(),
			})
		}
	}

	return prs, nil
}

func isDeployBranch(branch string, env *models.Environment) bool {
	for _, b := range env.ToEnvironmentType().GitDeployBranches {
		if b == branch {
			return true
		}
	}

	return false
}
