package config_test

import (
	"encoding/json"
	"os"

	"regexp"
	"testing"

	"github.com/companieshouse/lfp-error-reporter/config"

	. "github.com/smartystreets/goconvey/convey"
)

// key constants
const (
	PENALTYPAYMENTCOLLECTIONKEYCONST = `PPS_MONGODB_PAYABLE_RESOURCES_COLLECTION`
	DATABASEKEYCONST                 = `PPS_MONGODB_DATABASE`
	MONGODBURLKEYCONST               = `MONGODB_URL`
	SFTPSERVERKEYCONST               = `SFTP_SERVER`
	SFTPPORTKEYCONST                 = `SFTP_PORT`
	SFTPUSERNAMEKEYCONST             = `SFTP_USERNAME`
	SFTPPASSWORDKEYCONST             = `SFTP_PASSWORD`
	SFTPFILEPATHKEYCONST             = `SFTP_FILE_PATH`
)

// value constants
const (
	penaltyPaymentCollectionConst = `pps-collection`
	databaseConst                 = `oriel-db`
	mongoDbUrlConst               = `localhost:12344`
	sftpServerConst               = `sftp-server`
	sftpPortConst                 = `sftp-port`
	sftpUserNameConst             = `sftp-username`
	sftpPasswordConst             = `sftp-password`
	sftpFilePathConst             = `sftp-file-path`
	nameConst                     = `chs-log-test`
	namespaceConst                = `chs-namespace-test`
	configConst                   = `config`
)

func TestConfig(t *testing.T) {
	t.Parallel()
	os.Clearenv()
	var (
		err           error
		configuration *config.Config
		envVars       = map[string]string{
			PENALTYPAYMENTCOLLECTIONKEYCONST: penaltyPaymentCollectionConst,
			DATABASEKEYCONST:                 databaseConst,
			MONGODBURLKEYCONST:               mongoDbUrlConst,
			SFTPSERVERKEYCONST:               sftpServerConst,
			SFTPPORTKEYCONST:                 sftpPortConst,
			SFTPUSERNAMEKEYCONST:             sftpUserNameConst,
			SFTPPASSWORDKEYCONST:             sftpPasswordConst,
			SFTPFILEPATHKEYCONST:             sftpFilePathConst,
		}
		builtConfig = config.Config{
			PayableResourcesCollection: penaltyPaymentCollectionConst,
			Database:                   databaseConst,
			MongoDBURL:                 mongoDbUrlConst,
			SFTPServer:                 sftpServerConst,
			SFTPPort:                   sftpPortConst,
			SFTPUserName:               sftpUserNameConst,
			SFTPPassword:               sftpPasswordConst,
			SFTPFilePath:               sftpFilePathConst,
		}
		sftpUserNameRegex = regexp.MustCompile(sftpUserNameConst)
		sftpPasswordRegex = regexp.MustCompile(sftpPasswordConst)
		mongoDbUrlRegex   = regexp.MustCompile(mongoDbUrlConst)
		sftpServerRegex   = regexp.MustCompile(sftpServerConst)
		sftpFilePathRegex = regexp.MustCompile(sftpFilePathConst)
	)

	// set test env variables
	for varName, varValue := range envVars {
		os.Setenv(varName, varValue)
		defer os.Unsetenv(varName)
	}

	Convey("Given an environment with no environment variables set", t, func() {

		Convey("Then configuration should be nil", func() {
			So(configuration, ShouldBeNil)
		})

		Convey("When the config values are retrieved", func() {

			Convey("Then there should be no error returned, and values are as expected", func() {
				configuration, err = config.Get()

				So(err, ShouldBeNil)
				So(configuration, ShouldResemble, &builtConfig)
			})

			Convey("The generated JSON string from configuration should not contain sensitive data", func() {
				jsonByte, err := json.Marshal(builtConfig)

				So(err, ShouldBeNil)
				So(sftpUserNameRegex.Match(jsonByte), ShouldEqual, false)
				So(sftpPasswordRegex.Match(jsonByte), ShouldEqual, false)
				So(mongoDbUrlRegex.Match(jsonByte), ShouldEqual, false)
				So(sftpServerRegex.Match(jsonByte), ShouldEqual, false)
				So(sftpFilePathRegex.Match(jsonByte), ShouldEqual, false)
			})
		})
	})
}
