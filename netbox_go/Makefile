# Project variables
BINARY_NAME=netbox_api
MAIN_PATH=cmd/api/main.go
BUILD_DIR=bin

# Go toolchain
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod

# Generation tools
SQLC=sqlc
MOCKGEN=mockgen

.PHONY: all build clean generate-sqlc generate-mocks generate test

# Default target
all: generate build

## Build the project
build:
	@echo "Building binary..."
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

## Generate SQLC code
generate-sqlc:
	@echo "Generating SQLC code..."
	cd netbox_go && $(SQLC) generate

## Generate Gomock mocks
generate-mocks:
	@echo "Generating mocks..."
	# This command finds all interfaces in internal/domain and generates mocks for them
	# Adjust the path if your interfaces are located elsewhere
	$(GOCMD) install github.com/golang/mock/mockgen
	find netbox_go/internal/domain -name "*.go" | xargs grep -l "type .* interface" | xargs -I {} $(MOCKGEN) -source={} -destination=netbox_go/internal/infrastructure/mocks/{} . {} = mocks_$(shell basename {})

## Run all generators
generate: generate-sqlc generate-mocks

## Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

## Run tests
test:
	$(GOCMD) test ./...
