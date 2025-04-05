variable "project_id" {
  type = string
}

variable "service_account_email" {
  description = "The email address of the service account."
  type        = string
}

variable "roles" {
  description = "The roles to grant to the service account."
  type        = list(string)
  default = [
    "roles/run.developer",
    "roles/artifactregistry.writer"
  ]
}
