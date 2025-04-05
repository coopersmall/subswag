terraform {
  backend "gcs" {
    bucket = "tf_state-dev-us-central1"
    prefix = "us-central1"
  }
}
