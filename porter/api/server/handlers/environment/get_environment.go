package environment

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/porter-dev/porter/internal/telemetry"

	"github.com/porter-dev/porter/api/server/handlers"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/apierrors"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/server/shared/requestutils"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/models"
	"gorm.io/gorm"
)

type GetEnvironmentHandler struct {
	handlers.PorterHandlerWriter
}

func NewGetEnvironmentHandler(
	config *config.Config,
	writer shared.ResultWriter,
) *GetEnvironmentHandler {
	return &GetEnvironmentHandler{
		PorterHandlerWriter: handlers.NewDefaultPorterHandler(config, nil, writer),
	}
}

func (c *GetEnvironmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.NewSpan(r.Context(), "serve-get-environment")
	defer span.End()

	project, _ := ctx.Value(types.ProjectScope).(*models.Project)
	cluster, _ := ctx.Value(types.ClusterScope).(*models.Cluster)

	telemetry.WithAttributes(span,
		telemetry.AttributeKV{Key: "project-id", Value: project.ID},
		telemetry.AttributeKV{Key: "cluster-id", Value: cluster.ID},
	)

	envID, reqErr := requestutils.GetURLParamUint(r, "environment_id")
	if reqErr != nil {
		_ = telemetry.Error(ctx, span, reqErr, "could not get environment id from url")
		c.HandleAPIError(w, r, reqErr)
		return
	}

	telemetry.WithAttributes(span, telemetry.AttributeKV{Key: "environment-id", Value: envID})

	env, err := c.Repo().Environment().ReadEnvironmentByID(project.ID, cluster.ID, envID)
	if err != nil {
		_ = telemetry.Error(ctx, span, err, "could not read environment by id")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.HandleAPIError(w, r, apierrors.NewErrNotFound(fmt.Errorf("no such environment with ID: %d", envID)))
			return
		}

		c.HandleAPIError(w, r, apierrors.NewErrInternal(fmt.Errorf("error reading environment with ID: %d. Error: %w", envID, err)))
		return
	}

	c.WriteResult(w, r, env.ToEnvironmentType())
}
