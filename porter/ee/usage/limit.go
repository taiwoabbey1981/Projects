//go:build ee
// +build ee

package usage

import (
	"errors"

	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/models"
	"github.com/porter-dev/porter/internal/repository"
	"gorm.io/gorm"
)

func GetLimit(repo repository.Repository, proj *models.Project) (limit *types.ProjectUsage, err error) {
	// query for the project limit; if not found, no limits
	limitModel, err := repo.ProjectUsage().ReadProjectUsage(proj.ID)

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// place existing users without usage on enterprise plan
		copyBasic := types.EnterprisePlan
		limit = &copyBasic
	} else if err != nil {
		return nil, err
	} else {
		limit = limitModel.ToProjectUsageType()
	}

	return limit, nil
}
