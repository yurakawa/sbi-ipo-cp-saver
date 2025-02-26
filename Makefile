.DEFAULT_GOAL := help
PROJECT_NAME := $(shell basename "$(PWD)")
# GCP_PROJECT := $(GCP_PROJECT_ID)
PROJECT_NUMBER := $(shell gcloud projects describe  ${GCP_PROJECT} --format="value(projectNumber)")
# REGION := $(REGINON)
JOB_NAME := $(PROJECT_NAME)-$(SBI_USERNAME)
PASSWORD_SECRET_NAME := projects/$(GCP_PROJECT_ID)/secrets/$(SBI_USERNAME)-password
TORIHIKI_PASSWORD_SECRET_NAME := projects/$(GCP_PROJECT_ID)/secrets/$(SBI_USERNAME)-torihiki-password
SERVICE_ACCOUNT_NAME := sbiipocpsaver-sa-runoncloudrun

.PHONY: help
help:
	@echo "\033[32m$(PROJECT_NAME):\033[0m"
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: go_module
go_module: ## Run go mod tidy
	go mod tidy

.PHONY: test
test: ## Run go test
	go test ./...

.PHONY: errcheck
errcheck: ## Run errcheck
	errcheck -blank -asserts ./...

.PHONY: run
run: ## Run the application on local
	go run main.go

### The following commands are for Google Cloud Run

.PHONY: initializing
initializing:
	create_registry \
	create_service_account \ 
	create_secrets \
	configure \
	update \
	create_schedule

.PHONY: update
update: build_image \
	push_image \
	deploy

.PHONY: create_registry
create_registry: ## Create Artifact Registry
	gcloud auth configure-docker $(REGION)-docker.pkg.dev
	-@gcloud artifacts repositories create $(PROJECT_NAME) \
	      --repository-format=docker \
	      --location=$(REGION) \
	      --description="$(PROJECT_NAME)"  \
	      --async

.PHONY: build_image
build_image: ## Build with buildpacks
	docker build --platform linux/amd64 --no-cache -t $(PROJECT_NAME) -f Dockerfile  .
	docker tag $(PROJECT_NAME) $(REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/$(PROJECT_NAME)/$(PROJECT_NAME)

.PHONY: push_image 
push_image: ## Push container image to Artifact Registry
	docker push $(REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/$(PROJECT_NAME)/$(PROJECT_NAME):latest

.PHONY: create_secrets
create_secrets:
	gcloud secrets delete $(PASSWORD_SECRET_NAME) --project=$(GCP_PROJECT_ID) || true
	gcloud secrets delete $(TORIHIKI_PASSWORD_SECRET_NAME) --project=$(GCP_PROJECT_ID) || true
  
	# Create secrets in GCP Secret Manager for password only
	echo "Creating secrets in GCP Secret Manager..." \
	printf $(SBI_PASSWORD) | gcloud secrets create $(PASSWORD_SECRET_NAME) \
		--data-file=- \
		--project=$(GCP_PROJECT_ID) \
		--replication-policy=automatic || true
	
	printf $(SBI_TORIHIKI_PASSWORD) | gcloud secrets create $(TORIHIKI_PASSWORD_SECRET_NAME) \
		--data-file=- \
		--project=$(GCP_PROJECT_ID) \
		--replication-policy=automatic || true

.PHONY: create_secrets
update_secrets:
	printf $(SBI_PASSWORD) | gcloud secrets versions add $(PASSWORD_SECRET_NAME) \
		--data-file=- \
		--project=$(GCP_PROJECT_ID)

	printf $(SBI_TORIHIKI_PASSWORD) | gcloud secrets versions add $(TORIHIKI_PASSWORD_SECRET_NAME) \
		--data-file=- \
		--project=$(GCP_PROJECT_ID)



.PHONY: deploy
deploy: ## Deploy to Google Cloud Run
	gcloud run jobs deploy $(JOB_NAME) \
	--image $(REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/$(PROJECT_NAME)/$(PROJECT_NAME):latest \
	--command '/main' \
	--region $(REGION) \
	--memory 1Gi \
	--cpu 2 \
	--service-account=${SERVICE_ACCOUNT_NAME}@${GCP_PROJECT_ID}.iam.gserviceaccount.com \
	--set-env-vars "SBI_USERNAME=$(SBI_USERNAME)" \
	--set-env-vars "ENV=gcp" \
	--set-env-vars "GCP_PROJECT_ID=$(GCP_PROJECT_ID)" 

.PHONY: create_service_account
create_service_account: ## Create a service account and grant access to secrets
	-@gcloud iam service-accounts create ${SERVICE_ACCOUNT_NAME} \
	    --project=${GCP_PROJECT_ID}

	-@gcloud secrets add-iam-policy-binding ${PASSWORD_SECRET_NAME} \
	    --member="serviceAccount:${SERVICE_ACCOUNT_NAME}@${GCP_PROJECT_ID}.iam.gserviceaccount.com" \
	    --role="roles/secretmanager.secretAccessor"

	-@gcloud secrets add-iam-policy-binding ${TORIHIKI_PASSWORD_SECRET_NAME} \
	    --member="serviceAccount:${SERVICE_ACCOUNT_NAME}@${GCP_PROJECT_ID}.iam.gserviceaccount.com" \
	    --role="roles/secretmanager.secretAccessor"

.PHONY: create_schedule
create_schedule: ## Creating a schedule after deleting the existing schedule
	gcloud scheduler jobs create http $(JOB_NAME) \
	--location $(REGION) \
	--schedule "0 9 * * *" \
	--http-method=POST \
	--uri=https://$(REGION)-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/$(GCP_PROJECT)/jobs/$(JOB_NAME):run \
	--oauth-service-account-email=$(PROJECT_NUMBER)-compute@developer.gserviceaccount.com

.PHONY: exec
exec: ## Exec on Google Cloud Run
	gcloud beta run jobs execute $(JOB_NAME) \
	--region $(REGION)

