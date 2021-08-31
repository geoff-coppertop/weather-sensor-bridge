include .env

# Go related variables.
GOBASE=$(shell pwd | sed 's/ /\\ /g')
GOPATH="$(GOBASE)/vendor:$(GOBASE)"
GOBIN=$(GOBASE)/bin

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

# Redirect error output to a file, so we can show it in development mode.
STDERR=/tmp/weather-sensor-bridge-stderr.txt

## clean: Clean build files.
clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean -mod=mod
	@rm -rf bin
	@rm -rf out

deps:
	@echo "  >  Getting binary dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod download

## test: Generate and run all unit tests
test: clean deps
	@echo "  >  Running tests..."
	@mkdir -p out
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go test -v -coverprofile=./out/coverage.out -mod=mod ./...

## coverage: Show unit test coverage report
coverage: test
	@echo "  >  Parsing coverage..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go tool cover -html=./out/coverage.out

compile: clean deps test
	@echo "  >  Building binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -mod=mod -o $(GOBIN)/weather-sensor-bridge $(GOBASE)/cmd/weather-sensor-bridge/main.go

## build: Compile the binary.
build:
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s compile 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in weather-sensor-bridge:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo