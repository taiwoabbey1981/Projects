package project

import (
	"net/http"

	"github.com/porter-dev/porter/api/server/shared/apierrors"

	"github.com/porter-dev/porter/api/server/handlers"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/models"
)

type CollaboratorsListHandler struct {
	handlers.PorterHandlerWriter
}

func NewCollaboratorsListHandler(
	config *config.Config,
	writer shared.ResultWriter,
) *CollaboratorsListHandler {
	return &CollaboratorsListHandler{
		PorterHandlerWriter: handlers.NewDefaultPorterHandler(config, nil, writer),
	}
}

func (p *CollaboratorsListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proj, _ := r.Context().Value(types.ProjectScope).(*models.Project)

	roles, err := p.Repo().Project().ListProjectRoles(proj.ID)
	if err != nil {
		p.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	roleMap := make(map[uint]*models.Role)
	idArr := make([]uint, 0)

	for _, role := range roles {
		roleCp := role
		roleMap[role.UserID] = &roleCp
		idArr = append(idArr, role.UserID)
	}

	users, err := p.Repo().User().ListUsersByIDs(idArr)
	if err != nil {
		p.HandleAPIError(w, r, apierrors.NewErrInternal(err))
		return
	}

	var res types.ListCollaboratorsResponse = make([]*types.Collaborator, 0)

	for _, user := range users {
		res = append(res, &types.Collaborator{
			ID:        roleMap[user.ID].ID,
			Kind:      string(roleMap[user.ID].Kind),
			UserID:    roleMap[user.ID].UserID,
			Email:     user.Email,
			ProjectID: roleMap[user.ID].ProjectID,
		})
	}

	p.WriteResult(w, r, res)
}
