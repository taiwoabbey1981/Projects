package porter_app

import (
	"net/http"
	"strings"

	"github.com/porter-dev/porter/api/server/authz"
	"github.com/porter-dev/porter/api/server/handlers"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/apierrors"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/server/shared/features"
	"github.com/porter-dev/porter/api/server/shared/requestutils"
	"github.com/porter-dev/porter/api/types"
	utils "github.com/porter-dev/porter/api/utils/porter_app"
	"github.com/porter-dev/porter/internal/helm"
	"github.com/porter-dev/porter/internal/models"
	"github.com/porter-dev/porter/internal/telemetry"
	"gopkg.in/yaml.v2"
)

type RollbackPorterAppHandler struct {
	handlers.PorterHandlerReadWriter
	authz.KubernetesAgentGetter
}

func NewRollbackPorterAppHandler(
	config *config.Config,
	decoderValidator shared.RequestDecoderValidator,
	writer shared.ResultWriter,
) *RollbackPorterAppHandler {
	return &RollbackPorterAppHandler{
		PorterHandlerReadWriter: handlers.NewDefaultPorterHandler(config, decoderValidator, writer),
		KubernetesAgentGetter:   authz.NewOutOfClusterAgentGetter(config),
	}
}

func (c *RollbackPorterAppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.NewSpan(r.Context(), "serve-rollback-porter-app")
	defer span.End()
	cluster, _ := ctx.Value(types.ClusterScope).(*models.Cluster)

	request := &types.RollbackPorterAppRequest{}
	if ok := c.DecodeAndValidate(w, r, request); !ok {
		err := telemetry.Error(ctx, span, nil, "error decoding request")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
		return
	}

	appName, reqErr := requestutils.GetURLParamString(r, types.URLParamPorterAppName)
	if reqErr != nil {
		err := telemetry.Error(ctx, span, reqErr, "error getting stack name from url")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
		return
	}
	telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "stack-name", Value: appName})
	namespace := utils.NamespaceFromPorterAppName(appName)

	helmAgent, err := c.GetHelmAgent(ctx, r, cluster, namespace)
	if err != nil {
		err = telemetry.Error(ctx, span, err, "error getting helm agent")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	k8sAgent, err := c.GetAgent(r, cluster, namespace)
	if err != nil {
		err = telemetry.Error(ctx, span, err, "error getting k8s agent")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	helmReleaseFromRequestedRevision, err := helmAgent.GetRelease(ctx, appName, request.Revision, false)
	if err != nil {
		err = telemetry.Error(ctx, span, err, "error getting helm release for requested revision")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	latestHelmRelease, err := helmAgent.GetRelease(ctx, appName, 0, false)
	if err != nil {
		err = telemetry.Error(ctx, span, err, "error getting latest helm release")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	valuesYaml, err := yaml.Marshal(helmReleaseFromRequestedRevision.Config)
	if err != nil {
		err = telemetry.Error(ctx, span, err, "error marshalling helm release config to yaml")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	imageInfo := attemptToGetImageInfoFromRelease(helmReleaseFromRequestedRevision.Config)
	if imageInfo.Tag == "" {
		imageInfo.Tag = "latest"
	}

	porterApp, err := c.Repo().PorterApp().ReadPorterAppByName(cluster.ID, appName)
	if err != nil {
		err = telemetry.Error(ctx, span, err, "error getting porter app")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}
	injectLauncher := strings.Contains(porterApp.Builder, "heroku") ||
		strings.Contains(porterApp.Builder, "paketo")

	registries, err := c.Repo().Registry().ListRegistriesByProjectID(cluster.ProjectID)
	if err != nil {
		err = telemetry.Error(ctx, span, err, "error listing registries")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	chart, values, _, err := parse(
		ctx,
		ParseConf{
			PorterAppName: appName,
			ImageInfo:     imageInfo,
			ServerConfig:  c.Config(),
			ProjectID:     cluster.ProjectID,
			Namespace:     namespace,
			SubdomainCreateOpts: SubdomainCreateOpts{
				k8sAgent:      k8sAgent,
				dnsRepo:       c.Repo().DNSRecord(),
				dnsClient:     c.Config().DNSClient,
				appRootDomain: c.Config().ServerConf.AppRootDomain,
				stackName:     appName,
			},
			InjectLauncherToStartCommand: injectLauncher,
			FullHelmValues:               string(valuesYaml),
		},
	)
	if err != nil {
		err = telemetry.Error(ctx, span, err, "error parsing helm chart")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}

	conf := &helm.InstallChartConfig{
		Chart:      chart,
		Name:       appName,
		Namespace:  namespace,
		Values:     values,
		Cluster:    cluster,
		Repo:       c.Repo(),
		Registries: registries,
	}
	_, err = helmAgent.UpgradeInstallChart(ctx, conf, c.Config().DOConf, c.Config().ServerConf.DisablePullSecretsInjection)
	if err != nil {
		err = telemetry.Error(ctx, span, err, "error upgrading application")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusBadRequest))
		return
	}

	if features.AreAgentDeployEventsEnabled(k8sAgent) {
		serviceDeploymentStatusMap := getServiceDeploymentMetadataFromValues(values, types.PorterAppEventStatus_Progressing)
		_, err = createNewPorterAppDeployEvent(ctx, serviceDeploymentStatusMap, porterApp.ID, latestHelmRelease.Version+1, imageInfo.Tag, c.Repo().PorterAppEvent())
	} else {
		_, err = createOldPorterAppDeployEvent(ctx, types.PorterAppEventStatus_Success, porterApp.ID, latestHelmRelease.Version+1, imageInfo.Tag, c.Repo().PorterAppEvent())
	}
	if err != nil {
		err = telemetry.Error(ctx, span, err, "error creating porter app event")
		c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(err, http.StatusInternalServerError))
		return
	}
}
