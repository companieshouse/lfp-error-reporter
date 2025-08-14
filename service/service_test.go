package service

import (
	"errors"
	"testing"
	"time"

	"github.com/companieshouse/lfp-error-reporter/config"
	"github.com/companieshouse/lfp-error-reporter/dao"
	"github.com/companieshouse/lfp-error-reporter/models"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	expectedFileNamePrefix string = "CHS-PENALTY-PAYMENT-E5-ERRORS-"
	expectedCSVFileSuffix  string = ".csv"
	reconciliationDate     string = "2019-01-01"
	timeFormatLayout       string = "2006-01-02"
)

func createMockService(cfg *config.Config, mockDao *dao.MockDAO) *ServiceImpl {

	return &ServiceImpl{
		Config: cfg,
		DAO:    mockDao,
	}
}

func TestUnitGetFailingPaymentCSV(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	startTime, _ := time.Parse(timeFormatLayout, reconciliationDate)

	cfg := config.Config{}
	reconciliationMetaData := models.ReconciliationMetaData{
		ReconciliationDate: reconciliationDate,
		StartTime:          startTime,
	}

	Convey("Subject: Success", t, func() {

		mockDao := dao.NewMockDAO(mockCtrl)

		svc := createMockService(&cfg, mockDao)

		Convey("Given penalty payment data is successfully fetched", func() {

			var penaltyPayments models.PenaltyList
			mockDao.EXPECT().GetPenaltyPaymentData(&reconciliationMetaData).Return(penaltyPayments, nil).Times(1)

			Convey("Then no errors are returned", func() {

				failingPaymentCSV, err := svc.GetFailingPaymentCSV(&reconciliationMetaData)
				So(err, ShouldBeNil)

				Convey("And a CSV is successfully constructed", func() {

					So(failingPaymentCSV, ShouldNotBeNil)
					So(failingPaymentCSV.FileName, ShouldEqual, expectedFileNamePrefix+reconciliationMetaData.StartTime.AddDate(0, 0, -1).Format(timeFormatLayout)+expectedCSVFileSuffix)
				})
			})
		})
	})

	Convey("Subject: Failure to retrieve penalty payment data", t, func() {

		mockDao := dao.NewMockDAO(mockCtrl)

		svc := createMockService(&cfg, mockDao)

		Convey("Given an error when fetching Penalty payment data", func() {

			var penaltyPayments models.PenaltyList
			mockDao.EXPECT().GetPenaltyPaymentData(&reconciliationMetaData).Return(penaltyPayments, errors.New("failure to fetch penalty payment data")).Times(1)

			Convey("Then errors are returned", func() {

				failingPaymentCSV, err := svc.GetFailingPaymentCSV(&reconciliationMetaData)
				So(err, ShouldNotBeNil)

				Convey("And no CSV is constructed", func() {
					So(failingPaymentCSV.Data, ShouldBeNil)
				})
			})
		})
	})
}
