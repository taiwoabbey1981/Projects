package gorm

import (
	"context"

	"github.com/porter-dev/porter/internal/models"
	"github.com/porter-dev/porter/internal/repository"
	"gorm.io/gorm"
)

// PorterAppRepository uses gorm.DB for querying the database
type PorterAppRepository struct {
	db *gorm.DB
}

// NewPorterAppRepository returns a PorterAppRepository which uses
// gorm.DB for querying the database
func NewPorterAppRepository(db *gorm.DB) repository.PorterAppRepository {
	return &PorterAppRepository{db}
}

func (repo *PorterAppRepository) CreatePorterApp(a *models.PorterApp) (*models.PorterApp, error) {
	if err := repo.db.Create(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

func (repo *PorterAppRepository) ListPorterAppByClusterID(clusterID uint) ([]*models.PorterApp, error) {
	apps := []*models.PorterApp{}

	if err := repo.db.Where("cluster_id = ?", clusterID).Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}

// ReadPorterAppByID returns a PorterApp by its ID
func (repo *PorterAppRepository) ReadPorterAppByID(ctx context.Context, id uint) (*models.PorterApp, error) {
	app := &models.PorterApp{}

	if err := repo.db.Where("id = ?", id).Limit(1).Find(&app).Error; err != nil {
		return nil, err
	}

	return app, nil
}

func (repo *PorterAppRepository) ReadPorterAppByName(clusterID uint, name string) (*models.PorterApp, error) {
	app := &models.PorterApp{}

	if err := repo.db.Where("cluster_id = ? AND name = ?", clusterID, name).Limit(1).Find(&app).Error; err != nil {
		return nil, err
	}

	return app, nil
}

// ReadPorterAppsByProjectIDAndName returns a list of PorterApps by project ID and name. Multiple apps can have the same name and project id
// if they are in different clusters.
func (repo *PorterAppRepository) ReadPorterAppsByProjectIDAndName(projectID uint, name string) ([]*models.PorterApp, error) {
	apps := []*models.PorterApp{}

	if err := repo.db.Where("project_id = ? AND name = ?", projectID, name).Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}

func (repo *PorterAppRepository) UpdatePorterApp(app *models.PorterApp) (*models.PorterApp, error) {
	if err := repo.db.Save(app).Error; err != nil {
		return nil, err
	}

	return app, nil
}

func (repo *PorterAppRepository) DeletePorterApp(app *models.PorterApp) (*models.PorterApp, error) {
	if err := repo.db.Delete(&app).Error; err != nil {
		return nil, err
	}

	return app, nil
}
