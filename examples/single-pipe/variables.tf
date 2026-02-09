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
  default = {}
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
