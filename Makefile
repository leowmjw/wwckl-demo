test:
	go test ./...

build:
	go build ./... && ./wwckl-demo

deps:
	go get github.com/davecgh/go-spew/spew
	go get github.com/hhrutter/pdfcpu/pkg/pdfcpu
	go get github.com/hhrutter/pdfcpu/pkg/api

tools:
	# Get tabula-jar into the extractor bin? --> https://github.com/tabulapdf/tabula-java
	# Get pdfminer script into bin?
	# Get pdfcpu tooling into bin ..
	go get github.com/hhrutter/pdfcpu/cmd/...

.PHONY: build
