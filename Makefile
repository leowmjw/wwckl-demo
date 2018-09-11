test:
	go test ./...

build:
	go build ./...

deps:
	go get -u github.com/davecgh/go-spew/spew

tools:
	# Get tabula-jar into the extractor bin? --> https://github.com/tabulapdf/tabula-java
	# Get pdfminer script into bin?
	# Get pdfcpu tooling into bin ..
	go get github.com/hhrutter/pdfcpu/cmd/...

.PHONY: build