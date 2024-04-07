package test

import (
	"context"
	"errors"
	"strings"

	"github.com/porter-dev/porter/internal/models"
	"github.com/porter-dev/porter/internal/repository"
)

type PorterAppRepository struct {
	canQuery       bool
	failingMethods string
}

func NewPorterAppRepository(canQuery bool, failingMethods ...string) repository.PorterAppRepository {
	return &PorterAppRepository{canQuery, strings.Join(failingMethods, ",")}
}

func (repo *PorterAppRepository) ReadPorterAppByName(clusterID uint, name string) (*models.PorterApp, error) {
	return nil, errors.New("cannot write database")
}

// ReadPorterAppsByProjectIDAndName is a test method that is not implemented
func (repo *PorterAppRepository) ReadPorterAppsByProjectIDAndName(projectID uint, name string) ([]*models.PorterApp, error) {
	return nil, errors.New("cannot write database")
}

func (repo *PorterAppRepository) CreatePorterApp(app *models.PorterApp) (*models.PorterApp, error) {
	return nil, errors.New("cannot write database")
}

func (repo *PorterAppRepository) UpdatePorterApp(app *models.PorterApp) (*models.PorterApp, error) {
	return nil, errors.New("cannot write database")
}

// ListPorterAppByClusterID is a test method that is not implemented
func (repo *PorterAppRepository) ListPorterAppByClusterID(clusterID uint) ([]*models.PorterApp, error) {
	return nil, errors.New("cannot write database")
}

func (repo *PorterAppRepository) DeletePorterApp(app *models.PorterApp) (*models.PorterApp, error) {
	return nil, errors.New("cannot write database")
}

// ReadPorterAppByID is a test method that is not implemented
func (repo *PorterAppRepository) ReadPorterAppByID(ctx context.Context, id uint) (*models.PorterApp, error) {
	return nil, errors.New("cannot read database")
}
