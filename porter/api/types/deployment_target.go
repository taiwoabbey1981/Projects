package types

import (
	"time"

	"github.com/google/uuid"
)

// DeploymentTarget is a struct that represents a unique cluster, namespace pair that a Porter app is deployed to.
type DeploymentTarget struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uint      `json:"project_id"`
	ClusterID uint      `json:"cluster_id"`

	Name         string    `json:"name"`
	Namespace    string    `json:"namespace"`
	IsPreview    bool      `json:"is_preview"`
	IsDefault    bool      `json:"is_default"`
	CreatedAtUTC time.Time `json:"created_at"`
	UpdatedAtUTC time.Time `json:"updated_at"`
}
