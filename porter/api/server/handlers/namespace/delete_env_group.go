package namespace

import (
	"fmt"
	"net/http"

	"github.com/porter-dev/porter/api/server/authz"
	"github.com/porter-dev/porter/api/server/handlers"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/apierrors"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/kubernetes"
	"github.com/porter-dev/porter/internal/kubernetes/envgroup"
	"github.com/porter-dev/porter/internal/models"
)

type DeleteEnvGroupHandler struct {
	handlers.PorterHandlerReadWriter
	authz.KubernetesAgentGetter
}

func NewDeleteEnvGroupHandler(
	config *config.Config,
	decoderValidator shared.RequestDecoderValidator,
	writer shared.ResultWriter,
) *DeleteEnvGroupHandler {
	return &DeleteEnvGroupHandler{
		PorterHandlerReadWriter: handlers.NewDefaultPorterHandler(config, decoderValidator, writer),
		KubernetesAgentGetter:   authz.NewOutOfClusterAgentGetter(config),
	}
}

func (c *DeleteEnvGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := &types.DeleteEnvGroupRequest{}

	if ok := c.DecodeAndValidate(w, r, request); !ok {
		return
	}

	namespace := r.Context().Value(types.NamespaceScope).(string)
	cluster, _ := r.Context().Value(types.ClusterScope).(*models.Cluster)

	agent, err := c.GetAgent(r, cluster, namespace)
	if err != nil {
		c.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	// get the env group: if it's MetaVersion=2, return an error
	envGroup, err := envgroup.GetEnvGroup(agent, request.Name, namespace, 0)
	if err != nil {
		c.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	if envGroup != nil && envGroup.MetaVersion == 1 {
		if err := deleteV1ConfigMap(agent, request.Name, namespace); err != nil {
			c.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}
	} else if envGroup != nil && envGroup.MetaVersion == 2 {
		if len(envGroup.Applications) != 0 {
			c.HandleAPIError(w, r, apierrors.NewErrPassThroughToClient(
				fmt.Errorf("env group must not have any connected applications"),
				http.StatusNotFound,
			))

			return
		} else if err = envgroup.DeleteEnvGroup(agent, request.Name, namespace); err != nil {
			c.HandleAPIError(w, r, apierrors.NewErrInternal(err))
			return
		}
	}
}

func deleteV1ConfigMap(agent *kubernetes.Agent, name, namespace string) error {
	if err := agent.DeleteLinkedSecret(name, namespace); err != nil {
		return err
	}

	if err := agent.DeleteConfigMap(name, namespace); err != nil {
		return err
	}

	return nil
}
