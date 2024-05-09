package scraper

import (
	"VLN-backend/config"
	googlescraper "VLN-backend/internal/google-scraper"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type similarityRequest struct {
	Prompt string   `json:"prompt"`
	Images []string `json:"images"`
}

type similarityResponse struct {
	Score            float64 `json:"score"`
	ImageURL         string  `json:"url"`
	ImageDescription string  `json:"description"`
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

type imageScore struct {
	SimilarityScore float64 `json:"similarity_score"`
	QualityScore    float64 `json:"quality_score"`
	ImagePath       string  `json:"image_path"`
}

func getSimScores(images []*googlescraper.Image, q_keyword string) ([]similarityResponse, error) {
	/*
		Calls the ML process_similarity to a list of {keyword and list of images} --> get the text vector similarity between
		our query keyword and the downloaded images' caption

		Stores the result as a list of similarityResponse
	*/

	// Call api/process_similarity
	similarity_api := config.GetMLService() + "/api/process_similarity"
	image_urls := make([]string, 0)
	for _, image := range images {
		if image == nil {
			continue
		}
		image_urls = append(image_urls, image.ImagePath)
	}
	request_body, _ := json.Marshal((similarityRequest{
		Prompt: q_keyword,
		Images: image_urls,
	}))

	resp, err := http.Post(similarity_api, "application/json", bytes.NewBuffer(request_body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Process VLNML Response (api return a list of submitted taskS metadata)
	var vlnml_res []VLNMLResponse
	vlnml_resp_format_err := json.Unmarshal([]byte(string(resp_body)), &vlnml_res)
	if vlnml_resp_format_err != nil {
		log.Println("ERROR when trying to map inital VLNML response", err)
		return nil, err
	}

	for _, res := range vlnml_res {
		task_id := res.TaskID
		log.Println("TASK_ID = ", task_id)

		// Start an interval to check the ML status every 2 seconds
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			status_api := config.GetMLService() + "/api/status/" + task_id
			statusResp, err := http.Get(status_api)
			if err != nil {
				return nil, err
			}
			defer statusResp.Body.Close()

			statusBody, err := io.ReadAll(statusResp.Body)
			if err != nil {
				return nil, err
			}

			var statusResponse StatusResponse
			if err := json.Unmarshal(statusBody, &statusResponse); err != nil {
				return nil, err
			}

			if statusResponse.Status == "SUCCESS" {
				result_api := config.GetMLService() + "/api/result/" + task_id
				resultResp, err := http.Get(result_api)
				if err != nil {
					return nil, err
				}
				defer resultResp.Body.Close()

				resultBody, err := io.ReadAll(resultResp.Body)

				var similarity_scores []similarityResponse

				json_format_err := json.Unmarshal([]byte(string(resultBody)), &similarity_scores)
				if json_format_err != nil {
					return nil, err
				}
				return similarity_scores, nil
			}

			if statusResponse.Status == "FAIL" {
				return nil, errors.New("task failed")
			}
		}
	}
	return nil, nil
}

func getQualityScore(width float64, height float64) float64 {
	// function that caculates the quality scored based on the image height width
	// todo needs implementation

	return width * height / 100
}

func getImageScores(sim_scores []similarityResponse, images []*googlescraper.Image) []imageScore {
	//format the images into a list of imageScore for sorting

	image_scores := make([]imageScore, 0)
	for _, sim_resp := range sim_scores {
		image_path := sim_resp.ImageURL
		for _, image := range images {
			if image.ImagePath == image_path {
				image_scores = append(image_scores, imageScore{
					SimilarityScore: sim_resp.Score,
					QualityScore:    getQualityScore(image.Width, image.Height),
					ImagePath:       image.ImagePath,
				})
			}
		}
	}
	return image_scores
}

func filterImages(images []*googlescraper.Image, q_keyword string) ([]imageScore, error) {
	similarity_scores, err := getSimScores(images, q_keyword)
	if err != nil {
		return nil, err
	}

	image_scores := getImageScores(similarity_scores, images)

	//todo the weighting function needs implementation
	sort.Slice(image_scores, func(i, j int) bool {
		s1, s2 := image_scores[i], image_scores[j]
		return s2.QualityScore*s2.SimilarityScore < s1.QualityScore*s1.SimilarityScore
	})
	return image_scores, nil
}
func sanitizeQueryString(query string) string {
	query = strings.TrimSpace(query)
	query = strings.ReplaceAll(query, " ", "_")

	illegalChars := []string{`\`, `/`, `:`, `*`, `?`, `"`, `<`, `>`, `|`}
	for _, char := range illegalChars {
		query = strings.ReplaceAll(query, char, "_")
	}
	if len(query) > 255 {
		query = query[:255]
	}

	return query
}
func GetImages(c *gin.Context) {
	/*
	   request format in params:
	   {
	       query: query key word for google
	   }

	   response: returns the top image for the frame based on the scraped result
	*/

	// query, key := c.GetQuery("query")
	// if !key {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"Error: ": "Invalid request parameter for getting image",
	// 	})
	// 	c.Abort()
	// }

	//////WHAT MINH ADDED TO SCRAPER/////////////////////
	// query, err := docprep.ExtractKeywords()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Error while extracting keywords",
	// 	})
	// 	return
	// }
	/////////////////////////////////////////////////////
	query, key := c.GetQuery("query")
	if !key {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error: ": "Invalid request parameter for getting image",
		})
		c.Abort()
	}
	log.Println("Requested data: ", query)

	images, err := googlescraper.ImageSearch(query)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"Error: ": "Error while getting image from google",
		})
		c.Abort()
	}

	images, e := googlescraper.DownloadImages(sanitizeQueryString(query), images[:10])
	if e != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"Status: ": "Error while downloading image from google",
			"Error: ":  e,
		})
		c.Abort()
	}

	image_scores, err := filterImages(images, query)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"Status: ": "Error while filtering downloaded images",
			"Error":    err,
		})
		c.Abort()
	}

	log.Println("[", query, "] Image scores are: ", image_scores)

	c.JSON(http.StatusOK, gin.H{
		"image_scores": image_scores,
		"status":       http.StatusOK,
	})
}
