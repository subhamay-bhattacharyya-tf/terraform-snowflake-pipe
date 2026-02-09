output "pipe_names" {
  description = "The names of the created pipes"
  value       = module.pipes.pipe_names
}

output "pipe_fully_qualified_names" {
  description = "The fully qualified names of the pipes"
  value       = module.pipes.pipe_fully_qualified_names
}

output "pipe_notification_channels" {
  description = "The notification channels for the pipes"
  value       = module.pipes.pipe_notification_channels
}

output "pipe_owners" {
  description = "The owners of the pipes"
  value       = module.pipes.pipe_owners
}

output "pipes" {
  description = "All pipe resources"
  value       = module.pipes.pipes
}
