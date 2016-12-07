export S3_CATALOG_LOCATION?=s3://helion-service-manager/build/dev-catalog
export SKIP_SDL_VALIDATION?="hce hcf hsc"

.PHONY: help test

default: help

help:
	@echo "These 'make' targets are available."
	@echo
	@echo "  tools              Installs tools needed to run"
	@echo "  test               Run the unit tests"
	@echo "  test-format        Run the formatting tests"
	@echo "  publish-services   Publish service sdls"

all: tools test

tools:
	@echo "$(OK_COLOR)==> Getting required tools$(NO_COLOR)"
	go get github.com/tools/godep

test: tools test-format
	@echo "$(OK_COLOR)==> Testing code$(NO_COLOR)"
	go test $$(go list ./... | grep -v vendor) -v

test-format:
	@echo "$(OK_COLOR)==> Checking code with gofmt$(NO_COLOR)"
	./scripts/testFmt.sh tests

publish-services:
	@echo "$(OK_COLOR)==> Publish service sdls$(NO_COLOR)"
	./scripts/publish-services.bash
