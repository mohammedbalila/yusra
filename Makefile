# Makefile for yusra Go project
PACKAGE_NAME := yusra
VERSION := 0.0.1

# Build for all platforms
build-all: build-linux build-windows build-macos

build:
	@echo "Building the project"
	@CGO_ENABLED=1 go build -o "$(PACKAGE_NAME)" -ldflags "-X main.Version=$(VERSION)" ./

# Build for Linux optimized
build-linux:
	@echo "Building the project for Linux"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o "$(PACKAGE_NAME)-linux" -ldflags "-X main.Version=$(VERSION)" ./

build-macos:
	@echo "Building the project for MacOS"
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o "$(PACKAGE_NAME)-macos" -ldflags "-X main.Version=$(VERSION)" ./

# Build for Windows optimized
build-windows:
	@echo "Building the project for Windows"
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o "$(PACKAGE_NAME)-windows".exe -ldflags "-X main.Version=$(VERSION)" ./

# Clean the build artifacts
clean:
	@echo "Cleaning up..."
	@rm $(PACKAGE_NAME)
	@go clean
	@echo "Cleanup complete."
