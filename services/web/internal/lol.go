package internal

import (
	"bytes"
	"fmt"
	"log"

	"github.com/dslipak/pdf"
)

func ReadFile() {
	// Opens a PDF file and extracts text to a reader
	r, err := pdf.Open("example.pdf")
	if err != nil {
		log.Fatal(err)
	}

	// GetPlainText returns a reader containing the text
	b, err := r.GetPlainText()
	if err != nil {
		log.Fatal(err)
	}

	// Read into a buffer to output
	var buf bytes.Buffer
	buf.ReadFrom(b)
	fmt.Println(buf.String())
}
