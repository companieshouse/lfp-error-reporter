resource "aws_cloudwatch_event_rule" "lfp_error_reporter" {
  name        = "${var.service}-${var.environment}"
  description = "Call lfp error reporter lambda"
  schedule_expression =var.cron_schedule
}

resource "aws_cloudwatch_event_target" "call_lfp_error_reporter_lambda" {
    rule = aws_cloudwatch_event_rule.lfp_error_reporter.name
    target_id = "${var.service}-${var.environment}"
    arn = var.lambda_arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_lfp_error_reporter" {
    statement_id = "AllowLambdaExecutionFromCloudWatch"
    action = "lambda:InvokeFunction"
    function_name = "${var.service}-${var.environment}"
    principal = "events.amazonaws.com"
    source_arn = aws_cloudwatch_event_rule.lfp_error_reporter.arn
}
