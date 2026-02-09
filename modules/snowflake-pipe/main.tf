# Snowflake Pipe Resource
# Creates and manages one or more Snowflake pipes based on the pipe_configs map


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
