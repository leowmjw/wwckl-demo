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
- 

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