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

		Convey("Given a lfp CSV is constructed successfully", func() {

			var LFPCSV models.CSV
			mockService.EXPECT().GetLFPCSV(&reconciliationMetaData).Return(LFPCSV, nil).Times(1)

			Convey("And the CSV is uploaded successfully", func() {

				csvs := []models.CSV{LFPCSV}
				mockFileTransfer.EXPECT().UploadCSVFiles(csvs).Return(nil).Times(1)

				Convey("Then the request is successful", func() {

					err := lambda.Execute(&reconciliationMetaData)
					So(err, ShouldBeNil)
				})
			})
		})
	})

	Convey("Subject: Failure to construct lfp CSV", t, func() {

		mockService := service.NewMockService(mockCtrl)
		mockFileTransfer := filetransfer.NewMockFileTransfer(mockCtrl)

		lambda := createMockLambda(&cfg, mockService, mockFileTransfer)

		Convey("Given a failure when constructing a transactions CSV", func() {

			var lfpCSV models.CSV
			mockService.EXPECT().GetLFPCSV(&reconciliationMetaData).Return(lfpCSV, errors.New("failed to construct lfp CSV")).Times(1)

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

		Convey("Given a lfp CSV is constructed successfully", func() {

			var lfpCSV models.CSV
			mockService.EXPECT().GetLFPCSV(&reconciliationMetaData).Return(lfpCSV, nil).Times(1)

			Convey("But the CSV is not uploaded successfully", func() {

				csvs := []models.CSV{lfpCSV}
				mockFileTransfer.EXPECT().UploadCSVFiles(csvs).Return(errors.New("failure to upload CSV")).Times(1)

				Convey("Then the request is unsuccessful", func() {

					err := lambda.Execute(&reconciliationMetaData)
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}
