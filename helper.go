package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/gen2brain/go-fitz"

	"log"

	"github.com/leowmjw/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"

	"github.com/unidoc/unidoc/pdf/extractor"
	unidocpdf "github.com/unidoc/unidoc/pdf/model"

	rscpdf "github.com/ledongthuc/pdf"
)

// All extracted from pdfcpu .. da best!
func contentObjNrs(ctx *pdfcpu.PDFContext, page int) ([]int, error) {

	objNrs := []int{}

	log.Println("PAGE: ", page)
	d, _, err := ctx.PageDict(page)
	if err != nil {
		return nil, err
	}

	obj, found := d.Find("Contents")
	if !found || obj == nil {
		return nil, nil
	}

	var objNr int

	indRef, ok := obj.(pdfcpu.PDFIndirectRef)
	if ok {
		objNr = indRef.ObjectNumber.Value()
	}

	obj, err = ctx.Dereference(obj)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, nil
	}

	switch obj := obj.(type) {

	case pdfcpu.PDFStreamDict:

		objNrs = append(objNrs, objNr)

	case pdfcpu.PDFArray:

		for _, obj := range obj {

			indRef, ok := obj.(pdfcpu.PDFIndirectRef)
			if !ok {
				return nil, errors.Errorf("missing indref for page tree dict content no page %d", page)
			}

			sd, err := ctx.DereferenceStreamDict(obj)
			if err != nil {
				return nil, err
			}

			if sd == nil {
				continue
			}

			objNrs = append(objNrs, indRef.ObjectNumber.Value())

		}

	}

	return objNrs, nil
}

func exploreContent(ctx *pdfcpu.PDFContext, selectedPages pdfcpu.IntSet) error {

	visited := pdfcpu.IntSet{}

	for p, v := range selectedPages {

		// Page has been chosen for exploration ..
		if v {
			objNrs, err := contentObjNrs(ctx, p)
			if err != nil {
				log.Fatal("context_ERR for page:", p)
				return err
			}

			if objNrs == nil {
				log.Println("objNrs is NIL!!")
				continue
			}

			for _, objNr := range objNrs {

				if visited[objNr] {
					log.Println("VISITED BEFOREE:", objNr)
					continue
				}

				visited[objNr] = true

				b, err := pdfcpu.ExtractContentData(ctx, objNr)
				if err != nil {
					log.Fatal("EXTRACT_ERR:", err)
					return err
				}

				if b == nil {
					log.Println("Nothing to do with: ", objNr)
					continue
				}

				// log.Println(string(b[:]))
			}
		}
	}
	return nil
}

func doExtractContent(ctx *pdfcpu.PDFContext, selectedPages pdfcpu.IntSet) error {

	visited := pdfcpu.IntSet{}

	for p, v := range selectedPages {

		if v {

			// log.Info.Printf("writing content for page %d\n", p)

			objNrs, err := contentObjNrs(ctx, p)
			if err != nil {
				return err
			}

			if objNrs == nil {
				continue
			}

			for _, objNr := range objNrs {

				if visited[objNr] {
					continue
				}

				visited[objNr] = true

				b, err := pdfcpu.ExtractContentData(ctx, objNr)
				if err != nil {
					return err
				}

				if b == nil {
					continue
				}

				fileName := fmt.Sprintf("%s/%d_%d.txt", ctx.Write.DirName, p, objNr)

				err = ioutil.WriteFile(fileName, b, os.ModePerm)
				if err != nil {
					return err
				}

			}

		}

	}

	return nil
}

func exploreContentWithFitz(fileName string) error {

	doc, err := fitz.New(fileName)
	if err != nil {
		log.Fatal("Fitz_ERR:", err)
	}
	defer doc.Close()

	log.Println("Number of pages: ", doc.NumPage())
	for i := 0; i < 3; i++ {
		pageText, exerr := doc.Text(i)
		if exerr != nil {
			log.Fatal("exText_ERR: ", exerr)
		}
		log.Println(pageText)
	}

	return nil
}

