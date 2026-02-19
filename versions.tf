# -----------------------------------------------------------------------------
# Terraform Snowflake Pipe Module
# -----------------------------------------------------------------------------
# This module creates and manages Snowflake pipes using a map-based
# configuration. It supports creating single or multiple pipes with
# auto-ingest, AWS SNS integration, error integration, and storage
# integration settings in a single module call.
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
