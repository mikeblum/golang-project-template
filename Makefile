## golang-project-template build tooling

MAKEFLAGS += --silent

GOLANGCI_LINT_VERSION = v1.64.6

all: help

## help: Prints a list of available build targets.
help:
	@echo "\033[1;34mUsage:\033[0m make \033[1;36m<OPTIONS>\033[0m ... \033[1;36m<TARGETS>\033[0m"
	@echo ""
	@echo "\033[1;35mAvailable targets are:\033[0m"
	@echo ''
	@sed -n 's/^##//p' ${PWD}/Makefile | column -t -s ':' | sed -e 's/^/ /' | sed -e 's/\(.*\)/\x1b[32m\1\x1b[0m/'
	@echo
	@echo "\033[1;35mTargets run by default are:\033[0m \033[1;32m`sed -n 's/^all: //p' ./Makefile | sed -e 's/ /, /g' | sed -e 's/\(.*\), /\1, and /'`\033[0m"

## lint: Lint with golangci-lint
lint:
	docker run --rm -v $$(pwd):/repo -w /repo golangci/golangci-lint:${GOLANGCI_LINT_VERSION} golangci-lint run --verbose --color always ./...

## fmt: Format with gofmt
fmt:
	go fmt ./...

# tidy: Tidy with go mod tidy
tidy:
	go mod tidy

## pre-commit: Chain lint + test + scan
pre-commit: test lint vuln

## run: go run main.go
run:
	go run main.go

## test: Test with go test
test:
	go test -test.v -race -covermode=atomic -coverprofile=coverage.out ./... && go tool cover -html=coverage.out && rm coverage.out

## test-perf: Benchmark tests with go test -bench
test-perf:
	go test -test.v -benchmem -bench=. -coverprofile=coverage-bench.out ./... && go tool cover -html=coverage-bench.out && rm coverage-bench.out

## vuln: Scan against the Go vulnerability database
vuln:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

## badges: Generate README badges
badge badges:
	go test ./... -coverprofile coverage.out \
	COVERAGE=$(go tool cover -func=coverage.out | grep total: | grep -Eo '[0-9]+\.[0-9]+') \
	curl -sL "https://img.shields.io/static/v1?label=coverage&message=$$COVERAGE%&color=$$COLOR&logo=go" > images/badges/coverage.svg

.PHONY: lint fmt tidy pre-commit test test-perf vuln badges
