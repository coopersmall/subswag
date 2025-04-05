terraform {
  cloud {
    organization = "Microworlds" 
    workspaces { 
      name = "bootstrap" 
    }
  }
  required_version = "~> 1.5"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.10"
    }
  }
}

provider "google" {
  project      = var.project_id
  region       = var.region
}

module "bootstrap" {
  source = "../../modules/bootstrap"

  project_id  = var.project_id
  environment = var.environment
}
