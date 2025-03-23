output "enabled_services" {
  description = "List of enabled GCP APIs"
  value       = google_project_service.services
}
