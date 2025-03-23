variable "service_account_id" {
  type        = string
  description = "The service account ID"
}

variable "member" {
  type        = string
  description = "Principal in format principalSet://iam.googleapis.com/<pool-name>/<attribute>"
}
