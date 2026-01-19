# SNIP - Code Snippet Manager Makefile

# Variables
APP_NAME := snip
VERSION := $(shell git describe --tags --always --dirty 2>nul || echo dev)
BUILD_TIME := $(shell powershell -Command "Get-Date -Format 'yyyy-MM-dd_HH:mm:ss'")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>nul || echo unknown)

# Go variables
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt

# Build output directory
BUILD_DIR := build
DIST_DIR := dist

# Linker flags to inject version info
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -s -w"

# Detect OS
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    RM := cmd /C del /Q /F
    RMDIR := cmd /C rmdir /S /Q
    MKDIR := cmd /C mkdir
    SHELL := cmd
else
    DETECTED_OS := $(shell uname -s)
    RM := rm -f
    RMDIR := rm -rf
    MKDIR := mkdir -p
endif

.PHONY: all build build-all clean test coverage fmt vet lint install uninstall run help deps tidy

# Default target
all: clean fmt vet test build

# Build for current platform
build:
	@echo Building $(APP_NAME) for current platform...
	@$(MKDIR) $(BUILD_DIR) 2>nul || echo.
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME).exe .
	@echo Build complete: $(BUILD_DIR)/$(APP_NAME).exe

# Build for all platforms (Windows PowerShell)
build-all-windows:
	@echo Building for all platforms...
	@if not exist $(DIST_DIR) mkdir $(DIST_DIR)
	@echo Building for linux/amd64...
	@set GOOS=linux&& set GOARCH=amd64&& $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 .
	@echo Building for linux/arm64...
	@set GOOS=linux&& set GOARCH=arm64&& $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-linux-arm64 .
	@echo Building for darwin/amd64...
	@set GOOS=darwin&& set GOARCH=amd64&& $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-darwin-amd64 .
	@echo Building for darwin/arm64...
	@set GOOS=darwin&& set GOARCH=arm64&& $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-darwin-arm64 .
	@echo Building for windows/amd64...
	@set GOOS=windows&& set GOARCH=amd64&& $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-windows-amd64.exe .
	@echo All builds complete in $(DIST_DIR)/

# Build for all platforms (Unix/Linux/Mac)
build-all-unix:
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	@echo "Building for linux/amd64..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 .
	@echo "Building for linux/arm64..."
	@GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-linux-arm64 .
	@echo "Building for darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-darwin-amd64 .
	@echo "Building for darwin/arm64..."
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-darwin-arm64 .
	@echo "Building for windows/amd64..."
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-windows-amd64.exe .
	@echo "All builds complete in $(DIST_DIR)/"

# Auto-detect and run appropriate build-all
ifeq ($(OS),Windows_NT)
build-all: build-all-windows
else
build-all: build-all-unix
endif

# Clean build artifacts (Windows)
clean-windows:
	@echo Cleaning build artifacts...
	@if exist $(BUILD_DIR) rmdir /S /Q $(BUILD_DIR)
	@if exist $(DIST_DIR) rmdir /S /Q $(DIST_DIR)
	@$(GOCLEAN)
	@echo Clean complete

# Clean build artifacts (Unix)
clean-unix:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@$(GOCLEAN)
	@echo "Clean complete"

# Auto-detect and run appropriate clean
ifeq ($(OS),Windows_NT)
clean: clean-windows
else
clean: clean-unix
endif

# Run tests
test:
	@echo Running tests...
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	@echo Running tests with coverage...
	@$(MKDIR) $(BUILD_DIR) 2>nul || echo.
	$(GOTEST) -v -coverprofile=$(BUILD_DIR)/coverage.out ./...
	$(GOCMD) tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo Coverage report: $(BUILD_DIR)/coverage.html

# Format code
fmt:
	@echo Formatting code...
	$(GOFMT) ./...

# Run go vet
vet:
	@echo Running go vet...
	$(GOCMD) vet ./...

# Run linter (requires golangci-lint)
lint:
	@echo Running linter...
	@where golangci-lint >nul 2>nul && golangci-lint run ./... || echo golangci-lint not installed. Install from https://golangci-lint.run/usage/install/

# Install binary to system (Windows - requires admin)
install-windows: build
	@echo Installing $(APP_NAME) to C:\Program Files\$(APP_NAME)\...
	@if not exist "C:\Program Files\$(APP_NAME)" mkdir "C:\Program Files\$(APP_NAME)"
	@copy $(BUILD_DIR)\$(APP_NAME).exe "C:\Program Files\$(APP_NAME)\"
	@echo Installation complete. Add C:\Program Files\$(APP_NAME) to your PATH

# Install binary to system (Unix)
install-unix: build
	@echo "Installing $(APP_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/
	@echo "Installation complete"

# Auto-detect and run appropriate install
ifeq ($(OS),Windows_NT)
install: install-windows
else
install: install-unix
endif

# Uninstall binary from system
uninstall:
	@echo Uninstalling $(APP_NAME)...
ifeq ($(OS),Windows_NT)
	@if exist "C:\Program Files\$(APP_NAME)\$(APP_NAME).exe" del "C:\Program Files\$(APP_NAME)\$(APP_NAME).exe"
else
	@sudo rm -f /usr/local/bin/$(APP_NAME)
endif
	@echo Uninstall complete

# Run the application
run: build
	@echo Running $(APP_NAME)...
	@$(BUILD_DIR)\$(APP_NAME).exe

# Download dependencies
deps:
	@echo Downloading dependencies...
	$(GOGET) -v -t -d ./...

# Tidy dependencies
tidy:
	@echo Tidying dependencies...
	$(GOMOD) tidy

# Help target
help:
	@echo SNIP - Code Snippet Manager
	@echo.
	@echo Available targets:
	@echo   make build        - Build for current platform
	@echo   make build-all    - Build for all platforms
	@echo   make clean        - Remove build artifacts
	@echo   make test         - Run tests
	@echo   make coverage     - Run tests with coverage report
	@echo   make fmt          - Format code
	@echo   make vet          - Run go vet
	@echo   make lint         - Run golangci-lint
	@echo   make install      - Install binary to system
	@echo   make uninstall    - Remove binary from system
	@echo   make run          - Build and run the application
	@echo   make deps         - Download dependencies
	@echo   make tidy         - Tidy go.mod and go.sum
	@echo   make all          - Run clean, fmt, vet, test, and build
	@echo   make help         - Show this help message
	@echo.
	@echo Detected OS: $(DETECTED_OS)
