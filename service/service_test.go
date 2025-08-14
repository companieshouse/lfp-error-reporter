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
	expectedTPPSFileNamePrefix string = "CHS-PPS-CARD-ERRORS-"
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

func TestUnitGetPPSCSV(t *testing.T) {

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

			var ppss models.PenaltyList
			mockDao.EXPECT().GetPPSData(&reconciliationMetaData).Return(ppss, nil).Times(1)

			Convey("Then no errors are returned", func() {

				ppssCSV, err := svc.GetPPSCSV(&reconciliationMetaData)
				So(err, ShouldBeNil)

				Convey("And a CSV is successfully constructed", func() {

					So(ppssCSV, ShouldNotBeNil)
					So(ppssCSV.FileName, ShouldEqual, expectedTPPSFileNamePrefix+reconciliationMetaData.StartTime.AddDate(0, 0, -1).Format(timeFormatLayout)+expectedCSVFileSuffix)
				})
			})
		})
	})

	Convey("Subject: Failure to retrieve penalty payment data", t, func() {

		mockDao := dao.NewMockDAO(mockCtrl)

		svc := createMockService(&cfg, mockDao)

		Convey("Given an error when fetching Penalty payment data", func() {

			var ppss models.PenaltyList
			mockDao.EXPECT().GetPPSData(&reconciliationMetaData).Return(ppss, errors.New("failure to fetch penalty payment data")).Times(1)

			Convey("Then errors are returned", func() {

				ppssCSV, err := svc.GetPPSCSV(&reconciliationMetaData)
				So(err, ShouldNotBeNil)

				Convey("And no CSV is constructed", func() {
					So(ppssCSV.Data, ShouldBeNil)
				})
			})
		})
	})
}
