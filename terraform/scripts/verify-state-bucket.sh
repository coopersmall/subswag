#!/bin/bash
set -eo pipefail

if [ $# -ne 1 ]; then
  echo "Usage: $0 <bucket-name>"
  exit 1
fi

BUCKET_NAME="$1"

check_gsutil() {
  if ! command -v gsutil &> /dev/null; then
    echo "ERROR: gsutil not found. Install Google Cloud SDK first."
    exit 1
  fi
}

check_bucket_exists() {
  if ! gsutil ls -b "gs://${BUCKET_NAME}" &> /dev/null; then
    echo "ERROR: Bucket '${BUCKET_NAME}' does not exist."
    exit 1
  fi
  echo "Bucket '${BUCKET_NAME}' exists"
}

check_versioning() {
  echo "Checking versioning for gs://${BUCKET_NAME}"
  
  local status_output
  # Capture both stdout and stderr
  if ! status_output=$(gsutil versioning get "gs://${BUCKET_NAME}" 2>&1); then
    echo "ERROR: Failed to check versioning status:"
    echo "$status_output"
    exit 1
  fi

  # Check for Enabled status in the line about versioning
  if ! echo "$status_output" | rg -qw '^\s*gs://[^:]+:\s*Enabled\b'; then
    echo "ERROR: Versioning is NOT enabled for '${BUCKET_NAME}'"
    echo "Enable it with:"
    echo "  gsutil versioning set on gs://${BUCKET_NAME}"
    exit 1
  fi
  echo "Versioning is enabled"
}

check_state_files() {
  if ! gsutil ls "gs://${BUCKET_NAME}/*.tfstate*" &> /dev/null; then
    echo "WARNING: No existing state files found in '${BUCKET_NAME}/' (okay if first run)"
    return 0
  fi
  echo "State files found in '${BUCKET_NAME}/'"
}

main() {
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "🚀 Verifying Terraform state bucket"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  check_gsutil
  check_bucket_exists
  check_versioning
  check_state_files
  echo "All critical checks passed. Bucket is ready for Terraform state management."
}

main
exit 0

