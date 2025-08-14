// Package lambda contains the core lambda executable code
package lambda

import (
	"time"

	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/lfp-error-reporter/config"
	"github.com/companieshouse/lfp-error-reporter/filetransfer"
	"github.com/companieshouse/lfp-error-reporter/models"
	"github.com/companieshouse/lfp-error-reporter/service"
)

const dateFormat = "2006-01-02"

// Lambda provides AWS lambda execution functionality
type Lambda struct {
	Config       *config.Config
	Service      service.Service
	FileTransfer filetransfer.FileTransfer
}

// New returns a new Lambda using the provided configs
func New(cfg *config.Config) *Lambda {

	return &Lambda{
		Config:       cfg,
		Service:      service.New(cfg),
		FileTransfer: filetransfer.New(cfg),
	}
}

// Execute handles lambda execution
func (lambda *Lambda) Execute(reconciliationMetaData *models.ReconciliationMetaData) error {
	if reconciliationMetaData.ReconciliationDate == "" {

		reconciliationDateTime := time.Now()
		reconciliationMetaData.ReconciliationDate = reconciliationDateTime.Format(dateFormat)

		startTime := reconciliationDateTime.Truncate(24 * time.Hour)
		reconciliationMetaData.StartTime = startTime
		reconciliationMetaData.EndTime = startTime.Add(24 * time.Hour)
	} else {

		startTime, err := time.Parse(dateFormat, reconciliationMetaData.ReconciliationDate)
		if err != nil {
			log.Error(err)
			return err
		}

		reconciliationMetaData.StartTime = startTime
		reconciliationMetaData.EndTime = startTime.Add(24 * time.Hour)
	}

	log.Info("Penalty payment error reporting lambda executing. Getting penalties with e5 errors for date: " + reconciliationMetaData.ReconciliationDate + ". Creating penalty payment error CSV.")

	failingPaymentsCSV, err := lambda.Service.GetFailingPaymentCSV(reconciliationMetaData)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("Penalty payment error CSV constructed.")
	log.Trace("Penalty payment error CSV", log.Data{"failing_payment_csv": failingPaymentsCSV})

	err = lambda.FileTransfer.UploadCSVFiles([]models.CSV{failingPaymentsCSV})
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("CSV's successfully uploaded. Lambda execution finished.")

	return nil
}