func exploreContentWithUnidoc(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := unidocpdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	numPage, err := pdfReader.GetNumPages()
	if err != nil {
		log.Fatal("pgNum_ERR:", err)
	}
	log.Println("Document has ", numPage, " page(s)")
	for i := 0; i < 3; i++ {
		pdfPage, err := pdfReader.GetPage(i + 1)
		if err != nil {
			log.Fatal("getPage_ERR: ", err)
		}
		// Below for structure ..
		// spew.Dump(pdfPage.Contents)
		// perr := processPage(pdfPage)
		// if perr != nil {
		// 	log.Fatal("process_ERR: ", perr)
		// }

		ex, err := extractor.New(pdfPage)
		if err != nil {
			return err
		}

		text, err := ex.ExtractText()
		if err != nil {
			return err
		}

		fmt.Println("------------------------------")
		fmt.Printf("Page %d:\n", i+1)
		fmt.Printf("\"%s\"\n", text)
		fmt.Println("------------------------------")

	}
	return nil
}

func exploreRotatedWithUnidoc(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := unidocpdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	numPage, err := pdfReader.GetNumPages()
	if err != nil {
		log.Fatal("pgNum_ERR:", err)
	}
	log.Println("Document has ", numPage, " page(s)")
	for i := 0; i < 3; i++ {
		pdfPage, err := pdfReader.GetPage(i + 1)
		if err != nil {
			log.Fatal("getPage_ERR: ", err)
		}
		// Below for structure ..
		spew.Dump(pdfPage.Contents)
		perr := processPage(pdfPage)
		if perr != nil {
			log.Fatal("process_ERR: ", perr)
		}

	}
	return nil
}

func processPage(page *unidocpdf.PdfPage) error {
	mBox, err := page.GetMediaBox()
	if err != nil {
		return err
	}
	pageWidth := mBox.Urx - mBox.Llx
	pageHeight := mBox.Ury - mBox.Lly

	fmt.Printf(" Page: %+v\n", page)
	if page.Rotate != nil {
		fmt.Printf(" Page rotation: %v\n", *page.Rotate)
	} else {
		fmt.Printf(" Page rotation: 0\n")
	}
	fmt.Printf(" Page mediabox: %+v\n", page.MediaBox)
	fmt.Printf(" Page height: %f\n", pageHeight)
	fmt.Printf(" Page width: %f\n", pageWidth)

	return nil
}

func exploreContentWithRSCPdf(fileName string) error {
	f, r, err := rscpdf.Open(fileName)
	if err != nil {
		log.Fatal("RSCOpen_ERR: ", err)
	}
	defer f.Close()

	// Below to read mulitple ..
	// var buf bytes.Buffer

	// b, err := r.GetPlainText()
	// if err != nil {
	// 	log.Fatal("PlainTXT_ERR: ", err)
	// }

	// n, err := buf.ReadFrom(b)
	// if err != nil {
	// 	log.Fatal("ReadFrm_ERR:", err)
	// }
	// log.Println("Read number of ", n)
	// log.Println(buf.String())

	numPage := r.NumPage()
	log.Println("Document has ", numPage, " page(s)")
	var textBuilder bytes.Buffer
	for i := 0; i < 3; i++ {
		// Must be page number; not zero indexed .. fireezes ..
		p := r.Page(i + 1)
		if p.V.IsNull() {
			continue
		}
		lineStr, err := p.GetPlainText(nil)
		if err != nil {
			log.Fatal("GetTXT_ERR:", err)
		}
		textBuilder.WriteString(lineStr)
	}

	log.Println("+++++++++++++++++++++++++ INDIVIDUAL ++++++++++++++++++")
	log.Println(textBuilder.String())

	log.Println("======================= STYLED ========================")

	for pageIndex := 1; pageIndex < 3; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		// spew.Dump(p)
		log.Println("PAGE: ", pageIndex)
		var lastTextStyle rscpdf.Text
		texts := p.Content().Text
		// If check empty for texts; the klen should show it does not have any text; needs OCR ..
		for _, text := range texts {
			if isSameSentence(text.S, lastTextStyle.S) {
				lastTextStyle.S = lastTextStyle.S + text.S
			} else {
				fmt.Printf("Font: %s, Font-size: %f, x: %f, y: %f, content: %s \n", lastTextStyle.Font, lastTextStyle.FontSize, lastTextStyle.X, lastTextStyle.Y, lastTextStyle.S)
				lastTextStyle = text
			}
		}
		log.Println(lastTextStyle.S)
	}

	return nil
}

func isSameSentence(currentTxt string, ongoingTxt string) bool {
	// If previous non-space; is '.' will end the sentence ..
	// log.Println("NEW:", currentTxt, " ONGOING: ", ongoingTxt)
	return true
}
