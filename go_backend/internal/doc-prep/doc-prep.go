package docprep

import (
	"VLN-backend/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dslipak/pdf"
)

type KeywordExtractionRequest struct {
	RawText string `json:"raw_text"`
}

// PDF Reader adapted from https://pkg.go.dev/github.com/dslipak/pdf#section-readme
func readPdf(path string) (*string, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return nil, err
	}
	buf.ReadFrom(b)
	result := buf.String()
	return &result, nil
}

func ExtractKeywords() (*string, error) {
	//Read in PDF and process to raw text (string)
	//pdfData: String
	pdfData, err := readPdf("testPDF.pdf")
	if err != nil {
		return nil, err
	}

	requestBody, err := json.Marshal(KeywordExtractionRequest{RawText: *pdfData})
	if err != nil {
		return nil, err
	}

	keywordsAPI := config.GetMLService() + "/api/keyword_extractor"
	resp, err := http.Post(keywordsAPI, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	preprocessedText := string(respBody)

	fmt.Println("PREPROCESSED TEXT: ")
	fmt.Println(preprocessedText)
	fmt.Println("END PROCESSED TEXT")
	return &preprocessedText, nil
}
