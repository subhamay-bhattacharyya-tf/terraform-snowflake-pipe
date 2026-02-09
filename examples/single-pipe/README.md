# Basic Example - Single Snowflake Pipe

This example demonstrates how to create a single Snowflake pipe using the `snowflake-pipe` module.

## Usage

```hcl
module "pipe" {
  source = "../../modules/snowflake-pipe"

  pipe_configs = {
    "my_pipe" = {
      database       = "MY_DATABASE"
      schema         = "MY_SCHEMA"
      name           = "MY_PIPE"
      copy_statement = "COPY INTO MY_DATABASE.MY_SCHEMA.MY_TABLE FROM @MY_DATABASE.MY_SCHEMA.MY_STAGE"
      auto_ingest    = false
      comment        = "My test pipe"
    }
  }
}
```

## Requirements

- Terraform >= 1.3.0
- Snowflake provider >= 0.87.0
- Existing database, schema, stage, and target table in Snowflake

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| pipe_configs | Map of pipe configuration objects | map(object) | yes |

## Outputs

| Name | Description |
|------|-------------|
| pipe_names | The names of the created pipes |
| pipe_fully_qualified_names | The fully qualified names of the pipes |
| pipe_notification_channels | The notification channels for the pipes |
| pipe_owners | The owners of the pipes |
