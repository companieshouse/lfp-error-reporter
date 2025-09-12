variable "aws_region" {
  type        = string
  description = "AWS Region"
  default     = "eu-west-2"
}

variable "aws_profile" {
  type        = string
  description = "The AWS profile to use for deployment."
}

variable "handler" {
  type        = string
  default     = "lfp-error-reporter"
  description = "The entrypoint in the Lambda funtion."
}

variable "memory_megabytes" {
  type        = string
  default     = "320"
  description = "The amount of memory to allocate to the Lambda function"
}

variable "release_bucket_name" {
  type        = string
  default     = "release.ch.gov.uk"
  description = "The S3 release bucket location containing the function code. "
}

variable "release_artifact_key" {
  type        = string
  description = "The release artifact key for the Lambda function"
}

variable "runtime" {
  type        = string
  default     = "go1.x"
  description = "The Lambda function language / runtime."
}

variable "timeout_seconds" {
  type        = string
  default     = "6"
  description = "The amount of time the Lambda function has to run in seconds."
}

variable "environment" {
  type        = string
  description = "The name of the environment to deploy."
}

variable "service" {
  type        = string
  default     = "lfp-error-reporter"
  description = "The name of the service being deployed."
}

variable "cron_schedule" {
  type        = string
  description = "CloudWatch cron schedule expression for calling the Lambda function."
}

variable open_lambda_environment_variables {
  type        = map(string)
  description = "Lambda environment variables that do not require encryption."
  default     = {}
}

variable "vpc_name" {
  type        = string
  description = "The VPC in which to create resources"
  default     = "Test & Development"
}

variable "lambda_logs_retention_days" {
  type        = number
  description = "The number of days to retain Lambda logs in CloudWatch"
  default     = 7
}

variable "application_subnet_pattern_key" {
  type        = string
  description = "The vault key for the application subnet pattern"
}

variable "lambda_runtime" {
  type        = string
  description = "The lambda runtime to run the application"
  default     = "go1.x"
}
