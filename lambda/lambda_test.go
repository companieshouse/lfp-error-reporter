package lambda

import (
	"errors"
	"testing"

	"github.com/companieshouse/lfp-error-reporter/config"
	"github.com/companieshouse/lfp-error-reporter/filetransfer"
	"github.com/companieshouse/lfp-error-reporter/models"
	"github.com/companieshouse/lfp-error-reporter/service"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func createMockLambda(cfg *config.Config, mockService *service.MockService, mockFileTransfer *filetransfer.MockFileTransfer) *Lambda {

	return &Lambda{
		Config:       cfg,
		Service:      mockService,
		FileTransfer: mockFileTransfer,
	}
}

func TestUnitExecute(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cfg := config.Config{}
	reconciliationMetaData := models.ReconciliationMetaData{}

	Convey("Subject: Success", t, func() {

		mockService := service.NewMockService(mockCtrl)
		mockFileTransfer := filetransfer.NewMockFileTransfer(mockCtrl)

		lambda := createMockLambda(&cfg, mockService, mockFileTransfer)

		Convey("Given a penalty payment error CSV is constructed successfully", func() {

			var PPSCSV models.CSV
			mockService.EXPECT().GetPPSCSV(&reconciliationMetaData).Return(PPSCSV, nil).Times(1)

			Convey("And the CSV is uploaded successfully", func() {

				csvs := []models.CSV{PPSCSV}
				mockFileTransfer.EXPECT().UploadCSVFiles(csvs).Return(nil).Times(1)

				Convey("Then the request is successful", func() {

					err := lambda.Execute(&reconciliationMetaData)
					So(err, ShouldBeNil)
				})
			})
		})
	})

	Convey("Subject: Failure to construct penalty payment CSV", t, func() {

		mockService := service.NewMockService(mockCtrl)
		mockFileTransfer := filetransfer.NewMockFileTransfer(mockCtrl)

		lambda := createMockLambda(&cfg, mockService, mockFileTransfer)

		Convey("Given a failure when constructing a transactions CSV", func() {

			var ppsCSV models.CSV
			mockService.EXPECT().GetPPSCSV(&reconciliationMetaData).Return(ppsCSV, errors.New("failed to construct penalty payment CSV")).Times(1)

			Convey("And no CSV's are uploaded", func() {

				mockFileTransfer.EXPECT().UploadCSVFiles(gomock.Any()).Times(0)

				Convey("And the request is unsuccessful", func() {

					err := lambda.Execute(&reconciliationMetaData)
					So(err, ShouldNotBeNil)
				})
			})
		})
	})

	Convey("Subject: Failure to upload CSV's", t, func() {

		mockService := service.NewMockService(mockCtrl)
		mockFileTransfer := filetransfer.NewMockFileTransfer(mockCtrl)

		lambda := createMockLambda(&cfg, mockService, mockFileTransfer)

		Convey("Given a penalty payment CSV is constructed successfully", func() {

			var ppsCSV models.CSV
			mockService.EXPECT().GetPPSCSV(&reconciliationMetaData).Return(ppsCSV, nil).Times(1)

			Convey("But the CSV is not uploaded successfully", func() {

				csvs := []models.CSV{ppsCSV}
				mockFileTransfer.EXPECT().UploadCSVFiles(csvs).Return(errors.New("failure to upload CSV")).Times(1)

				Convey("Then the request is unsuccessful", func() {

					err := lambda.Execute(&reconciliationMetaData)
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}
