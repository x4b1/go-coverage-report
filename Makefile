SHELL=/bin/bash -e -o pipefail
PWD = $(shell pwd)

BINARY_NAME=go-coverage-report

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build tools
set_opts: .PHONY

default: tidy install-tools

## generate:
generate: go/generate

go/generate:
	$(call print-target)
	@go generate ./...

## Deps:
download: ## Downloads the dependencies
	$(call print-target)
	@go mod download

update-deps:
	$(call print-target)
	@go get -u && go mod tidy
	@cd tools && go get -u && go mod tidy

tidy: ## Install development tools
	$(call print-target)
	@go mod tidy
	@cd tools && go mod tidy

install-tools: ## go install tools
	$(call print-target)
	@cd tools && go install $(shell cd tools && go list -f '{{ join .Imports " " }}' -tags=tools)

## Lint:
lint: download ## Lint go code
	$(call print-target)
	@golangci-lint run
lint-fix: download ## Lint go code and try to fix issues
	$(call print-target)
	@golangci-lint run --fix

## Test:
test: ## Runs all tests
	$(call print-target)
	@go test -count 1 ./...

coverage: out/report.json ## Displays coverage per func on cli
	$(call print-target)
	@go tool cover -func=out/cover.out

html-coverage: out/report.json ## Displays the coverage results in the browser
	$(call print-target)
	@go tool cover -html=out/cover.out

.PHONY: out/report.json
out/report.json:
	$(call print-target)
	@mkdir -p out
	@go test -count 1 ./... -covermode=atomic -coverpkg=./... -coverprofile=out/cover.out --json

test-build: ## Tests whether the code compiles
	$(call print-target)
	@go build -o /dev/null ./...

clean: ## Cleans up everything
	$(call print-target)
	@rm -rf bin out

## Migrations:
migration-create: ## Creates a new migration usage: `migration-create name=<migration name>`
	$(call print-target)
	@migrate create -dir ./deployment/migrations -ext sql $(name)

## Build:
build: ## Build project and put the output binary in bin/
	$(call print-target)
	mkdir -p bin
	@go build -o bin/$(BINARY_NAME) .

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef
