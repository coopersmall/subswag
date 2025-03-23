resource "google_storage_bucket" "state" {
  name          = var.name
  location      = var.region
  force_destroy = false

  versioning {
    enabled = var.versioning
  }
}
