package gitinstallation

import (
	"net/http"

	"github.com/porter-dev/porter/api/server/handlers"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/models/integrations"
)

type GitInstallationGetHandler struct {
	handlers.PorterHandlerWriter
}

func NewGitInstallationGetHandler(
	config *config.Config,
	writer shared.ResultWriter,
) *GitInstallationGetHandler {
	return &GitInstallationGetHandler{
		PorterHandlerWriter: handlers.NewDefaultPorterHandler(config, nil, writer),
	}
}

func (c *GitInstallationGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ga, _ := r.Context().Value(types.GitInstallationScope).(*integrations.GithubAppInstallation)

	c.WriteResult(w, r, ga.ToGitInstallationType())
}
