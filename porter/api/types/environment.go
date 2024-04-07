package types

import "time"

type Environment struct {
	ID                uint     `json:"id"`
	ProjectID         uint     `json:"project_id"`
	ClusterID         uint     `json:"cluster_id"`
	GitInstallationID uint     `json:"git_installation_id"`
	GitRepoOwner      string   `json:"git_repo_owner"`
	GitRepoName       string   `json:"git_repo_name"`
	GitRepoBranches   []string `json:"git_repo_branches"`

	Name                 string            `json:"name"`
	Mode                 string            `json:"mode"`
	DeploymentCount      uint              `json:"deployment_count"`
	LastDeploymentStatus string            `json:"last_deployment_status"`
	NewCommentsDisabled  bool              `json:"new_comments_disabled"`
	NamespaceLabels      map[string]string `json:"namespace_labels,omitempty"`
	GitDeployBranches    []string          `json:"git_deploy_branches"`
}

type CreateEnvironmentRequest struct {
	Name               string            `json:"name" form:"required"`
	Mode               string            `json:"mode" form:"oneof=auto manual" default:"manual"`
	DisableNewComments bool              `json:"disable_new_comments"`
	GitRepoBranches    []string          `json:"git_repo_branches"`
	NamespaceLabels    map[string]string `json:"namespace_labels"`
	GitDeployBranches  []string          `json:"git_deploy_branches"`
}

type GitHubMetadata struct {
	DeploymentID int64  `json:"gh_deployment_id"`
	PRName       string `json:"gh_pr_name"`
	RepoName     string `json:"gh_repo_name"`
	RepoOwner    string `json:"gh_repo_owner"`
	CommitSHA    string `json:"gh_commit_sha"`
	PRBranchFrom string `json:"gh_pr_branch_from"`
	PRBranchInto string `json:"gh_pr_branch_into"`
}

type DeploymentStatus string

const (
	DeploymentStatusCreated  DeploymentStatus = "created"
	DeploymentStatusCreating DeploymentStatus = "creating"
	DeploymentStatusUpdating DeploymentStatus = "updating"
	DeploymentStatusInactive DeploymentStatus = "inactive"
	DeploymentStatusTimedOut DeploymentStatus = "timed_out"
	DeploymentStatusFailed   DeploymentStatus = "failed"
)

