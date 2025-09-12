locals {

  additional_iam_policies_json = [data.aws_iam_policy_document.lfp_error_reporter_execution.json]

  # todo - should I move this to the pipeline?
  application_subnet_pattern = local.vault_secrets["APPLICATION_SUBNET_PATTERN"]

  lambda_env_vars = merge(local.vault_secrets, var.open_lambda_environment_variables)

  lambda_vpc_access_subnet_ids = data.aws_subnets.application.ids

  vault_secrets = data.vault_generic_secret.configuration.data
}
