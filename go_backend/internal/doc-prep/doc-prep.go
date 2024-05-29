package docprep

import (
	"VLN-backend/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dslipak/pdf"
)

type KeywordExtractionRequest struct {
	RawText string `json:"raw_text"`
}
type VLNMLResponse struct {
	TaskID    string `json:"task_id"`
	TaskName  string `json:"task_name"`
	Status    string `json:"status"`
	URLResult string `json:"url_result"`
}
type StatusResponse struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
}
type KeywordResponse struct {
	Keywords  []string `json:"sentences"`
	Sentences []string `json:"keywords"`
}

// PDF Reader adapted from https://pkg.go.dev/github.com/dslipak/pdf#section-readme
func readPdf(path string) (*string, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF: %w", err)
	}

	log.Println("Starting extraction...")
	var buf bytes.Buffer

	// Get the plain text from the PDF.
	b, err := r.GetPlainText()
	if err != nil {
		return nil, fmt.Errorf("failed to get plain text from PDF: %w", err)
	}

	// Read the plain text into the buffer.
	_, err = buf.ReadFrom(b)
	if err != nil {
		return nil, fmt.Errorf("failed to read from buffer: %w", err)
	}

	result := buf.String()
	return &result, nil
}

func ExtractKeywords(file string) (KeywordResponse, error) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
	}

	filepath := cwd + "/" + "uploads/" + file

	// Check if the file exists before attempting to read it.
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Fatalf("file does not exist: %s", filepath)
	} else {
		log.Printf("Reading PDF file from: %s", filepath)

	}

	// Read the PDF file.
	pdfData, err := readPdf(filepath)
	if err != nil {
		return KeywordResponse{}, err
	}

	requestBody, err := json.Marshal(KeywordExtractionRequest{RawText: *pdfData})
	if err != nil {
		return KeywordResponse{}, err
	}

	keywordsAPI := config.GetMLService() + "/api/process_keywords"
	resp, err := http.Post(keywordsAPI, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return KeywordResponse{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return KeywordResponse{}, err
	}

	var vlnml_res []VLNMLResponse
	vlnml_resp_format_err := json.Unmarshal([]byte(string(respBody)), &vlnml_res)
	if vlnml_resp_format_err != nil {
		log.Println("ERROR when trying to map inital VLNML response", err)
		return KeywordResponse{}, err
	}

	for _, res := range vlnml_res {
		task_id := res.TaskID
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			status_api := config.GetMLService() + "/api/status/" + task_id
			statusResp, err := http.Get(status_api)
			if err != nil {
				return KeywordResponse{}, err
			}
			defer statusResp.Body.Close()

			statusBody, err := io.ReadAll(statusResp.Body)
			if err != nil {
				return KeywordResponse{}, err
			}

			var statusResponse StatusResponse
			if err := json.Unmarshal(statusBody, &statusResponse); err != nil {
				return KeywordResponse{}, err
			}

			if statusResponse.Status == "SUCCESS" {
				result_api := config.GetMLService() + "/api/result/" + task_id
				resultResp, err := http.Get(result_api)
				if err != nil {
					return KeywordResponse{}, err
				}
				defer resultResp.Body.Close()

				resultBody, err := io.ReadAll(resultResp.Body)

				log.Println("resulBody = ", string(resultBody))
				var keywords KeywordResponse

				json_format_err := json.Unmarshal([]byte(string(resultBody)), &keywords)
				if json_format_err != nil {
					log.Println("ERROR while maping the result value", json_format_err)
					return KeywordResponse{}, json_format_err
				}
				return keywords, nil
			}

			if statusResponse.Status == "FAIL" {
				return KeywordResponse{}, errors.New("task failed")
			}
			// else, the ticker loop continue
		}
	}
	return KeywordResponse{}, nil
}