type Deployment struct {
	*GitHubMetadata

	ID                 uint             `json:"id"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	EnvironmentID      uint             `json:"environment_id"`
	Namespace          string           `json:"namespace"`
	Status             DeploymentStatus `json:"status"`
	Subdomain          string           `json:"subdomain"`
	PullRequestID      uint             `json:"pull_request_id"`
	InstallationID     uint             `json:"gh_installation_id"`
	LastWorkflowRunURL string           `json:"last_workflow_run_url"`
	LastErrors         string           `json:"last_errors"`
}

type CreateGHDeploymentRequest struct {
	ActionID uint `json:"action_id" form:"required"`
}

type CreateDeploymentRequest struct {
	*CreateGHDeploymentRequest
	*GitHubMetadata

	Namespace     string `json:"namespace" form:"required"`
	PullRequestID uint   `json:"pull_request_id"`
}

type SuccessfullyDeployedResource struct {
	ReleaseName string `json:"release_name" form:"required"`
	ReleaseType string `json:"release_type"`
}

type FinalizeDeploymentRequest struct {
	SuccessfulResources []*SuccessfullyDeployedResource `json:"successful_resources"`
	Subdomain           string                          `json:"subdomain"`
	PRNumber            uint                            `json:"pr_number"`
	Namespace           string                          `json:"namespace"`
}

type FinalizeDeploymentByClusterRequest struct {
	RepoOwner string `json:"gh_repo_owner" form:"required"`
	RepoName  string `json:"gh_repo_name" form:"required"`

	SuccessfulResources []*SuccessfullyDeployedResource `json:"successful_resources"`
	Subdomain           string                          `json:"subdomain"`
	PRNumber            uint                            `json:"pr_number"`
	Namespace           string                          `json:"namespace"`
}

type FinalizeDeploymentWithErrorsRequest struct {
	SuccessfulResources []*SuccessfullyDeployedResource `json:"successful_resources"`
	Errors              map[string]string               `json:"errors" form:"required"`
	PRNumber            uint                            `json:"pr_number"`
	Namespace           string                          `json:"namespace"`
}

type FinalizeDeploymentWithErrorsByClusterRequest struct {
	RepoOwner string `json:"gh_repo_owner" form:"required"`
	RepoName  string `json:"gh_repo_name" form:"required"`

	SuccessfulResources []*SuccessfullyDeployedResource `json:"successful_resources"`
	Errors              map[string]string               `json:"errors" form:"required"`
	PRNumber            uint                            `json:"pr_number"`
	Namespace           string                          `json:"namespace"`
}

type UpdateDeploymentRequest struct {
	*CreateGHDeploymentRequest

	PRBranchFrom string `json:"gh_pr_branch_from" form:"required"`
	CommitSHA    string `json:"commit_sha" form:"required"`
	PRNumber     uint   `json:"pr_number"`
	Namespace    string `json:"namespace"`
}

type UpdateDeploymentByClusterRequest struct {
	*CreateGHDeploymentRequest

	RepoOwner string `json:"gh_repo_owner" form:"required"`
	RepoName  string `json:"gh_repo_name" form:"required"`

	PRBranchFrom string `json:"gh_pr_branch_from" form:"required"`
	CommitSHA    string `json:"commit_sha" form:"required"`
	PRNumber     uint   `json:"pr_number"`
	Namespace    string `json:"namespace"`
}

type ListDeploymentRequest struct {
	EnvironmentID uint `schema:"environment_id"`
}

type UpdateDeploymentStatusRequest struct {
	*CreateGHDeploymentRequest

	PRBranchFrom string `json:"gh_pr_branch_from" form:"required"`
	Status       string `json:"status" form:"required,oneof=created creating inactive failed"`
	PRNumber     uint   `json:"pr_number"`
	Namespace    string `json:"namespace"`
}

type UpdateDeploymentStatusByClusterRequest struct {
	*CreateGHDeploymentRequest

	RepoOwner string `json:"gh_repo_owner" form:"required"`
	RepoName  string `json:"gh_repo_name" form:"required"`

	PRBranchFrom string `json:"gh_pr_branch_from" form:"required"`
	Status       string `json:"status" form:"required,oneof=created creating inactive failed"`
	PRNumber     uint   `json:"pr_number"`
	Namespace    string `json:"namespace"`
}

type DeleteDeploymentRequest struct {
	Namespace string `json:"namespace" form:"required"`
}

type GetDeploymentRequest struct {
	DeploymentID uint   `schema:"id"`
	PRNumber     uint   `schema:"pr_number"`
	Namespace    string `schema:"namespace"`
	Branch       string `schema:"branch"`
}

type PullRequest struct {
	Title      string    `json:"pr_title"`
	Number     uint      `json:"pr_number"`
	RepoOwner  string    `json:"repo_owner"`
	RepoName   string    `json:"repo_name"`
	BranchFrom string    `json:"branch_from"`
	BranchInto string    `json:"branch_into"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ToggleNewCommentRequest struct {
	Disable bool `json:"disable"`
}

type ListEnvironmentsResponse []*Environment

type ValidatePorterYAMLRequest struct {
	Branch string `schema:"branch"`
}

type ValidatePorterYAMLResponse struct {
	Errors []string `json:"errors"`
}

type UpdateEnvironmentSettingsRequest struct {
	Mode               string            `json:"mode" form:"oneof=auto manual"`
	DisableNewComments bool              `json:"disable_new_comments"`
	GitRepoBranches    []string          `json:"git_repo_branches"`
	NamespaceLabels    map[string]string `json:"namespace_labels"`
	GitDeployBranches  []string          `json:"git_deploy_branches"`
}
