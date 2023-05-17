.PHONY: dependencies test container container-push linux darwin windows env run clean help 

BINARY_NAME=releasebot
CONTAINERTAG=clanktron/releasebot
SRC=$(shell git ls-files ./cmd)
VERSION=0.1.0
GOENV=GOARCH=amd64 CGO_ENABLED=0
BUILD_FLAGS=-ldflags="-X 'main.Version=$(VERSION)'"

# Build the binary
releasebot: clean
	$(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(SRC)

dependencies:
	go mod tidy && go get -v -d ./...

# Test the binary
test: releasebot
	go test $(SRC)

# Build the container image
container:
	docker build -t $(CONTAINERTAG):$(VERSION) . && docker image tag $(CONTAINERTAG):$(VERSION) $(CONTAINERTAG):latest
	
# Push the binary
container-push: container
	docker push $(CONTAINERTAG):$(VERSION) && docker push $(CONTAINERTAG):latest 

# Build the binary for Linux
linux:
	GOOS=linux $(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux $(SRC)
# Build the binary for MacOS
darwin:
	GOOS=darwin $(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin $(SRC)
# Build the binary for Windows
windows:
	GOOS=windows $(GOENV) go build $(BUILD_FLAGS) -o $(BINARY_NAME)-windows $(SRC)

env:
	set -a; source .env; set +a

# Clean the binary
clean:
	rm -f $(BINARY_NAME)

# Show help
help:
	@printf "Available targets:\n"
	@printf "  releasebot 		Build the binary\n"
	@printf "  test 			Build and test the binary\n"
	@printf "  linux 		Build the binary for Linux\n"
	@printf "  darwin 		Build the binary for MacOS\n"
	@printf "  windows 		Build the binary for Windows\n"
	@printf "  container 		Build the container\n"
	@printf "  container-push 	Build and push the container\n"
	@printf "  env 			apply .env file in PWD\n"
	@printf "  clean 		Clean build results\n"
	@printf "  help 			Show help\n"
