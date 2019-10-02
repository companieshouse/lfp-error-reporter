# lfp-error-reporter
A lambda to report errors from the lfp service when communicating with e5. This service collates failed lfp payments from the payable resources collection and writes it to CSV's on an SFTP server.

### The Lambda function
This API runs on AWS Lambda. Release and deployment of the lambda is similar to that of other trunk based services. 
After performing relevant unit tests, a github release task is run and the zip is uploaded to the S3 release bucket. 
Terraform is then run to deploy the new version. 

### Terraform deployment
All dependent AWS resources are provisioned by Terraform and deployed from a concourse pipeline.
The pipeline is capable of deploying everything so manual deployment should not be necessary. For
instructions on Terraform provisioning, see [here](/terraform/README.md).

### Environment Variables
Environment variables required to execute the lambda:

Name                                             | Description                                                                                                   | Examples
------------------------------------------------ | --------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------
MONGODB_LFP_ERR_REP_TRANSACTIONS_COLLECTION      | The name of the collection within the error reporting database from which to fetch lfp data.                  | 'payable_resources'
MONGODB_LFP_ERR_REP_DATABASE                     | The name of the database containing the collection from which to fetch lfp data.                              | 'late_filing_penalties'
MONGODB_URL                                      | The Mongo database URL.                                                                                       | 'mongodb://<mongo_host>:27017
SFTP_SERVER                                      | The SFTP server host name.                                                                                    | 
SFTP_PORT                                        | The port over which to connect to the SFTP server.                                                            | '22'
SFTP_USERNAME                                    | The username of the SFTP server credentials.                                                                  | 
SFTP_PASSWORD                                    | The password of the SFTP server credentials.                                                                  |
SFTP_FILE_PATH                                   | The file path, relative to the root of the SFTP server, to which to upload CSV files.                         | 'uploadPath' (will result is CV's uploaded to directory: ~/uploadPath)
