module "api_services" {
  source     = "../api-services"
  project_id = var.project_id
  services = [
    "cloudresourcemanager.googleapis.com",
    "iam.googleapis.com",
    "sts.googleapis.com",
    "iamcredentials.googleapis.com",
  ]
}

module "github_service_account" {
  source                      = "../service-account"
  project_id                  = var.project_id
  service_account_id          = "github-sa-${var.environment}"
  service_account_name        = "GitHub SA ${var.environment}"
  service_account_description = "Service account for GitHub Actions"
  depends_on                  = [module.api_services]
}

module "github_workload_identity" {
  source                = "../workload-identity/create"
  project_id            = var.project_id
  pool_id               = "github-wif-pool-${var.environment}"
  pool_display_name     = "GitHub Pool ${var.environment}"
  pool_description      = "Workload identity pool for GitHub Actions (${var.environment})"
  provider_id           = "github-com-oidc"
  provider_display_name = "GitHub WIF Provider ${var.environment}"
  oidc_issuer_uri       = "https://token.actions.githubusercontent.com"
  attribute_condition   = "attribute.repository == '${var.repo}'"
  attribute_mappings = {
    # Default attributes used in:
    # https://registry.terraform.io/modules/Cyclenerd/wif-service-account/google/latest
    "google.subject"       = "assertion.sub"        # Subject
    "attribute.sub"        = "attribute.sub"        # Subject
    "attribute.repository" = "assertion.repository" # The repository from where the workflow is running
    # More
    # https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect#understanding-the-oidc-token
    "attribute.aud"                   = "attribute.aud"                   # Audience
    "attribute.iss"                   = "attribute.iss"                   # Issuer
    "attribute.actor"                 = "assertion.actor"                 # The personal account that initiated the workflow run.
    "attribute.actor_id"              = "assertion.actor_id"              # The ID of personal account that initiated the workflow run.
    "attribute.base_ref"              = "assertion.base_ref"              # The target branch of the pull request in a workflow run.
    "attribute.environment"           = "assertion.environment"           # The name of the environment used by the job.
    "attribute.event_name"            = "assertion.event_name"            # The name of the event that triggered the workflow run.
    "attribute.head_ref"              = "assertion.head_ref"              # The source branch of the pull request in a workflow run.
    "attribute.job_workflow_ref"      = "assertion.job_workflow_ref"      # For jobs using a reusable workflow, the ref path to the reusable workflow. For more information, see "Using OpenID Connect with reusable workflows."
    "attribute.job_workflow_sha"      = "assertion.job_workflow_sha"      # For jobs using a reusable workflow, the commit SHA for the reusable workflow file.
    "attribute.ref"                   = "assertion.ref"                   # (Reference) The git ref that triggered the workflow run.
    "attribute.ref_type"              = "assertion.ref_type"              # The type of ref, for example: "branch".
    "attribute.repository_visibility" = "assertion.repository_visibility" # The visibility of the repository where the workflow is running. Accepts the following values: internal, private, or public.
    "attribute.repository_id"         = "assertion.repository_id"         # The ID of the repository from where the workflow is running.
    "attribute.repository_owner"      = "assertion.repository_owner"      # The name of the organization in which the repository is stored.
    "attribute.repository_owner_id"   = "assertion.repository_owner_id"   # The ID of the organization in which the repository is stored.
    "attribute.run_id"                = "assertion.run_id"                # The ID of the workflow run that triggered the workflow.
    "attribute.run_number"            = "assertion.run_number"            # The number of times this workflow has been run.
    "attribute.run_attempt"           = "assertion.run_attempt"           # The number of times this workflow run has been retried.
    "attribute.runner_environment"    = "assertion.runner_environment"    # The type of runner used by the job. Accepts the following values: github-hosted or self-hosted.
    "attribute.workflow"              = "assertion.workflow"              # The name of the workflow.
    "attribute.workflow_ref"          = "assertion.workflow_ref"          # The ref path to the workflow. For example, octocat/hello-world/.github/workflows/my-workflow.yml@refs/heads/my_branch.
    "attribute.workflow_sha"          = "assertion.workflow_sha"          # The commit SHA for the workflow file.
  }
  depends_on = [module.api_services]
}

module "github_workload_binding" {
  source             = "../workload-identity/bind"
  service_account_id = module.github_service_account.service_account_name
  member             = "principalSet://iam.googleapis.com/${module.github_workload_identity.pool_name}/attribute.repository/${var.repo}"
  depends_on         = [module.github_workload_identity, module.github_service_account, module.api_services]
}

module "github_iam_roles" {
  source                = "../iam-roles"
  project_id            = var.project_id
  service_account_email = module.github_service_account.service_account_email
  roles                 = var.github_service_account_roles
  depends_on            = [module.github_workload_binding, module.api_services]
}
