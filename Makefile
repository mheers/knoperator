## This is a self-documented Makefile. For usage information, run `make help`:
##
## For more information, refer to https://suva.sh/posts/well-documented-makefiles/

SHELL := /bin/bash

all: help

##@ Building
build: docker ##  Builds the application (same as 'docker')

install-dependencies: ## installs dependencies
	go mod download

get-binary: docker ## Builds the knoperator application and outputs a single binary to 'knoperator'
	docker run --rm --entrypoint cat knoperator /bin/knoperator > knoperator
	chmod +x knoperator

set-version: ## Sets the version
	./ci/set-version.sh

docker: set-version ##  Builds the knoperator application
	docker buildx build --platform linux/amd64 -t knoperator --output type=docker .

docker-arm64: set-version ##  Builds the knoperator application for arm64
	docker buildx build --platform linux/arm64 -t knoperator --output type=docker .

docker-multi: set-version ##  Builds the knoperator application for amd64 and arm64
	docker buildx build --platform linux/amd64,linux/arm64 -t knoperator --push .

##@ Starting
start-knoperator: ##  Starts the application as docker container
	docker-compose up knoperator

start-mq: ##  Starts the mq
	docker-compose up nats

##@ Testing
test-unit: ## Starts unit tests
	go test ./... -coverprofile cover.out
	go tool cover -func cover.out
	rm cover.out

test-unit-race: ## Starts unit tests in race detection mode
	go test ./... -race -coverprofile cover.out
	go tool cover -func cover.out
	rm cover.out

##@ Scanning
scan: docker ## Runs a security scan
	trivy knoperator

##@ Helpers

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
