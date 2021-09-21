include .env

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

PLATFORM?=linux/arm64

# Redirect error output to a file, so we can show it in development mode.
STDERR=/tmp/weather-sensor-bridge-stderr.txt

## clean: Clean build files.
clean:
	@echo "  >  Cleaning build cache"
	@go clean --mod=mod
	@rm -rf out
	@rm -rf internal/mocks

deps:
	@echo "  >  Getting dependencies..."
	@go install github.com/golang/mock/mockgen@latest
	@go mod download

generate:
	@echo "  >  Generate code..."
	@go generate ./...

## test: Clean and run all unit tests
test: clean deps generate
	@echo "  >  Running tests..."
	@mkdir -p out
	@go test -v -coverprofile=./out/coverage.out --mod=mod ./...

## coverage: Show unit test coverage report
coverage: test
	@echo "  >  Parsing coverage..."
	@go tool cover -html=./out/coverage.out

compile: test
	@echo "  >  Building binary..."
	@go build --mod=mod -o ./out/weather-sensor-bridge ./cmd/weather-sensor-bridge/main.go

## build: Compile the binary.
build:
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s compile 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2

## docker-build: Builds the docker image, defaults to linux/amd64 platform can be specified by platform=<platform>.
docker-build: clean
	@echo "  >  Building docker image..."
	@docker buildx build \
		--platform $(PLATFORM) \
		-t ghcr.io/geoff-coppertop/weather-sensor-bridge:latest \
		--load .

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in weather-sensor-bridge:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
