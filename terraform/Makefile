BUCKET_NAME ?= terraform-state-bucket
ENV_PREFIX ?= environments

.PHONY: bootstrap
bootstrap: bootstrap-dev

.PHONY: bootstrap-dev
bootstrap-dev:
	@chmod +x ./scripts/bootstrap.sh
	@./scripts/bootstrap.sh $(shell gcloud config get-value project) "dev" "us-central1"

.PHONY: fmt
fmt:
	@terraform fmt -recursive

.PHONY: verify-state
verify-state: verify-bucket-state

.PHONY: verify-bucket-state
verify-bucket-state:
	@chmod +x ./scripts/verify-state-bucket.sh
	@echo "Verifying GCS Dev bucket configuration..."
	@./scripts/verify-state-bucket.sh "terraform_state_bucket-dev"
