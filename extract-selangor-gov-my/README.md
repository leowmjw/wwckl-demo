# Extract structured text from Selangor Gov (dewan.selangor.gov.my)

## Pre-reqs

- brew install imagemagick
- brew install tesseract --with-all-languages 
- brew install ghostscript (<-- gives you gs)
- brew install pdfsandwich (<-- improve detect non-OCR)
- pipenv install pdfminer.six

## Dependencies

- OSX - brew install poppler 
- invoice2data

## Setup

- While in this folder; invoke "pipenv shell" to get latest with the needed deps in a clean shell ..

## Other Tools available

- Tabula: https://tabula.technology/
- Tabula Java lib (including runnable jar): https://github.com/tabulapdf/tabula-java
- Tabula Python Pandas: https://github.com/chezou/tabula-py
- iTextPDF: https://itextpdf.com/
- PDFCpu - https://godoc.org/github.com/hhrutter/pdfcpu <-- use this to split up PDFs and perform other processing?
- https://github.com/UW-Deepdive-Infrastructure/table-extract
- https://github.com/ashima/pdf-table-extract
- If all else fails - https://github.com/WZBSocialScienceCenter/pdftabextract

## Running

Extract out the matching ..
```bash
invoice2data  --debug ./samples/Selangor-Mulut-1-20.pdf 
```

Extract out tables automatically ..
```bash
 java  -Dsun.java2d.cmm=sun.java2d.cmm.kcms.KcmsServiceProvider -jar ./bin/tabula.jar --lattice --guess --pages all \
    -o /tmp/myouput.csv ./samples/Selangor-Mulut-1-20.pdf 
```

Convert pure image PDF to OCR-ed PDF
```bash
pdfsandwich ./samples/Selangor-Mulut-1-20.pdf 2>&1 | less
```

Optimize
```bash
pdfcpu optimize -stats /tmp/stats.csv  ./samples/Selangor-Mulut-1-20.pdf \ 
    ./samples/Selangor-Mulut-1-20-optimized.pdf
```

Extract out images
```bash
# pdfcpu extract -mode image ./samples/Selangor-Mulut-1-20.pdf /tmp/pdfcpu
```

Split out individual
```bash
# pdfcpu extract -mode page ./samples/Selangor-Mulut-1-20.pdf /tmp/pdfcpu

```

## Articles

- https://datascience.blog.wzb.eu/2016/07/04/data-mining-pdfs-the-simple-cases/
- https://datascience.blog.wzb.eu/2016/07/08/data-mining-ocr-pdfs-getting-things-straight/#more-96
- https://github.com/UW-Deepdive-Infrastructure/blackstack
- https://datascience.blog.wzb.eu/2017/02/16/data-mining-ocr-pdfs-using-pdftabextract-to-liberate-tabular-data-from-scanned-documents/

## Research Papers

- http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.638.7400&rep=rep1&type=pdf