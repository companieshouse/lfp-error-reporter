module "lambda" {
  source = "git@github.com:companieshouse/terraform-modules.git//aws/lambda?ref=1.0.342"

  environment           = var.environment
  function_name         = var.service
  lambda_runtime        = var.lambda_runtime
  lambda_handler        = var.handler

  lambda_code_s3_bucket = var.release_bucket_name
  lambda_code_s3_key    = var.release_artifact_key

  lambda_memory_size         = var.memory_megabytes
  lambda_timeout_seconds     = var.timeout_seconds
  lambda_logs_retention_days = var.lambda_logs_retention_days

  lambda_env_vars = local.lambda_env_vars

  lambda_cloudwatch_event_rules = [
    {
      name                = "lfp-error-reporter-${var.environment}"
      description         = "Call penalty payment error reporter lambda"
      schedule_expression = var.cron_schedule
    }
  ]

#  additional_policies = local.additional_iam_policies_json

  lambda_sg_egress_rule = {
    from_port   = -1
    to_port     = -1
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lambda_vpc_access_subnet_ids = local.lambda_vpc_access_subnet_ids
  lambda_vpc_id                = data.aws_vpc.vpc.id
}
