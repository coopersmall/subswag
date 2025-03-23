variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "environment" {
  description = "The environment"
  type        = string
}

variable "region" {
  description = "The region"
  type        = string
}

variable "repo" {
  description = "GitHub repository in ORG/REPO format"
  type        = string
  default     = "coopersmall/subswag"
}

variable "github_service_account_roles" {
  description = "List of IAM roles to grant"
  type        = list(string)
  default = [
    "roles/container.clusterAdmin", # Create/manage clusters
    "roles/container.developer",    # Deploy workloads
    "roles/compute.networkUser",    # Manage networking components
    "roles/compute.securityAdmin",  # Firewall rules
    "roles/iam.serviceAccountUser", # Required for Workload Identity
    "roles/storage.objectAdmin"     # Only if using GCS for Terraform state
  ]
}
