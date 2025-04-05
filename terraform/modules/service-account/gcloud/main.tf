resource "google_service_account" "github_actions" {
  project      = var.project_id
  account_id   = var.service_account_id
  display_name = var.service_account_name
  description  = var.service_account_description
}
