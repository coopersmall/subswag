variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "services" {
  description = "List of GCP APIs to enable"
  type        = list(string)
  default = [
    "iam.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "sts.googleapis.com",
    "iamcredentials.googleapis.com"
  ]
}
