# -----------------------------------------------------------------------------
# Terraform Snowflake Pipe Module - Multiple Pipes Example
# -----------------------------------------------------------------------------
# This example demonstrates how to use the snowflake-pipe module
# to create multiple Snowflake pipes using a map of configurations.
# -----------------------------------------------------------------------------

terraform {
  required_version = ">= 1.3.0"

  required_providers {
    snowflake = {
      source  = "snowflakedb/snowflake"
      version = ">= 1.0.0"
    }
  }
}

provider "snowflake" {
  organization_name = var.snowflake_organization_name
  account_name      = var.snowflake_account_name
  user              = var.snowflake_user
  role              = var.snowflake_role
  authenticator     = "SNOWFLAKE_JWT"
  private_key       = var.snowflake_private_key

  preview_features_enabled = ["snowflake_pipe_resource"]
}
