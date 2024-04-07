package porter_app

import (
	"context"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/porter-dev/porter/internal/deployment_target"
	"github.com/porter-dev/porter/internal/kubernetes"
	"github.com/porter-dev/porter/internal/porter_app"

	"github.com/porter-dev/porter/api/server/shared/requestutils"
	"github.com/porter-dev/porter/internal/kubernetes/environment_groups"

	"github.com/porter-dev/api-contracts/generated/go/helpers"
	porterv1 "github.com/porter-dev/api-contracts/generated/go/porter/v1"

	"github.com/porter-dev/porter/api/server/authz"
	"github.com/porter-dev/porter/api/server/handlers"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/apierrors"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/models"
	"github.com/porter-dev/porter/internal/telemetry"
)

// UpdateAppEnvironmentHandler handles the /apps/{porter_app_name}/update-environment endpoint
type UpdateAppEnvironmentHandler struct {
	handlers.PorterHandlerReadWriter
	authz.KubernetesAgentGetter
}

// NewUpdateAppEnvironmentHandler returns a new UpdateAppEnvironmentHandler
func NewUpdateAppEnvironmentHandler(
	config *config.Config,
	decoderValidator shared.RequestDecoderValidator,
	writer shared.ResultWriter,
) *UpdateAppEnvironmentHandler {
	return &UpdateAppEnvironmentHandler{
		PorterHandlerReadWriter: handlers.NewDefaultPorterHandler(config, decoderValidator, writer),
		KubernetesAgentGetter:   authz.NewOutOfClusterAgentGetter(config),
	}
}

const (
	// LabelKey_AppName is the label key for the app name
	LabelKey_AppName = "porter.run/app-name"
	// LabelKey_DeploymentTargetID is the label key for the deployment target id
	LabelKey_DeploymentTargetID = "porter.run/deployment-target-id"
	// LabelKey_PorterManaged is the label key signifying the resource is managed by porter
	LabelKey_PorterManaged = "porter.run/managed"
)

// UpdateAppEnvironmentRequest represents the accepted fields on a request to the /apps/{porter_app_name}/environment-group endpoint
type UpdateAppEnvironmentRequest struct {
	Base64AppProto     string            `json:"b64_app_proto"`
	DeploymentTargetID string            `json:"deployment_target_id"`
	Variables          map[string]string `json:"variables"`
	Secrets            map[string]string `json:"secrets"`
	// HardUpdate is used to remove any variables that are not specified in the request.  If false, the request will only update the variables specified in the request,
	// and leave all other variables untouched.
	HardUpdate bool `json:"remove_missing"`
}

// UpdateAppEnvironmentResponse represents the fields on the response object from the /apps/{porter_app_name}/environment-group endpoint
type UpdateAppEnvironmentResponse struct {
	Base64AppProto string                                `json:"b64_app_proto"`
	EnvGroups      []environment_groups.EnvironmentGroup `json:"env_groups"`
}

