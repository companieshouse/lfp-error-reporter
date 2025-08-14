// Package config defines the environment variables and flags
package config

import (
	"github.com/companieshouse/gofigure"
)

// Config holds configuration details required to execute the lambda.
type Config struct {
	PayableResourcesCollection string `env:"PPS_MONGODB_PAYABLE_RESOURCES_COLLECTION"       flag:"mongodb-payable-resources-collection"              flagDesc:"MongoDB collection for penalty payment data"`
	Database                   string `env:"PPS_MONGODB_DATABASE"                           flag:"mongodb-database"                                  flagDesc:"MongoDB database for penalty payment data"`
	MongoDBURL                 string `env:"MONGODB_URL"                                    flag:"mongodb-url"                                       flagDesc:"MongoDB server URL" json:"-"`
	SFTPServer                 string `env:"SFTP_SERVER"                                    flag:"sftp-server"                                       flagDesc:"Name of the SFTP server" json:"-"`
	SFTPPort                   string `env:"SFTP_PORT"                                      flag:"sftp-port"                                         flagDesc:"Port to connect to the SFTP server"`
	SFTPUserName               string `env:"SFTP_USERNAME"                                  flag:"sftp-username"                                     flagDesc:"Username of SFTP server" json:"-"`
	SFTPPassword               string `env:"SFTP_PASSWORD"                                  flag:"sftp-password"                                     flagDesc:"Password of SFTP server" json:"-"`
	SFTPFilePath               string `env:"SFTP_FILE_PATH"                                 flag:"sftp-file-path"                                    flagDesc:"File path on the SFTP server" json:"-"`
}

// Get returns configuration details marshalled into a Config struct
func Get() (*Config, error) {

	cfg := &Config{}

	err := gofigure.Gofigure(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
