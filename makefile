.PHONY: container test linux darwin windows env run clean help 

BINARY_NAME=releasebot
CONTAINERTAG=clanktron/releasebot
SRC=$(shell git ls-files ./cmd)
VERSION=0.1.0
BUILD_FLAGS=-ldflags="-X 'main.Version=$(VERSION)'"

# Build the binary
releasebot:
	go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(SRC)

# Test the binary
test: releasebot
	go test $(SRC)

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

env:
	set -a; source .env; set +a

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
	@printf "Available targets:\n"
	@printf "  build 		Build the binary\n"
	@printf "  linux 		Build the binary for Linux\n"
	@printf "  darwin 		Build the binary for MacOS\n"
	@printf "  windows 		Build the binary for Windows\n"
	@printf "  container 		Build the container\n"
	@printf "  clean 		Clean the binary\n"
	@printf "  help 			Show help\n"
