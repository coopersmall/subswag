output "bound_members" {
  value = keys(google_service_account_iam_member.workload_identity_binding)
}
