package service

import (
	"reflect"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/lfp-error-reporter/config"
	"github.com/companieshouse/lfp-error-reporter/dao"
	"github.com/companieshouse/lfp-error-reporter/models"
)

const lfpFileNamePrefix = "CHS-LFP-CARD-ERRORS-"
const csvFileSuffix = ".csv"

// Service provides functionality by which to fetch lfp error CSV's
type Service interface {
	GetLFPCSV(reconciliationMetaData *models.ReconciliationMetaData) (models.CSV, error)
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

// GetLFPCSV retrieves lfp data and constructs a CSV
func (s *ServiceImpl) GetLFPCSV(reconciliationMetaData *models.ReconciliationMetaData) (models.CSV, error) {

	var csv models.CSV

	log.Info("Fetching lfp data.")

	penalties, err := s.DAO.GetLFPData(reconciliationMetaData)
	if err != nil {
		return csv, err
	}

	log.Info("Successfully retrieved lfp data.")
	log.Trace("LFP data", log.Data{"lfp_data": penalties})

	// Convert LFP data to format required for CSV
	var penaltyErrorDataList models.PenaltyErrorDataList
	for _, p := range penalties.Penalties {
		keys := reflect.ValueOf(p.Data.Transactions).MapKeys()
		penaltyErrorData := models.PenaltyErrorData{
			TransactionDate: p.Data.CreatedAt,
			PUON:            p.Reference,
			CompanyNumber:   p.CompanyNumber,
			MadeUpDate:      p.Data.Transactions[keys[0].String()].MadeUpDate,
			Value:           p.Data.Transactions[keys[0].String()].Amount,
		}
		penaltyErrorDataList.Penalties = append(penaltyErrorDataList.Penalties, penaltyErrorData)
	}

	csv = constructCSV(penaltyErrorDataList, lfpFileNamePrefix, reconciliationMetaData)

	return csv, nil
}

// constructCSV marshals CSVable data into a CSV, accompanied by a file name
func constructCSV(data models.CSVable, fileNamePrefix string, reconciliationMetaData *models.ReconciliationMetaData) models.CSV {

	return models.CSV{
		Data:     data,
		FileName: fileNamePrefix + reconciliationMetaData.ReconciliationDate + csvFileSuffix,
	}
}
