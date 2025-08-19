// Package service contains the logic that retrieves the penalty payment error data and constructs the CSV file
package service

import (
	"reflect"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/lfp-error-reporter/config"
	"github.com/companieshouse/lfp-error-reporter/dao"
	"github.com/companieshouse/lfp-error-reporter/models"
)

const (
	penaltyPaymentErrorFileNamePrefix string = "CHS-PENALTY-PAYMENT-E5-ERRORS-"
	csvFileSuffix                     string = ".csv"
	YYYYMMDD                          string = "2006-01-02"
)

// Service provides functionality by which to fetch penalty payment error CSV's
type Service interface {
	GetFailingPaymentCSV(reconciliationMetaData *models.ReconciliationMetaData) (models.CSV, error)
}

// ServiceImpl provides a concrete implementation of the Service interface
type ServiceImpl struct {
	Config *config.Config
	DAO    dao.DAO
}

// New returns a new, concrete implementation of the Service interface
func New(cfg *config.Config) *ServiceImpl {

	return &ServiceImpl{
		Config: cfg,
		DAO:    dao.New(cfg),
	}
}

// GetFailingPaymentCSV retrieves penalty payment data and constructs a CSV
func (s *ServiceImpl) GetFailingPaymentCSV(reconciliationMetaData *models.ReconciliationMetaData) (models.CSV, error) {

	var csv models.CSV

	log.Info("Fetching penalty payment data.")

	penalties, err := s.DAO.GetPenaltyPaymentData(reconciliationMetaData)
	if err != nil {
		return csv, err
	}

	log.Info("Successfully retrieved penalty payment data.")
	log.Trace("Penalty payment data", log.Data{"penalty_payment_data": penalties})

	// Convert Penalty payment data to format required for CSV
	var penaltyErrorDataList models.PenaltyErrorDataList
	for _, p := range penalties.Penalties {
		keys := reflect.ValueOf(p.Data.Transactions).MapKeys()
		penaltyErrorData := models.PenaltyErrorData{
			TransactionDate: p.Data.CreatedAt,
			PUON:            p.PayableRef,
			CompanyNumber:   p.CustomerCode,
			MadeUpDate:      p.Data.Transactions[keys[0].String()].MadeUpDate,
			Value:           p.Data.Transactions[keys[0].String()].Amount,
		}
		penaltyErrorDataList.Penalties = append(penaltyErrorDataList.Penalties, penaltyErrorData)
	}

	csv = constructCSV(penaltyErrorDataList, penaltyPaymentErrorFileNamePrefix, reconciliationMetaData)

	return csv, nil
}

// constructCSV marshals CSVable data into a CSV, accompanied by a file name
func constructCSV(data models.CSVable, fileNamePrefix string, reconciliationMetaData *models.ReconciliationMetaData) models.CSV {

	return models.CSV{
		Data:     data,
		FileName: fileNamePrefix + reconciliationMetaData.StartTime.AddDate(0, 0, -1).Format(YYYYMMDD) + csvFileSuffix,
	}
}
