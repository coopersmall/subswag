#!/bin/bash
set -eo pipefail

# Validate parameters
if [ $# -ne 3 ]; then
  echo "Usage: $0 <project-id> <environment> <region>"
  exit 1
fi

PROJECT_ID="$1"
ENVIRONMENT="$2"
REGION="$3"

case $ENVIRONMENT in
  "dev")
    ROOT="$(cd "$(dirname "$0")/../setups/dev" && pwd)"
    ;;
  "prod")
    ROOT="$(cd "$(dirname "$0")/../setups/prod" && pwd)"
    ;;
  *)
    echo "Invalid environment: $ENVIRONMENT"
    exit 1
    ;;
esac

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
TF_DIR="${ROOT}/bootstrap"

# Authentication flow
authenticate_gcloud() {
  echo "Verifying gcloud login..."
  if ! gcloud auth list > /dev/null 2>&1; then
    echo "Please log in to gcloud first using 'gcloud auth login'"
    exit 1
  fi
  echo "gcloud login verified!"
  
  # Export credentials for Terraform
  export GOOGLE_OAUTH_ACCESS_TOKEN=$(gcloud auth print-access-token)
}

setup_backend() {
  local dir="$1"
  local bucket_name="$2"
  
  echo "Processing directory: $dir"
  cd "$dir"
  
  echo "Generating backend configuration..."
  cat > backend.tf <<EOF
terraform {
  backend "gcs" {
    bucket = "${bucket_name}"
    prefix = "$(basename "$dir")"
  }
}
EOF

  echo "Initializing and migrating state..."
  terraform init -migrate-state -force-copy -input=false

  echo "Cleaning up local state files..."
  rm -f terraform.tfstate* .terraform.lock.hcl
  rm -rf .terraform
  
  cd ..
}

verify_bucket() {
  echo "Verifying bucket..."
  "${SCRIPT_DIR}/verify-state-bucket.sh" "${BUCKET_NAME}"
}

main() {
  authenticate_gcloud
  
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "🚀 Bootstrapping Terraform state for ${ENVIRONMENT}"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  
  cd "${TF_DIR}"
  
  echo "Initializing Terraform..."
  terraform init -input=false

  echo "Validating configuration..."
  terraform validate

  echo "Applying bootstrap configuration..."
  terraform apply \
    -var "project_id=${PROJECT_ID}" \
    -var "environment=${ENVIRONMENT}" \
    -var "region=${REGION}" \
    -var "access_token=${GOOGLE_OAUTH_ACCESS_TOKEN}"
    
  BUCKET_NAME=$(terraform output -raw bucket_name)
  
  # Setup backend for all directories under ROOT
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "🚀 Setting up backends for all components"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

  cd "${ROOT}"
  echo "Processing components in: $(pwd)"
  for dir in *; do
    if [ -d "${dir}" ] && [ -f "${dir}/main.tf" ]; then
      setup_backend "${dir}" "${BUCKET_NAME}"
    fi
  done

  verify_bucket
}

main

