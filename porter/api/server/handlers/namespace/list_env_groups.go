package namespace

import (
	"net/http"

	"github.com/porter-dev/porter/api/server/authz"
	"github.com/porter-dev/porter/api/server/handlers"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/apierrors"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/kubernetes/envgroup"
	"github.com/porter-dev/porter/internal/models"
)

type ListEnvGroupsHandler struct {
	handlers.PorterHandlerWriter
	authz.KubernetesAgentGetter
}

func NewListEnvGroupsHandler(
	config *config.Config,
	writer shared.ResultWriter,
) *ListEnvGroupsHandler {
	return &ListEnvGroupsHandler{
		PorterHandlerWriter:   handlers.NewDefaultPorterHandler(config, nil, writer),
		KubernetesAgentGetter: authz.NewOutOfClusterAgentGetter(config),
	}
}

func (c *ListEnvGroupsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	namespace := r.Context().Value(types.NamespaceScope).(string)
	cluster, _ := r.Context().Value(types.ClusterScope).(*models.Cluster)

	agent, err := c.GetAgent(r, cluster, "")
	if err != nil {
		c.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// get all versioned config maps
	configMaps, err := agent.ListAllVersionedConfigMaps(namespace)
	if err != nil {
		c.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	res := make(types.ListEnvGroupsResponse, 0)

	for _, cm := range configMaps {
		eg, err := envgroup.ToEnvGroup(&cm)
		if err != nil {
			continue
		}

		res = append(res, &types.EnvGroupMeta{
			MetaVersion: eg.MetaVersion,
			CreatedAt:   eg.CreatedAt,
			Name:        eg.Name,
			Namespace:   eg.Namespace,
			Version:     eg.Version,
		})
	}

	// get all meta-version 1 configmaps
	configMapList, err := agent.ListConfigMaps(namespace)
	if err != nil {
		c.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	for _, v1CM := range configMapList.Items {
		eg, err := envgroup.ToEnvGroup(&v1CM)
		if err != nil {
			continue
		}

		res = append(res, &types.EnvGroupMeta{
			MetaVersion: eg.MetaVersion,
			CreatedAt:   eg.CreatedAt,
			Name:        eg.Name,
			Namespace:   eg.Namespace,
			Version:     eg.Version,
		})
	}

	c.WriteResult(w, r, res)
}
