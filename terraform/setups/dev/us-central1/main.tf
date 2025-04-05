terraform {
  cloud {
    organization = "Microworlds"
    workspaces {
      name = "microworlds-dev-us-central1"
    }
  }
}

provider "google" {
  project     = "inbound-ranger-453921-k7"
  region      = "us-central1"
}
