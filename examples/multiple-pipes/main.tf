# Example: Multiple Snowflake Pipes
#
# This example demonstrates how to use the snowflake-pipe module
# to create multiple Snowflake pipes using a map of configurations.

module "pipes" {
  source = "../../modules/snowflake-pipe"

  pipe_configs = var.pipe_configs
}
