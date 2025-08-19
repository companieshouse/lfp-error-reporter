# Penalty Payment Error Reporter
A lambda to report errors from the penalty payment service when communicating with E5 Financial Management System. This service collates failed E5 payment actions from the payable resources collection and writes it to a CSV file on an SFTP server. More information can be found [here](https://companieshouse.atlassian.net/wiki/spaces/TEAM8/pages/2824667427/LFP+Error+Reporter)

> [!IMPORTANT]  
> Testing can only be carried out from CIDEV as an instance of AWS CloudWatch is required and is documented [here](https://companieshouse.atlassian.net/wiki/spaces/TEAM8/pages/2824667427/LFP+Error+Reporter#Testing-on-AWS).

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

Name                                             | Description                                                                                              | Examples
------------------------------------------------ |----------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------
PPS_MONGODB_PAYABLE_RESOURCES_COLLECTION         | The name of the collection within the error reporting database from which to fetch penalty payment data. | 'payable_resources'
PPS_MONGODB_DATABASE                             | The name of the database containing the collection from which to fetch penalty payment data.             | 'financial_penalties'
MONGODB_URL                                      | The Mongo database URL.                                                                                  | 'mongodb://<mongo_host>:27017
SFTP_SERVER                                      | The SFTP server host name.                                                                               | 
SFTP_PORT                                        | The port over which to connect to the SFTP server.                                                       | '22'
SFTP_USERNAME                                    | The username of the SFTP server credentials.                                                             | 
SFTP_PASSWORD                                    | The password of the SFTP server credentials.                                                             |
SFTP_FILE_PATH                                   | The file path, relative to the root of the SFTP server, to which to upload CSV files.                    | 'uploadPath' (will result in CV's uploaded to directory: ~/uploadPath)
