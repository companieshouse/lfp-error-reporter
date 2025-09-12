data "vault_generic_secret" "configuration" {
  path = "applications/${var.aws_profile}/${var.environment}/${var.service}/lambda_environment_variables"
}

data "aws_vpc" "vpc" {
  filter {
    name   = "tag:Name"
    values = [var.vpc_name]
  }
}

data "aws_subnets" "application" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.vpc.id]
  }

  filter {
    name   = "tag:Name"
    values = [local.application_subnet_pattern]
  }
}

data "aws_iam_policy_document" "lfp_error_reporter_execution" {
  statement {
    effect = "Allow"

    actions = [
      "s3:PutAccountPublicAccessBlock",
      "s3:GetAccountPublicAccessBlock",
      "s3:ListAllMyBuckets",
      "s3:HeadBucket",
      "s3:GetObject",
      "ec2:CreateNetworkInterface",
      "ec2:DescribeNetworkInterfaces",
      "ec2:DeleteNetworkInterface", # todo - delete network interfaces ???
      "logs:DescribeQueries",       # todo - do we even need all these logs: ...
      "logs:GetLogRecord",          #        the new module has its own log policy,
      "logs:PutDestinationPolicy",  #        but with much less permissions than this though ?
      "logs:StopQuery",
      "logs:TestMetricFilter",
      "logs:DeleteDestination",
      "logs:CreateLogGroup",        # todo - the new lambda module creates its own logging policy...
      "logs:GetLogDelivery",
      "logs:ListLogDeliveries",
      "logs:CreateLogDelivery",
      "logs:DeleteResourcePolicy",
      "logs:PutResourcePolicy",
      "logs:DescribeExportTasks",
      "logs:GetQueryResults",
      "logs:UpdateLogDelivery",
      "logs:CancelExportTask",
      "logs:DeleteLogDelivery",
      "logs:PutDestination",
      "logs:DescribeResourcePolicies",
      "logs:DescribeDestinations"
    ]

    resources = [
      "*"
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "s3:*",
      "logs:*"
    ]

    resources = [
      "arn:aws:logs:::log-group:/aws/lambda/${var.service}",
      "arn:aws:logs:*:*:log-group:*:*:*",
    ]
  }
}
