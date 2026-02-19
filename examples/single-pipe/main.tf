# -----------------------------------------------------------------------------
# Terraform Snowflake Pipe Module - Single Pipe Example
# -----------------------------------------------------------------------------
# This example demonstrates how to use the snowflake-pipe module
# to create a single Snowflake pipe.
# -----------------------------------------------------------------------------

module "pipe" {
  source = "../.."

  pipe_configs = var.pipe_configs
}
