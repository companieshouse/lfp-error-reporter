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
	expectedTLFPFileNamePrefix string = "CHS-LFP-CARD-ERRORS-"
	expectedCSVFileSuffix      string = ".csv"
	reconciliationDate         string = "2019-01-01"
	timeFormatLayout           string = "2006-01-02"
)

func createMockService(cfg *config.Config, mockDao *dao.MockDAO) *ServiceImpl {

	return &ServiceImpl{
		Config: cfg,
		DAO:    mockDao,
	}
}

func TestUnitGetLFPCSV(t *testing.T) {

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

		Convey("Given lfp data is successfully fetched", func() {

			var lfps models.PenaltyList
			mockDao.EXPECT().GetLFPData(&reconciliationMetaData).Return(lfps, nil).Times(1)

			Convey("Then no errors are returned", func() {

				lfpsCSV, err := svc.GetLFPCSV(&reconciliationMetaData)
				So(err, ShouldBeNil)

				Convey("And a CSV is successfully constructed", func() {

					So(lfpsCSV, ShouldNotBeNil)
					So(lfpsCSV.FileName, ShouldEqual, expectedTLFPFileNamePrefix+reconciliationMetaData.StartTime.AddDate(0, 0, -1).Format(timeFormatLayout)+expectedCSVFileSuffix)
				})
			})
		})
	})

	Convey("Subject: Failure to retrieve lfp data", t, func() {

		mockDao := dao.NewMockDAO(mockCtrl)

		svc := createMockService(&cfg, mockDao)

		Convey("Given an error when fetching LFP data", func() {

			var lfps models.PenaltyList
			mockDao.EXPECT().GetLFPData(&reconciliationMetaData).Return(lfps, errors.New("failure to fetch lfp data")).Times(1)

			Convey("Then errors are returned", func() {

				lfpsCSV, err := svc.GetLFPCSV(&reconciliationMetaData)
				So(err, ShouldNotBeNil)

				Convey("And no CSV is constructed", func() {
					So(lfpsCSV.Data, ShouldBeNil)
				})
			})
		})
	})
}
