package integrations

import (
	"github.com/porter-dev/porter/api/types"
	"gorm.io/gorm"
)

// GithubAppInstallation is an instance of the porter github app
// we need to store account/installation id pairs in order to authenticate as the installation
type GithubAppInstallation struct {
	gorm.Model

	// Can belong to either a user or an organization
	AccountID int64 `json:"account_id" gorm:"unique"`

	// Installation ID (used for authentication)
	InstallationID int64 `json:"installation_id"`
}

func (r *GithubAppInstallation) ToGitInstallationType() *types.GitInstallation {
	return &types.GitInstallation{
		ID:             r.ID,
		AccountID:      r.AccountID,
		InstallationID: r.InstallationID,
	}
}
