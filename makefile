.PHONY: build container linux darwin windows run clean help 

BINARY_NAME=releasebot
CONTAINERTAG=clanktron/releasebot
SRC=./src/*
VERSION=0.1.0
BUILD_FLAGS=-ldflags="-X 'main.Version=$(VERSION)'"

# Build the binary
build-binary:
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(SRC)

# Build the binary
container:
	docker build -t $(CONTAINERTAG) .

# Build the binary for Linux
linux:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux $(SRC)
# Build the binary for MacOS
darwin:
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin $(SRC)
# Build the binary for Windows
windows:
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-windows $(SRC)

# run the executable
run:
	./$(BINARY_NAME)

# Clean the binary
clean:
	rm -f $(BINARY_NAME)

# Run tests
# test:
#	go test ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  build           Build the binary"
	@echo "  linux     	     Build the binary for Linux"
	@echo "  darwin    	     Build the binary for MacOS"
	@echo "  windows   	     Build the binary for Windows"
	@echo "  container   	 Build the container
	@echo "  clean           Clean the binary"
	@echo "  help            Show help"
