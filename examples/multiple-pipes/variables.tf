variable "pipe_configs" {
  description = "Map of configuration objects for Snowflake pipes"
  type = map(object({
    database          = string
    schema            = string
    name              = string
    copy_statement    = string
    auto_ingest       = optional(bool, false)
    aws_sns_topic_arn = optional(string, null)
    error_integration = optional(string, null)
    integration       = optional(string, null)
    comment           = optional(string, null)
  }))
  default = {
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
      database       = "ANALYTICS_DB"
      schema         = "RAW"
      name           = "EVENTS_PIPE"
      copy_statement = "COPY INTO ANALYTICS_DB.RAW.EVENTS FROM @ANALYTICS_DB.RAW.S3_STAGE/events/"
      auto_ingest    = false
      comment        = "Pipe for loading event data from S3"
    }
  }
}

# Snowflake authentication variables
variable "snowflake_organization_name" {
  description = "Snowflake organization name"
  type        = string
  default     = null
}

variable "snowflake_account_name" {
  description = "Snowflake account name"
  type        = string
  default     = null
}

variable "snowflake_user" {
  description = "Snowflake username"
  type        = string
  default     = null
}

variable "snowflake_role" {
  description = "Snowflake role"
  type        = string
  default     = null
}

variable "snowflake_private_key" {
  description = "Snowflake private key for key-pair authentication"
  type        = string
  sensitive   = true
  default     = null
}
