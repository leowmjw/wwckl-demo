build:
	go build ./... && ./wwckl-demo

test:
	go test ./...

deps:
	go get github.com/davecgh/go-spew/spew
	go get github.com/hhrutter/pdfcpu/pkg/pdfcpu
	go get github.com/hhrutter/pdfcpu/pkg/api
	# Try out alternatives like unidoc, docconv,pdf ..
	go get github.com/ledongthuc/pdf
	# Using MuPDF engine ..
	go get github.com/gen2brain/go-fitz

tools:
	# Get tabula-jar into the extractor bin? --> https://github.com/tabulapdf/tabula-java
	# Get pdfminer script into bin?
	# Get pdfcpu tooling into bin ..
	go get github.com/hhrutter/pdfcpu/cmd/...
	# Needs license to create PDFs; but extract text?
	go get github.com/unidoc/unidoc/...
	# Too much complexity??
	go get code.sajari.com/docconv/...

.PHONY: build
