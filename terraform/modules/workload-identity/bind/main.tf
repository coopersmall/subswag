resource "google_service_account_iam_member" "workload_identity_binding" {
  service_account_id = var.service_account_id
  role               = "roles/iam.workloadIdentityUser"
  member             = var.member
}
