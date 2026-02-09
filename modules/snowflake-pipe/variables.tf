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

  validation {
    condition     = alltrue([for k, pipe in var.pipe_configs : length(pipe.database) > 0])
    error_message = "Database name must not be empty."
  }

  validation {
    condition     = alltrue([for k, pipe in var.pipe_configs : length(pipe.schema) > 0])
    error_message = "Schema name must not be empty."
  }

  validation {
    condition     = alltrue([for k, pipe in var.pipe_configs : length(pipe.name) > 0])
    error_message = "Pipe name must not be empty."
  }

  validation {
    condition     = alltrue([for k, pipe in var.pipe_configs : length(pipe.copy_statement) > 0])
    error_message = "Copy statement must not be empty."
  }

  validation {
    condition = alltrue([
      for k, pipe in var.pipe_configs :
      !pipe.auto_ingest || (pipe.auto_ingest && (pipe.aws_sns_topic_arn != null || pipe.integration != null))
    ])
    error_message = "When auto_ingest is true, either aws_sns_topic_arn or integration must be provided."
  }
}
