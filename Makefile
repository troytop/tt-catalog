.PHONY: help test

export CATALOG_LOCATION_DEVELOPMENT?=s3://helion-service-manager/release/catalog-templates/HCP-v1/stable-1/services
export CATALOG_LOCATION_STABLE?=s3://stackato-4/catalog-templates/stable-1/services
export IDL_LOCATION?=s3://helion-service-manager/release/instance-definition/stable-1

default: help

help:
	@echo "These 'make' targets are available."
	@echo
	@echo "  all                Run install tools and tests"
	@echo "  tools              Installs tools needed to run"
	@echo "  test               Run the unit tests"
	@echo "  test-format        Run the formatting tests"
	@echo "  publish-catalog    Publish the catalog to s3 location"
	@echo "  sign-development   Generate signatures for development bucket"
	@echo "  sign-stable        Generate signatures for stable bucket"

all: tools test

tools:
	@echo "$(OK_COLOR)==> Getting required tools$(NO_COLOR)"
	go get github.com/tools/godep

test: tools test-format
	@echo "$(OK_COLOR)==> Testing code$(NO_COLOR)"
	SKIP_SDL_VALIDATION="hce hcf hsc hcp" godep go test ./... -v

test-format:
	@echo "$(OK_COLOR)==> Checking code with gofmt$(NO_COLOR)"
	./scripts/testFmt.sh tests

publish-catalog:
	@echo "$(OK_COLOR)===> Publish catalog to s3 location @ ${}$(NO_COLOR)"
	./scripts/publish-catalog-bucket.sh

sign-development:
	@echo "$(OK_COLOR)===> Publish catalog to s3 location @ ${}$(NO_COLOR)"
	./scripts/generate-signature-development.sh '${PASSPHRASE}'

sign-stable:
	@echo "$(OK_COLOR)===> Publish catalog to s3 location @ ${}$(NO_COLOR)"
	./scripts/generate-signature-stable.sh '${PASSPHRASE}'