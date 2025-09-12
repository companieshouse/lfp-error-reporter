locals {
  vault_secrets                = data.vault_generic_secret.configuration.data
  application_subnet_pattern   = local.vault_secrets[var.application_subnet_pattern_key]
  lambda_vpc_access_subnet_ids = data.aws_subnets.application.ids
  additional_iam_policies_json = [data.aws_iam_policy_document.lfp_error_reporter_execution.json]
}
