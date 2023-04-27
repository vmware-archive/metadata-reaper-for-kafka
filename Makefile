# Image URL to use all building/pushing image targets
TAG ?= latest
IMG ?= vmwaresaas.jfrog.io/vdp/source/kafka-metadata-reaper

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: test
test: ## Run go test against code.
	go test ./... -count=1 -v

.PHONY: test-ci
test-ci: ## Run go test against code.
	go test --tags=ci ./... -count=1 -v

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: tidy
tidy: ## Run go mod tidy
	go mod tidy

.PHONY: vendor
vendor: go.mod go.sum
	go mod vendor

##@ Build

.PHONY: build
build: tidy fmt vet vendor  ## Build manager binary.
	go build -o bin/kafka-metadata-reaper

.PHONY: run-topics-reaper
run-topics-reaper: fmt vet
	go run main.go topics

run-topics-reaper-no-dryrun: fmt vet
	go run main.go topics --dry-run=false

.PHONY: build-image
build-image: vendor ## Build docker image with the manager.
	docker build -t ${IMG}:${TAG} $(DOCKER_BUILD_ARGS) .

.PHONY: push-image
push-image: build-image ## Push docker image with the manager.
	docker push ${IMG}:${TAG}

.PHONY: start-environ
start-environ:
	 docker-compose -f compose.yaml up
