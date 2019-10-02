package dao

import (
	"github.com/companieshouse/lfp-error-reporter/models"
)

// DAO provides access to the database
type DAO interface {
	GetLFPData(reconciliationMetaData *models.ReconciliationMetaData) (models.PenaltyList, error)
}
