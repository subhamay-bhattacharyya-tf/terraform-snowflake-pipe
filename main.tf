# -----------------------------------------------------------------------------
# Terraform Snowflake Pipe Module
# -----------------------------------------------------------------------------
# This module creates and manages Snowflake pipes using a map-based
# configuration. It supports creating single or multiple pipes with
# auto-ingest, AWS SNS integration, error integration, and storage
# integration settings in a single module call.
# -----------------------------------------------------------------------------

resource "snowflake_pipe" "this" {
  for_each = var.pipe_configs

  database       = each.value.database
  schema         = each.value.schema
  name           = each.value.name
  copy_statement = each.value.copy_statement

  auto_ingest       = each.value.auto_ingest
  aws_sns_topic_arn = each.value.aws_sns_topic_arn
  error_integration = each.value.error_integration
  integration       = each.value.integration
  comment           = each.value.comment
}
