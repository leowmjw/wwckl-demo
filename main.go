//usr/bin/env test -x $0 && (([ ! -x "${0}c" ] || [ "$0" -nt "${0}c" ]) && go build -o "${0}c" ./...; "${0}c" $@); exit "$?"
// Great snippet from dsypniewski: https://gist.github.com/posener/73ffd326d88483df6b1cb66e8ed1e0bd; just chmod +x and run!  Works for multi files too ..
package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"

	"github.com/leowmjw/pdfcpu/pkg/api"
	"github.com/leowmjw/pdfcpu/pkg/pdfcpu"
)

func main() {
	log.Println("WWCKL Demo!!!")

	// Open up sample PDFs into byte stream ioutil?

	// Extract out each page one by one
	// Inside each page; determine the page boundaries; extract out initial metadata; question numbers into
	//	simple topic structure
	// Split all into individual pages into workspace area
	// Reform page per topic based on simple topic structure
	// Attempt to extract out more structure based on type information recognized ..

	sourceFileName := "./extract-selangor-gov-my/samples/Selangor-Mulut-1-20.pdf"
	// sourceFileName := "./extract-selangor-gov-my/samples/Selangor-Tulis-1-20.pdf"
	// sourceFileName := "./extract-selangor-gov-my/samples/Selangor-Penyata-JP-PBT.pdf"
	// sourceFileName := "./extract-selangor-gov-my/samples/Selangor-Maklumbalas.pdf"
	pdfctx, readerr := pdfcpu.ReadPDFFile(sourceFileName, pdfcpu.NewDefaultConfiguration())
	if readerr != nil {
		log.Fatal("ERR:", readerr)
	}
	// Needs to verify first otherwise page count is not in there ..
	valerr := pdfcpu.ValidateXRefTable(pdfctx.XRefTable)
	if valerr != nil {
		log.Fatal("val_ERR: ", valerr)
	}
	log.Println("Document has ", pdfctx.PageCount, " page(s)")
	pdfref, pgerr := pdfctx.Pages()
	if pgerr != nil {
		log.Fatal("ERR:", pgerr)
	}
	spew.Dump(pdfref)
	log.Println("Name:", pdfctx.Read.FileName, " Size:", pdfctx.Read.FileSize)
	// data, exerr := pdfcpu.ExtractContentData(pdfctx, 0)
	spew.Println(api.ParsePageSelection("1-50"))
	// spew.Dump(pdfcpu.ExtractContentData(pdfctx, 1))

	pageSelection := pdfcpu.IntSet{}
	pageSelection[1] = true
	pageSelection[2] = true
	exerr := exploreContent(pdfctx, pageSelection)
	if exerr != nil {
		log.Fatal("explore_ERR: ", exerr)
	}

	// Experiment #2: Use rsc/pdf, vs unidoc
	// Not so good; and have watermark license ..
	// exploreContentWithUnidoc(sourceFileName)
	exploreContentWithRSCPdf(sourceFileName)
	rotatedFileName := "./extract-selangor-gov-my/samples/Selangor-Penyata-JP-PBT.pdf"
	exploreRotatedWithUnidoc(rotatedFileName)
	// Experiment #3: go-fitz vs docconv
	// Fitx is pretty good
	// exploreContentWithFitz(sourceFileName)
}
