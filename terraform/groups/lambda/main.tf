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
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lambda_vpc_access_subnet_ids = local.lambda_vpc_access_subnet_ids
  lambda_vpc_id                = data.aws_vpc.vpc.id
}

# temporary move blocks

moved {
  from = module.lambda.aws_lambda_function.lfp_error_reporter
  to   = module.lambda.aws_lambda_function.lambda
}

moved {
  from = module.security-group.aws_security_group.lfp_error_reporter
  to   = module.lambda.aws_security_group.lambda_sg
}

moved {
  from = module.lambda-roles.aws_iam_role.lfp_error_reporter_execution
  to   = module.lambda.aws_iam_role.lambda_execution
}

moved {
  from = module.cloud-watch.aws_cloudwatch_event_rule.lfp_error_reporter
  to   = module.lambda.aws_cloudwatch_event_rule.lambda_cloudwatch_event_rules["lfp-error-reporter-${var.environment}"]
}

moved {
  from = module.cloud-watch.aws_cloudwatch_event_target.call_lfp_error_reporter_lambda
  to   = module.lambda.aws_cloudwatch_event_target.lambda_target["lfp-error-reporter-${var.environment}"]
}

moved {
  from = module.cloud-watch.aws_lambda_permission.allow_cloudwatch_to_call_lfp_error_reporter
  to   = module.lambda.aws_lambda_permission.allow_cloudwatch_to_call_lambda["lfp-error-reporter-${var.environment}"]
}
