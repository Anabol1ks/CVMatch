package parser

import (
	"bytes"

	"github.com/ledongthuc/pdf"
)

func ExtractTextFromPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var textBuilder bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	_, err = textBuilder.ReadFrom(b)
	if err != nil {
		return "", err
	}
	return textBuilder.String(), nil
}
