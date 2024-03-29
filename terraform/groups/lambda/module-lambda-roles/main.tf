# ------------------------------------------------------------------------------
# Policy Documents
# ------------------------------------------------------------------------------
data "aws_iam_policy_document" "lfp_error_reporter_trust" {
  statement {
    effect = "Allow"

    actions = [
      "sts:AssumeRole",
    ]

    principals {
      type = "Service"

      identifiers = [
        "lambda.amazonaws.com",
      ]
    }
  }
}

data "aws_iam_policy_document" "lfp_error_reporter_execution" {
  statement {
    effect = "Allow"

    actions = [
      "logs:DescribeQueries",
      "logs:GetLogRecord",
      "s3:PutAccountPublicAccessBlock",
      "logs:PutDestinationPolicy",
      "logs:StopQuery",
      "logs:TestMetricFilter",
      "logs:DeleteDestination",
      "logs:CreateLogGroup",
      "logs:GetLogDelivery",
      "logs:ListLogDeliveries",
      "logs:CreateLogDelivery",
      "logs:DeleteResourcePolicy",
      "logs:PutResourcePolicy",
      "logs:DescribeExportTasks",
      "s3:GetAccountPublicAccessBlock",
      "logs:GetQueryResults",
      "s3:ListAllMyBuckets",
      "logs:UpdateLogDelivery",
      "logs:CancelExportTask",
      "logs:DeleteLogDelivery",
      "s3:HeadBucket",
      "s3:GetObject",
      "logs:PutDestination",
      "logs:DescribeResourcePolicies",
      "logs:DescribeDestinations",
      "ec2:CreateNetworkInterface",
      "ec2:DescribeNetworkInterfaces",
      "ec2:DeleteNetworkInterface"
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

# ------------------------------------------------------------------------------
# Roles
# ------------------------------------------------------------------------------
resource "aws_iam_role" "lfp_error_reporter_execution" {
  name               = "${var.service}-execution-${var.environment}"
  assume_role_policy = data.aws_iam_policy_document.lfp_error_reporter_trust.json
}

# ------------------------------------------------------------------------------
# Role Policies
# ------------------------------------------------------------------------------
resource "aws_iam_role_policy" "lfp_error_reporter_execution" {
  name   = "lfp_error_reporter_execution"
  role   = aws_iam_role.lfp_error_reporter_execution.id
  policy = data.aws_iam_policy_document.lfp_error_reporter_execution.json
}

output "execution_role" {
  value = aws_iam_role.lfp_error_reporter_execution.arn
}
