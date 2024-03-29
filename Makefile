.DEFAULT_GOAL := help
PROJECT_NAME := $(shell basename "$(PWD)")
GCP_PROJECT := $(shell gcloud config get-value project)
PROJECT_NUMBER := $(shell gcloud projects describe  ${GCP_PROJECT} --format="value(projectNumber)")
##  Please change the region
REGION := asia-northeast1
JOB_NAME := $(PROJECT_NAME)-$(SBI_USERNAME)

.PHONY: help
help:
	@echo "\033[32m$(PROJECT_NAME):\033[0m"
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: go_module
go_module: ## Run go mod tidy
	go mod tidy

.PHONY: build
build: ## Build with buildpacks
	docker build --no-cache --platform amd64 -t $(PROJECT_NAME) -f Dockerfile  .
	docker tag $(PROJECT_NAME) asia.gcr.io/$(GCP_PROJECT)/$(PROJECT_NAME)

.PHONY: push
push: ## Push container image to GCR (Google Cloud Registry)
	docker push asia.gcr.io/$(GCP_PROJECT)/$(PROJECT_NAME)

.PHONY: deploy
deploy: ## Deploy to Google Cloud Run
	gcloud beta run jobs deploy $(JOB_NAME) \
	--image asia.gcr.io/$(GCP_PROJECT)/$(PROJECT_NAME):latest \
	--command '/main' \
	--region $(REGION) \
	--set-env-vars "SBI_USERNAME=$(SBI_USERNAME)" \
	--set-env-vars "SBI_PASSWORD=$(SBI_PASSWORD)" \
	--set-env-vars "SBI_TORIHIKI_PASSWORD=$(SBI_TORIHIKI_PASSWORD)"

.PHONY: exec
exec: ## Exec on Google Cloud Run
	gcloud beta run jobs execute $(JOB_NAME) \
	--region $(REGION)

.PHONY: test
test: ## Run go test
	go test ./...

.PHONY: errcheck
errcheck: ## Run errcheck
	errcheck -blank -asserts ./...

.PHONY: create_schedule
create_schedule: ## Run errcheck
	gcloud scheduler jobs create http $(JOB_NAME) \
	--location $(REGION) \
	--schedule "0 9 * * *" \
	--http-method=POST \
	--uri=https://$(REGION)-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/$(GCP_PROJECT)/jobs/$(JOB_NAME):run \
	--oauth-service-account-email=$(PROJECT_NUMBER)-compute@developer.gserviceaccount.com


