// Package dao provides access to the database.
package dao

import (
	"github.com/companieshouse/lfp-error-reporter/models"
)

// DAO provides access to the database
type DAO interface {
	GetPenaltyPaymentData(reconciliationMetaData *models.ReconciliationMetaData) (models.PenaltyList, error)
}
