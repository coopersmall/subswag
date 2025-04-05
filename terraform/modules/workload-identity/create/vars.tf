variable "project_id" {
  description = "The project ID"
}

variable "pool_id" {
  description = "The workload identity pool ID"
}

variable "pool_display_name" {
  description = "The display name of the workload identity pool"
}

variable "pool_description" {
  description = "The description of the workload identity pool"
}

variable "provider_id" {
  description = "The workload identity pool provider ID"
}

variable "provider_display_name" {
  description = "The display name of the workload identity pool provider"
}

variable "oidc_issuer_uri" {
  description = "The OIDC issuer URI"
}

variable "attribute_mappings" {
  type        = map(string)
  description = "The attribute mappings for the workload identity pool provider"
}

variable "attribute_condition" {
  type        = string
  description = "(Optional) Workload Identity Pool Provider attribute condition expression"
  default     = null
}
