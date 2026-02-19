# -----------------------------------------------------------------------------
# Terraform Snowflake Pipe Module - Multiple Pipes Example
# -----------------------------------------------------------------------------
# This example demonstrates how to use the snowflake-pipe module
# to create multiple Snowflake pipes using a map of configurations.
# -----------------------------------------------------------------------------

module "pipes" {
  source = "../.."

  pipe_configs = var.pipe_configs
}
