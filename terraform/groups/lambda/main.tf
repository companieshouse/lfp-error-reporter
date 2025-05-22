terraform {
  required_version = ">= 1.3.0, < 2.0.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.72.0, < 5.73"
    }
    vault = {
      source  = "hashicorp/vault"
      version = ">= 3.18.0, < 5.0"
    }
  }

  backend "s3" {
    encrypt = true
  }
}

provider "aws" {
  region = var.aws_region
}

# PH: we should not be using remote state - this is not best practice
# (see my note here: https://companieshouse.atlassian.net/wiki/spaces/~642153ecf1b529dfa98e4afc/pages/5150343173/Terraform+Lambda+Module#PR-Review)
data "terraform_remote_state" "network_remote_state" {
  backend = "s3"
  config = {
    bucket = var.remote_state_bucket
    key    = var.remote_state_key
    region = var.aws_region
  }
}

# PH: related to the above, get rid of this and
# PH: this relates back to the last convo I had with JW (see locals.tf in efs-doc-proc)
# PH: basically, we can filter the necessary resources by tag:Name etc and use those
locals {
  test_and_development_vpc_id     = data.terraform_remote_state.network_remote_state.outputs.vpc_id
  test_and_development_subnet_ids = split(",", data.terraform_remote_state.network_remote_state.outputs.application_ids)
}

module "lambda" {
  source                            = "./module-lambda"
  service                           = var.service
  handler                           = var.handler
  memory_megabytes                  = var.memory_megabytes
  runtime                           = var.runtime
  timeout_seconds                   = var.timeout_seconds
  release_version                   = var.release_version
  release_bucket_name               = var.release_bucket_name
  execution_role                    = module.lambda-roles.execution_role
  open_lambda_environment_variables = var.open_lambda_environment_variables
  aws_profile                       = var.aws_profile
  subnet_ids                        = local.test_and_development_subnet_ids
  security_group_ids                = [module.security-group.lambda_into_vpc_id]
  environment                       = var.environment
}

module "lambda-roles" {
  source      = "./module-lambda-roles"
  service     = var.service
  environment = var.environment
}

module "security-group" {
  source      = "./module-security-group"
  vpc_id      = local.test_and_development_vpc_id
  environment = var.environment
  service     = var.service
}

module "cloud-watch" {
  source        = "./module-cloud-watch"
  service       = var.service
  lambda_arn    = module.lambda.lambda_arn
  environment   = var.environment
  cron_schedule = var.cron_schedule
}
