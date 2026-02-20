# Multiple Pipes Example

This example demonstrates how to create multiple Snowflake pipes using the `snowflake-pipe` module with a map of configurations.

## Usage

```hcl
module "pipes" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-pipe"

  pipe_configs = {
    "orders_pipe" = {
      database       = "ANALYTICS_DB"
      schema         = "RAW"
      name           = "ORDERS_PIPE"
      copy_statement = "COPY INTO ANALYTICS_DB.RAW.ORDERS FROM @ANALYTICS_DB.RAW.S3_STAGE/orders/"
      auto_ingest    = false
      comment        = "Pipe for loading order data from S3"
    }
    "customers_pipe" = {
      database       = "ANALYTICS_DB"
      schema         = "RAW"
      name           = "CUSTOMERS_PIPE"
      copy_statement = "COPY INTO ANALYTICS_DB.RAW.CUSTOMERS FROM @ANALYTICS_DB.RAW.S3_STAGE/customers/"
      auto_ingest    = false
      comment        = "Pipe for loading customer data from S3"
    }
    "events_pipe" = {
      database          = "ANALYTICS_DB"
      schema            = "RAW"
      name              = "EVENTS_PIPE"
      copy_statement    = "COPY INTO ANALYTICS_DB.RAW.EVENTS FROM @ANALYTICS_DB.RAW.S3_STAGE/events/"
      auto_ingest       = true
      aws_sns_topic_arn = "arn:aws:sns:us-east-1:123456789012:snowflake-events"
      comment           = "Pipe for auto-ingesting event data from S3"
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
- Existing external stages pointing to cloud storage (S3, Azure Blob, GCS)
- Target tables with schemas matching the staged files
- For auto-ingest pipes: AWS SNS topic or storage integration configured (optional)

## Provider Configuration

The Snowflake provider requires `preview_features_enabled` to use the pipe resource:

```hcl
provider "snowflake" {
  # ... other configuration ...
  preview_features_enabled = ["snowflake_pipe_resource"]
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| pipe_configs | Map of pipe configuration objects | `map(object)` | yes |
| snowflake_organization_name | Snowflake organization name | `string` | yes |
| snowflake_account_name | Snowflake account name | `string` | yes |
| snowflake_user | Snowflake username | `string` | yes |
| snowflake_role | Snowflake role | `string` | no |
| snowflake_private_key | Snowflake private key for key-pair authentication | `string` | yes |

## pipe_configs Object Properties

| Property | Type | Default | Description |
|----------|------|---------|-------------|
| database | string | - | Database where the pipe resides (required) |
| schema | string | - | Schema where the pipe resides (required) |
| name | string | - | Pipe identifier (required) |
| copy_statement | string | - | COPY statement for the pipe (required) |
| auto_ingest | bool | false | Enable auto-ingest for the pipe |
| aws_sns_topic_arn | string | null | AWS SNS topic ARN for auto-ingest notifications |
| error_integration | string | null | Name of the notification integration for error notifications |
| integration | string | null | Name of the storage integration for auto-ingest |
| comment | string | null | Description of the pipe |

## Outputs

| Name | Description |
|------|-------------|
| pipe_names | The names of the created pipes |
| pipe_fully_qualified_names | The fully qualified names of the pipes |
| pipe_notification_channels | The notification channels for the pipes |
| pipe_owners | The owners of the pipes |
| pipes | All pipe resources |

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

## Auto-Ingest Configuration

For pipes with `auto_ingest = true`, you can optionally configure event notifications from your cloud storage:

### AWS S3
1. Create an SNS topic
2. Configure S3 bucket event notifications to publish to the SNS topic
3. Grant Snowflake access to the SNS topic
4. Use the SNS topic ARN in `aws_sns_topic_arn`

### Using Storage Integration
Alternatively, use a storage integration:
```hcl
"auto_ingest_pipe" = {
  database       = "MY_DB"
  schema         = "MY_SCHEMA"
  name           = "AUTO_PIPE"
  copy_statement = "COPY INTO MY_DB.MY_SCHEMA.MY_TABLE FROM @MY_DB.MY_SCHEMA.MY_STAGE"
  auto_ingest    = true
  integration    = "MY_STORAGE_INTEGRATION"
  comment        = "Auto-ingest pipe using storage integration"
}
```
