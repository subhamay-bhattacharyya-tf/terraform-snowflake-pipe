output "pipe_names" {
  description = "The names of the created pipes."
  value       = { for k, v in snowflake_pipe.this : k => v.name }
}

output "pipe_fully_qualified_names" {
  description = "The fully qualified names of the pipes."
  value       = { for k, v in snowflake_pipe.this : k => v.fully_qualified_name }
}

output "pipe_notification_channels" {
  description = "The notification channels for the pipes (for SNS integration)."
  value       = { for k, v in snowflake_pipe.this : k => v.notification_channel }
}

output "pipe_owners" {
  description = "The owners of the pipes."
  value       = { for k, v in snowflake_pipe.this : k => v.owner }
}

output "pipes" {
  description = "All pipe resources."
  value       = snowflake_pipe.this
}
