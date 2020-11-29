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
	rm juryrig
inspect: build
	golint ./...
update:
	go get -u ./...
pre-commit: update clean coverage.txt inspect
	go mod tidy
clean:
	rm -f juryrig
	rm -f coverage.txt