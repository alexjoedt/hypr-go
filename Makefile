## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
	
## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v
	
## test: run the test suite with the race detector and report coverage
.PHONY: test
test:
	go test -race -vet=off -coverprofile=cover.out ./...
	go tool cover -func=cover.out | tail -n 1

## lint: run golangci-lint
.PHONY: lint
lint:
	golangci-lint run ./...

## audit: run quality control checks
.PHONY: audit
audit:
	go vet ./...
	golangci-lint run ./...
	go test -race -vet=off ./...
	go mod verify
	
## build: compile the library and the example binary
.PHONY: build
build:
	go mod verify
	go build ./...
	go build -ldflags='-s' -o=./bin/example ./examples

## run: run the example against the running Hyprland instance
.PHONY: run
run: build
	./bin/example
