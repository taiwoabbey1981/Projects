package repository

import (
	"context"

	"github.com/porter-dev/porter/internal/models"
)

// PorterAppRepository represents the set of queries on the PorterApp model
type PorterAppRepository interface {
	ReadPorterAppByID(ctx context.Context, id uint) (*models.PorterApp, error)
	ReadPorterAppByName(clusterID uint, name string) (*models.PorterApp, error)
	ReadPorterAppsByProjectIDAndName(projectID uint, name string) ([]*models.PorterApp, error)
	CreatePorterApp(app *models.PorterApp) (*models.PorterApp, error)
	ListPorterAppByClusterID(clusterID uint) ([]*models.PorterApp, error)
	UpdatePorterApp(app *models.PorterApp) (*models.PorterApp, error)
	DeletePorterApp(app *models.PorterApp) (*models.PorterApp, error)
}
