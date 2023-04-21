.PHONY: build build-linux build-darwin build-windows run clean help

BINARY_NAME=releasebot

# Set the source files
SRC=./src/*

# Set the version
VERSION=0.1.0

# Set the build flags
BUILD_FLAGS=-ldflags="-X 'main.Version=$(VERSION)'"

# Build the binary
build:
	go build $(BUILD_FLAGS) -o ./build/$(BINARY_NAME) $(SRC)

# Build the binary for Linux
build-linux:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux $(SRC)
# Build the binary for MacOS
build-darwin:
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin $(SRC)
# Build the binary for Windows
build-windows:
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-windows $(SRC)

# run the executable
run:
	./build/$(BINARY_NAME)

# Clean the binary
clean:
	rm -rf build/*

# Run tests
# test:
#	go test ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  build           Build the binary"
	@echo "  build-linux     Build the binary for Linux"
	@echo "  build-darwin    Build the binary for MacOS"
	@echo "  build-windows   Build the binary for Windows"
	@echo "  clean           Clean the binary"
	@echo "  help            Show help"
