terraform {
  backend "gcs" {
    bucket = "terraform_state_bucket-dev"
    prefix = "bootstrap"
  }
}
