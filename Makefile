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


.PHONY: configure_and_schedule
configure_and_schedule: configure \
	build \
	push \
	deploy \
	create_schedule

.PHONY: configure
configure:
	gcloud auth configure-docker $(REGION)-docker.pkg.dev
	-@gcloud artifacts repositories create $(PROJECT_NAME) \
	      --repository-format=docker \
	      --location=$(REGION) \
	      --description="$(PROJECT_NAME)"  \
	      --async


.PHONY: build
build: ## Build with buildpacks
	docker build --platform linux/amd64 --no-cache -t $(PROJECT_NAME) -f Dockerfile  .
	docker tag $(PROJECT_NAME) $(REGION)-docker.pkg.dev/$(GCP_PROJECT)/$(PROJECT_NAME)/$(PROJECT_NAME)

.PHONY: push
push: ## Push container image to Artifact Registry
	docker push $(REGION)-docker.pkg.dev/$(GCP_PROJECT)/$(PROJECT_NAME)/$(PROJECT_NAME):latest

.PHONY: deploy
deploy: ## Deploy to Google Cloud Run
	gcloud beta run jobs deploy $(JOB_NAME) \
	--image $(REGION)-docker.pkg.dev/$(GCP_PROJECT)/$(PROJECT_NAME)/$(PROJECT_NAME):latest \
	--command '/main' \
	--region $(REGION) \
	--memory 1Gi \
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


