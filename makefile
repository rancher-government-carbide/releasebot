.PHONY: dependencies test lint container container-push linux darwin windows clean help 

BINARY_NAME=releasebot
CONTAINER_NAME=clanktron/releasebot
SRC=$(shell git ls-files ./cmd)
VERSION=0.1.0
COMMIT_HASH=$(shell git rev-parse HEAD)
GOENV=GOARCH=amd64 CGO_ENABLED=0
BUILD_FLAGS=-ldflags="-X 'main.Version=$(VERSION)'"
TEST_FLAGS=-v -cover -count 2
CONTAINER_CLI=nerdctl

# Build the binary
$(BINARY_NAME):
	$(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(SRC)

dependencies:
	go mod tidy && go get -v -d ./...

# Test the binary
test: $(BINARY_NAME)
	go test $(TEST_FLAGS) $(SRC) 
	make clean

# Run linters
lint: $(BINARY_NAME)
	go vet $(SRC)
	staticcheck $(SRC)
	make clean

# Build the container image
container: clean
	$(CONTAINER_CLI) build -t $(CONTAINER_NAME):$(COMMIT_HASH) -t $(CONTAINER_NAME):latest .
	
# Push the binary
container-push: container
	$(CONTAINER_CLI) push $(CONTAINER_NAME):$(COMMIT_HASH) && $(CONTAINER_CLI) push $(CONTAINER_NAME):latest

# Build the binary for Linux
linux:
	GOOS=linux $(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux $(SRC)
# Build the binary for MacOS
darwin:
	GOOS=darwin $(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin $(SRC)
# Build the binary for Windows
windows:
	GOOS=windows $(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME)-windows $(SRC)

# Clean the binary
clean:
	rm -f $(BINARY_NAME)

# Show help
help:
	@printf "Available targets:\n"
	@printf "  $(BINARY_NAME) 		Build the binary\n"
	@printf "  test 			Build and test the binary\n"
	@printf "  linux 		Build the binary for Linux\n"
	@printf "  darwin 		Build the binary for MacOS\n"
	@printf "  windows 		Build the binary for Windows\n"
	@printf "  container 		Build the container\n"
	@printf "  container-push 	Build and push the container\n"
	@printf "  env 			apply .env file in PWD\n"
	@printf "  clean 		Clean build results\n"
	@printf "  help 			Show help\n"
