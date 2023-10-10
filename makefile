.PHONY: check test lint container container-push dependencies clean help

BINARY_NAME=releasebot
CONTAINER_NAME=clanktron/releasebot
SRC=./cmd
VERSION=0.1.2
COMMIT_HASH=$(shell git rev-parse HEAD)
GOENV=CGO_ENABLED=0
BUILD_FLAGS=-ldflags="-X 'main.Version=$(VERSION)'"
TEST_FLAGS=-v -cover -count 1
CONTAINER_CLI=nerdctl
DATA_FOLDER=./data

# Build the binary
$(BINARY_NAME):
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(SRC)

check: test lint

# Test the binary
test:
	go test $(TEST_FLAGS) $(SRC) 

# Run linters
lint:
	go vet $(SRC)
	staticcheck $(SRC)

# Build the container image
container: clean
	$(CONTAINER_CLI) build -t $(CONTAINER_NAME):$(COMMIT_HASH) -t $(CONTAINER_NAME):latest .
	
# Push the binary
container-push: container
	$(CONTAINER_CLI) push $(CONTAINER_NAME):$(COMMIT_HASH) && $(CONTAINER_CLI) push $(CONTAINER_NAME):latest

dependencies:
	go mod tidy && go get -v -d ./...

# Clean the binary
clean:
	rm -rf $(BINARY_NAME) $(DATA_FOLDER)

# Show help
help:
	@printf "Available targets:\n"
	@printf "  $(BINARY_NAME) 		Build the binary (default)\n"
	@printf "  test 			Run all unit tests\n"
	@printf "  lint 			Run go vet and staticcheck\n"
	@printf "  check 		Test and lint the binary\n"
	@printf "  container 		Build the container\n"
	@printf "  container-push 	Build and push the container\n"
	@printf "  clean 		Clean build results and data folder\n"
	@printf "  help 			Show help\n"
