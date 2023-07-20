// Package config defines the environment variables and flags
package config

import (
	"fmt"

	"github.com/companieshouse/gofigure"
)

// Config holds configuration details required to execute the lambda.
type Config struct {
	LFPCollection string `env:"MONGODB_LFP_ERR_REP_TRANSACTIONS_COLLECTION"       flag:"mongodb-lfp-err-rep-lfp-collection"                flagDesc:"MongoDB collection for lfp data"`
	Database      string `env:"MONGODB_LFP_ERR_REP_DATABASE"                      flag:"MongoDB database for lfp data"                     flagDesc:"MongoDB database for lfp data"`
	MongoDBURL    string `env:"MONGODB_URL"                                       flag:"mongodb-url"                                       flagDesc:"MongoDB server URL" json:"-"`
	SFTPServer    string `env:"SFTP_SERVER"                                       flag:"sftp-server"                                       flagDesc:"Name of the SFTP server" json:"-"`
	SFTPPort      string `env:"SFTP_PORT"                                         flag:"sftp-port"                                         flagDesc:"Port to connect to the SFTP server"`
	SFTPUserName  string `env:"SFTP_USERNAME"                                     flag:"sftp-username"                                     flagDesc:"Username of SFTP server" json:"-"`
	SFTPPassword  string `env:"SFTP_PASSWORD"                                     flag:"sftp-password"                                     flagDesc:"Password of SFTP server" json:"-"`
	SFTPFilePath  string `env:"SFTP_FILE_PATH"                                    flag:"sftp-file-path"                                    flagDesc:"File path on the SFTP server" json:"-"`
}

// Get returns configuration details marshalled into a Config struct
func Get() (*Config, error) {

	cfg := &Config{}

	err := gofigure.Gofigure(cfg)
	if err != nil {
		return nil, err
	}
	fmt.Println("a ...any========>", cfg)
	return cfg, nil
}