// ServeHTTP updates or creates the environment group for an app
func (c *UpdateAppEnvironmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.NewSpan(r.Context(), "serve-update-app-env-group")
	defer span.End()
	r = r.Clone(ctx)
	project, _ := ctx.Value(types.ProjectScope).(*models.Project)
	cluster, _ := ctx.Value(types.ClusterScope).(*models.Cluster)

	appName, reqErr := requestutils.GetURLParamString(r, types.URLParamPorterAppName)
	if reqErr != nil {
		err := telemetry.Error(ctx, span, nil, "error parsing porter app name")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
		return
	}
	telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "app-name", Value: appName})

	request := &UpdateAppEnvironmentRequest{}
	if ok := c.DecodeAndValidate(w, r, request); !ok {
		err := telemetry.Error(ctx, span, nil, "invalid request")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
		return
	}
	porterApp, err := c.Config().Repo.PorterApp().ReadPorterAppByName(cluster.ID, appName)
	if err != nil {
		err := telemetry.Error(ctx, span, nil, "error getting porter app by name")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
		return
	}
	if porterApp.ID == 0 {
		err := telemetry.Error(ctx, span, nil, "porter app not found")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusNotFound))
		return
	}
	telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "porter-app-id", Value: porterApp.ID})

	if request.DeploymentTargetID == "" {
		err := telemetry.Error(ctx, span, nil, "must provide deployment target id")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
		return
	}
	telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "deployment-target-id", Value: request.DeploymentTargetID})

	appProto := &porterv1.PorterApp{}

	if request.Base64AppProto == "" {
		if appName == "" {
			err := telemetry.Error(ctx, span, nil, "app name is empty and no base64 proto provided")
			c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
			return
		}

		appProto.Name = appName
	} else {
		decoded, err := base64.StdEncoding.DecodeString(request.Base64AppProto)
		if err != nil {
			err := telemetry.Error(ctx, span, err, "error decoding base yaml")
			c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
			return
		}

		err = helpers.UnmarshalContractObject(decoded, appProto)
		if err != nil {
			err := telemetry.Error(ctx, span, err, "error unmarshalling app proto")
			c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
			return
		}
	}

	if appProto.Name == "" {
		err := telemetry.Error(ctx, span, nil, "app proto name is empty")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
		return
	}

	deploymentTarget, err := deployment_target.DeploymentTargetDetails(ctx, deployment_target.DeploymentTargetDetailsInput{
		ProjectID:          int64(project.ID),
		ClusterID:          int64(cluster.ID),
		DeploymentTargetID: request.DeploymentTargetID,
		CCPClient:          c.Config().ClusterControlPlaneClient,
	})
	if err != nil {
		err := telemetry.Error(ctx, span, err, "error getting deployment target details")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	namespace := deploymentTarget.Namespace
	isPreview := deploymentTarget.IsPreview

	telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "is-preview", Value: isPreview})
	telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "namespace", Value: namespace})
	telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "hard-update", Value: request.HardUpdate})

	appEnvGroupName, err := porter_app.AppEnvGroupName(ctx, appName, request.DeploymentTargetID, cluster.ID, c.Repo().PorterApp())
	if err != nil {
		err := telemetry.Error(ctx, span, err, "error getting app env group name")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	agent, err := c.GetAgent(r, cluster, "")
	if err != nil {
		err := telemetry.Error(ctx, span, err, "unable to connect to kubernetes cluster")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	latestEnvironmentGroup, err := environment_groups.LatestBaseEnvironmentGroup(ctx, agent, appEnvGroupName)
	if err != nil {
		err := telemetry.Error(ctx, span, err, "unable to get latest base environment group")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "env-group-exists", Value: latestEnvironmentGroup.Name != ""})

	previewTemplateEnvName, err := porter_app.AppTemplateEnvGroupName(ctx, appName, cluster.ID, c.Repo().PorterApp())
	if err != nil {
		err := telemetry.Error(ctx, span, err, "error getting preview template env name")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	// filter out preview template and app env groups
	filteredEnvGroups := []*porterv1.EnvGroup{}
	for _, envGroup := range appProto.EnvGroups {
		if envGroup.GetName() != previewTemplateEnvName && envGroup.GetName() != appEnvGroupName {
			filteredEnvGroups = append(filteredEnvGroups, envGroup)
		}
	}

	if latestEnvironmentGroup.Name != "" {
		sameEnvGroup := true
		for key, newValue := range request.Variables {
			if existingValue, ok := latestEnvironmentGroup.Variables[key]; !ok || existingValue != newValue {
				sameEnvGroup = false
			}
		}
		for key, newValue := range request.Secrets {
			// We cannot check if the values are the same because the existing secrets are substituted with dummy values. However, if the new value is a dummy value, then it is unchanged.
			if _, ok := latestEnvironmentGroup.SecretVariables[key]; !ok || newValue != environment_groups.EnvGroupSecretDummyValue {
				sameEnvGroup = false
			}
		}
		if request.HardUpdate {
			for key := range latestEnvironmentGroup.Variables {
				if _, ok := request.Variables[key]; !ok {
					sameEnvGroup = false
				}
			}
			for key := range latestEnvironmentGroup.SecretVariables {
				if _, ok := request.Secrets[key]; !ok {
					sameEnvGroup = false
				}
			}
		}
		telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "same-env-group", Value: sameEnvGroup})

		if sameEnvGroup {
			// even if the env group is the same, we still need to sync the latest versions of the other env groups
			syncInp := syncLatestEnvGroupVersionsInput{
				envGroups:          filteredEnvGroups,
				appName:            appName,
				namespace:          namespace,
				deploymentTargetID: request.DeploymentTargetID,
				k8sAgent:           agent,
			}
			latestEnvGroups, err := syncLatestEnvGroupVersions(ctx, syncInp)
			if err != nil {
				err := telemetry.Error(ctx, span, err, "error syncing latest env group versions")
				c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
				return
			}

			latestEnvGroups = append(latestEnvGroups, environment_groups.EnvironmentGroup{
				Name:    latestEnvironmentGroup.Name,
				Version: latestEnvironmentGroup.Version,
			})

			var protoEnvGroups []*porterv1.EnvGroup
			for _, envGroup := range latestEnvGroups {
				protoEnvGroups = append(protoEnvGroups, &porterv1.EnvGroup{
					Name:    envGroup.Name,
					Version: int64(envGroup.Version),
				})
			}
			appProto.EnvGroups = protoEnvGroups

			encodedApp, err := encodeAppProto(ctx, appProto)
			if err != nil {
				err := telemetry.Error(ctx, span, err, "error encoding app proto")
				c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
				return
			}

			res := &UpdateAppEnvironmentResponse{
				EnvGroups:      latestEnvGroups,
				Base64AppProto: encodedApp,
			}

			c.WriteResult(w, r, res)
			return
		}
	}

	// if this app does not have a default env group for this deployment target and is a preview
	// then use the preview template env group as the default
	// this should only run when the app is first deployed to a given deployment target
	if latestEnvironmentGroup.Name == "" && isPreview {
		latestEnvironmentGroup, err = environment_groups.LatestBaseEnvironmentGroup(ctx, agent, previewTemplateEnvName)
		if err != nil {
			err := telemetry.Error(ctx, span, err, "unable to get latest base environment group")
			c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
			return
		}
	}

	variables := make(map[string]string)
	secrets := make(map[string]string)

	if !request.HardUpdate {
		for key, value := range latestEnvironmentGroup.Variables {
			variables[key] = value
		}
		for key, value := range latestEnvironmentGroup.SecretVariables {
			secrets[key] = value
		}
	}

	for key, value := range request.Variables {
		if len(key) > 0 && len(value) > 0 {
			variables[key] = value
		}
	}
	for key, value := range request.Secrets {
		if len(key) > 0 && len(value) > 0 {
			secrets[key] = value
		}
	}

	envGroup := environment_groups.EnvironmentGroup{
		Name:            appEnvGroupName,
		Variables:       variables,
		SecretVariables: secrets,
		CreatedAtUTC:    time.Now().UTC(),
	}

	additionalEnvGroupLabels := map[string]string{
		LabelKey_AppName:                                  appName,
		LabelKey_DeploymentTargetID:                       request.DeploymentTargetID,
		environment_groups.LabelKey_DefaultAppEnvironment: "true",
		LabelKey_PorterManaged:                            "true",
	}

	err = environment_groups.CreateOrUpdateBaseEnvironmentGroup(ctx, agent, envGroup, additionalEnvGroupLabels)
	if err != nil {
		err := telemetry.Error(ctx, span, err, "unable to create or update base environment group")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	inp := environment_groups.SyncLatestVersionToNamespaceInput{
		BaseEnvironmentGroupName: appEnvGroupName,
		TargetNamespace:          namespace,
	}

	syncedAppEnvironment, err := environment_groups.SyncLatestVersionToNamespace(ctx, agent, inp, additionalEnvGroupLabels)
	if err != nil {
		err := telemetry.Error(ctx, span, err, "unable to create or update synced environment group")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}
	telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "env-group-versioned-name", Value: syncedAppEnvironment.EnvironmentGroupVersionedName})

	syncInp := syncLatestEnvGroupVersionsInput{
		envGroups:          filteredEnvGroups,
		appName:            appName,
		namespace:          namespace,
		deploymentTargetID: request.DeploymentTargetID,
		k8sAgent:           agent,
	}
	latestEnvGroups, err := syncLatestEnvGroupVersions(ctx, syncInp)
	if err != nil {
		err := telemetry.Error(ctx, span, err, "error syncing latest env group versions")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	split := strings.Split(syncedAppEnvironment.EnvironmentGroupVersionedName, ".")
	if len(split) != 2 {
		err := telemetry.Error(ctx, span, err, "unexpected environment group versioned name")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	version, err := strconv.Atoi(split[1])
	if err != nil {
		err := telemetry.Error(ctx, span, err, "error converting environment group version to int")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	latestEnvGroups = append(latestEnvGroups, environment_groups.EnvironmentGroup{
		Name:    split[0],
		Version: version,
	})

	var protoEnvGroups []*porterv1.EnvGroup
	for _, envGroup := range latestEnvGroups {
		protoEnvGroups = append(protoEnvGroups, &porterv1.EnvGroup{
			Name:    envGroup.Name,
			Version: int64(envGroup.Version),
		})
	}
	appProto.EnvGroups = protoEnvGroups

	encodedApp, err := encodeAppProto(ctx, appProto)
	if err != nil {
		err := telemetry.Error(ctx, span, err, "error encoding app proto")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	res := &UpdateAppEnvironmentResponse{
		EnvGroups:      latestEnvGroups,
		Base64AppProto: encodedApp,
	}

	c.WriteResult(w, r, res)
}

type syncLatestEnvGroupVersionsInput struct {
	// envGroups is the list of env groups to sync. We only need the names and will get the latest version of each from the porter-env-group ns
	envGroups []*porterv1.EnvGroup
	// appName is the name of the app
	appName string
	// namespace is the namespace to sync the latest versions to
	namespace string
	// deploymentTargetID is the id of the deployment target
	deploymentTargetID string
	// k8sAgent is the kubernetes agent
	k8sAgent *kubernetes.Agent
}

// syncLatestEnvGroupVersions syncs the latest versions of the env groups to the namespace where an app is deployed
func syncLatestEnvGroupVersions(ctx context.Context, inp syncLatestEnvGroupVersionsInput) ([]environment_groups.EnvironmentGroup, error) {
	ctx, span := telemetry.NewSpan(ctx, "sync-latest-env-group-versions")
	defer span.End()

	var envGroups []environment_groups.EnvironmentGroup

	if inp.deploymentTargetID == "" {
		return envGroups, telemetry.Error(ctx, span, nil, "deployment target id is empty")
	}
	if inp.appName == "" {
		return envGroups, telemetry.Error(ctx, span, nil, "app name is empty")
	}
	if inp.namespace == "" {
		return envGroups, telemetry.Error(ctx, span, nil, "namespace is empty")
	}
	if inp.k8sAgent == nil {
		return envGroups, telemetry.Error(ctx, span, nil, "k8s agent is nil")
	}

	for _, envGroup := range inp.envGroups {
		if envGroup == nil {
			continue
		}

		additionalEnvGroupLabels := map[string]string{
			LabelKey_AppName:            inp.appName,
			LabelKey_DeploymentTargetID: inp.deploymentTargetID,
			LabelKey_PorterManaged:      "true",
		}

		syncedEnvironment, err := environment_groups.SyncLatestVersionToNamespace(ctx, inp.k8sAgent, environment_groups.SyncLatestVersionToNamespaceInput{
			TargetNamespace:          inp.namespace,
			BaseEnvironmentGroupName: envGroup.GetName(),
		}, additionalEnvGroupLabels)
		if err != nil {
			telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "env-group-name", Value: envGroup.GetName()})
			return envGroups, telemetry.Error(ctx, span, err, "error syncing latest version to namespace")
		}

		split := strings.Split(syncedEnvironment.EnvironmentGroupVersionedName, ".")
		if len(split) != 2 {
			return envGroups, telemetry.Error(ctx, span, err, "unexpected environment group versioned name")
		}

		version, err := strconv.Atoi(split[1])
		if err != nil {
			return envGroups, telemetry.Error(ctx, span, err, "error converting environment group version to int")
		}

		envGroups = append(envGroups, environment_groups.EnvironmentGroup{
			Name:    split[0],
			Version: version,
		})
	}

	return envGroups, nil
}
