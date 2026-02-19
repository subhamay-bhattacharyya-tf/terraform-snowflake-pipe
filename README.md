# Terraform Snowflake Module - Pipe

![Release](https://github.com/subhamay-bhattacharyya-tf/terraform-snowflake-pipe/actions/workflows/ci.yaml/badge.svg)&nbsp;![Snowflake](https://img.shields.io/badge/Snowflake-29B5E8?logo=snowflake&logoColor=white)&nbsp;![Commit Activity](https://img.shields.io/github/commit-activity/t/subhamay-bhattacharyya-tf/terraform-snowflake-pipe)&nbsp;![Last Commit](https://img.shields.io/github/last-commit/subhamay-bhattacharyya-tf/terraform-snowflake-pipe)&nbsp;![Release Date](https://img.shields.io/github/release-date/subhamay-bhattacharyya-tf/terraform-snowflake-pipe)&nbsp;![Repo Size](https://img.shields.io/github/repo-size/subhamay-bhattacharyya-tf/terraform-snowflake-pipe)&nbsp;![File Count](https://img.shields.io/github/directory-file-count/subhamay-bhattacharyya-tf/terraform-snowflake-pipe)&nbsp;![Issues](https://img.shields.io/github/issues/subhamay-bhattacharyya-tf/terraform-snowflake-pipe)&nbsp;![Top Language](https://img.shields.io/github/languages/top/subhamay-bhattacharyya-tf/terraform-snowflake-pipe)&nbsp;![Custom Endpoint](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/bsubhamay/a424dacc7e3eb7a15caadd293998bc24/raw/terraform-snowflake-pipe.json?)

A Terraform module for creating and managing Snowflake pipes using a map of configuration objects. Supports creating single or multiple pipes with a single module call.

## Features

- Map-based configuration for creating single or multiple pipes
- Built-in input validation with descriptive error messages
- Sensible defaults for optional properties
- Outputs keyed by pipe identifier for easy reference
- Support for auto-ingest with AWS SNS or storage integration
- Support for error integration configuration

## Usage

### Single Pipe

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

### Multiple Pipes

```hcl
locals {
  pipes = {
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

module "pipes" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-pipe"

  pipe_configs = local.pipes
}
```

### Auto-Ingest with Storage Integration

```hcl
module "pipe" {
  source = "github.com/subhamay-bhattacharyya-tf/terraform-snowflake-pipe"

  pipe_configs = {
    "auto_ingest_pipe" = {
      database       = "MY_DATABASE"
      schema         = "MY_SCHEMA"
      name           = "AUTO_INGEST_PIPE"
      copy_statement = "COPY INTO MY_DATABASE.MY_SCHEMA.MY_TABLE FROM @MY_DATABASE.MY_SCHEMA.MY_STAGE"
      auto_ingest    = true
      integration    = "MY_STORAGE_INTEGRATION"
      comment        = "Auto-ingest pipe using storage integration"
    }
  }
}
```

## Examples

- [Single Pipe](examples/single-pipe) - Create a single pipe
- [Multiple Pipes](examples/multiple-pipes) - Create multiple pipes

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.3.0 |
| snowflake | >= 0.87.0 |

## Providers

| Name | Version |
|------|---------|
| snowflake | >= 0.87.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|----------|
| pipe_configs | Map of configuration objects for Snowflake pipes | `map(object)` | `{}` | no |

### pipe_configs Object Properties

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
| pipe_names | Map of pipe names keyed by identifier |
| pipe_fully_qualified_names | Map of fully qualified pipe names |
| pipe_notification_channels | Map of notification channels for SNS integration |
| pipe_owners | Map of pipe owners |
| pipes | All pipe resources |

## Validation

The module validates inputs and provides descriptive error messages for:

- Empty database name
- Empty schema name
- Empty pipe name
- Empty copy statement
- Auto-ingest enabled without aws_sns_topic_arn or integration

## Testing

The module includes Terratest-based integration tests:

```bash
cd test
go mod tidy
go test -v -timeout 30m
```

Required environment variables for testing:
- `SNOWFLAKE_ORGANIZATION_NAME` - Snowflake organization name
- `SNOWFLAKE_ACCOUNT_NAME` - Snowflake account name
- `SNOWFLAKE_USER` - Snowflake username
- `SNOWFLAKE_ROLE` - Snowflake role (e.g., "SYSADMIN")
- `SNOWFLAKE_PRIVATE_KEY` - Snowflake private key for key-pair authentication

## CI/CD Configuration

The CI workflow runs on:
- Push to `main`, `feature/**`, and `bug/**` branches (when `*.tf`, `examples/**`, or `test/**` changes)
- Pull requests to `main` (when `*.tf`, `examples/**`, or `test/**` changes)
- Manual workflow dispatch

The workflow includes:
- Terraform validation and format checking
- Examples validation
- Terratest integration tests (output displayed in GitHub Step Summary)
- Changelog generation (non-main branches)
- Semantic release (main branch only)

The CI workflow uses the following GitHub organization variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `TERRAFORM_VERSION` | Terraform version for CI jobs | `1.3.0` |
| `GO_VERSION` | Go version for Terratest | `1.21` |
| `SNOWFLAKE_ORGANIZATION_NAME` | Snowflake organization name | - |
| `SNOWFLAKE_ACCOUNT_NAME` | Snowflake account name | - |
| `SNOWFLAKE_USER` | Snowflake username | - |
| `SNOWFLAKE_ROLE` | Snowflake role (e.g., SYSADMIN) | - |

The following GitHub secrets are required for Terratest integration tests:

| Secret | Description | Required |
|--------|-------------|----------|
| `SNOWFLAKE_PRIVATE_KEY` | Snowflake private key for key-pair authentication | Yes |

## License

MIT License - See [LICENSE](LICENSE) for details.
