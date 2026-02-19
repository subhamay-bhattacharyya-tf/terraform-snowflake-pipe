# Single Pipe Example

This example demonstrates how to create a single Snowflake pipe using the `snowflake-pipe` module.

## Usage

```hcl
module "pipe" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-pipe"

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

| Name | Version |
|------|---------|
| terraform | >= 1.3.0 |
| snowflake | >= 1.0.0 |

## Prerequisites

- Existing database and schema in Snowflake
- Existing external stage pointing to cloud storage (S3, Azure Blob, GCS)
- Target table with schema matching the staged files

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| pipe_configs | Map of pipe configuration objects | `map(object)` | yes |
| snowflake_organization_name | Snowflake organization name | `string` | yes |
| snowflake_account_name | Snowflake account name | `string` | yes |
| snowflake_user | Snowflake username | `string` | yes |
| snowflake_role | Snowflake role | `string` | no |
| snowflake_private_key | Snowflake private key for key-pair authentication | `string` | yes |

## Outputs

| Name | Description |
|------|-------------|
| pipe_names | The names of the created pipes |
| pipe_fully_qualified_names | The fully qualified names of the pipes |
| pipe_notification_channels | The notification channels for the pipes |
| pipe_owners | The owners of the pipes |

## Running This Example

```bash
# Set environment variables
export TF_VAR_snowflake_organization_name="your_org"
export TF_VAR_snowflake_account_name="your_account"
export TF_VAR_snowflake_user="your_user"
export TF_VAR_snowflake_role="SYSADMIN"
export TF_VAR_snowflake_private_key="$(cat ~/.snowflake/rsa_key.p8)"

# Initialize and apply
terraform init
terraform plan
terraform apply
```
