module "github-wif" {
  source     = "../github-wif"
  project_id = var.project_id
  region     = var.region
  environment = var.environment
}
