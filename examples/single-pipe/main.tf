# Example: Single Snowflake Pipe
#
# This example demonstrates how to use the snowflake-pipe module
# to create a single Snowflake pipe.

module "pipe" {
  source = "../../modules/snowflake-pipe"

  pipe_configs = var.pipe_configs
}
