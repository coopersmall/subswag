output "bucket_id" {
  value = module.state.bucket_id
}

output "bucket_name" {
  value = module.state.bucket_name
}

output "github_service_account_email" {
  value = module.github_service_account.service_account_email
}
