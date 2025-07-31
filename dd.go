package main

import (
	"CVMatch/internal/parser"
	"log"
)

func main() {
	text, err := parser.ExtractTextFromPDF("./uploads/resume.pdf")
	if err != nil {
		log.Fatalf("Error extracting text: %v", err)
	}
	log.Println(text)

}
