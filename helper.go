package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"log"

	"github.com/leowmjw/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
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
