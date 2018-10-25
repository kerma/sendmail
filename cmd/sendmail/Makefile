.DEFAULT_GOAL:=help
SHELL:=/bin/bash

.PHONY: help deps build install

# https://suva.sh/posts/well-documented-makefiles/
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

deps:  ## Check dependencies
	$(info Getting dependencies)
	go get -v .

build: deps ## Build the project
	$(info Building the project)
	go build

test:  ## Run tests 
	go test

install: deps  ## Install to $GOPATH/bin
	go install

