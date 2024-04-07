//go:build ee
// +build ee

package repository

type EERepository interface {
	UserBilling() UserBillingRepository
	ProjectBilling() ProjectBillingRepository
}
