test:
	go test ./...

build:
	go build ./...

deps:
	go get -u github.com/davecgh/go-spew/spew

.PHONY: build