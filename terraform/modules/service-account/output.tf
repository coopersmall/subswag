output "service_account_id" {
  value = google_service_account.github_actions.account_id
}

output "service_account_name" {
  value = google_service_account.github_actions.name
}

output "service_account_email" {
  value = google_service_account.github_actions.email
}

