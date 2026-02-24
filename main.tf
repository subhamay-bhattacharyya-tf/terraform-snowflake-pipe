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

# Flatten grants for all pipes into a single map for iteration
locals {
  pipe_grants = flatten([
    for pipe_key, pipe in var.pipe_configs : [
      for grant in pipe.grants : {
        pipe_key   = pipe_key
        role_name  = grant.role_name
        privileges = grant.privileges
      }
    ]
  ])

  pipe_grants_map = {
    for grant in local.pipe_grants :
    "${grant.pipe_key}_${grant.role_name}" => grant
  }
}

resource "snowflake_grant_privileges_to_account_role" "pipe_grants" {
  for_each = local.pipe_grants_map

  privileges        = each.value.privileges
  account_role_name = each.value.role_name
  on_schema_object {
    object_type = "PIPE"
    object_name = snowflake_pipe.this[each.value.pipe_key].fully_qualified_name
  }
}
