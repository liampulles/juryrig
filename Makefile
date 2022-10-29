# Init variables
GOBIN := $(shell go env GOPATH)/bin

# Keep test at the top so that it is default when `make` is called.
# This is used by Travis CI.
coverage.txt:
	go test -race -coverprofile=coverage.txt -covermode=atomic -coverpkg=./pkg/...,./ ./...
view-cover: clean coverage.txt
	go tool cover -html=coverage.txt
test: build
	go test ./test/...
build:
	go build ./...
install: build
	go install ./...
inspect: build $(GOBIN)/golangci-lint
	golangci-lint run
update:
	go get -u ./...
pre-commit: update clean coverage.txt inspect
	go mod tidy
clean:
	rm -f coverage.txt

# Needed tools
GOLANGCI_VERSION := 1.50.1

$(GOBIN)/golangci-lint:
	$(MAKE) install-tools
install-tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v$(GOLANGCI_VERSION)
	rm -rf ./v$(GOLANGCI_VERSION)