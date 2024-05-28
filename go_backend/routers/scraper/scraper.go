package scraper

import (
	"VLN-backend/config"
	docprep "VLN-backend/internal/doc-prep"
	googlescraper "VLN-backend/internal/google-scraper"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		// log.Println("TASK_ID = ", task_id)

		// Start an interval to check the Requested ML sim image service every 2 seconds
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
			// else, the ticker loop continue
		}
	}
	return nil, nil
}

func getQualityScore(width float64, height float64) float64 {
	// function that caculates the quality scored based on the image height width
	// TODO needs implementation

	return (width * height) / 100.0
}

func getImageScores(sim_scores []similarityResponse, images []*googlescraper.Image) []imageScore {
	//format the images into a list of imageScore for sorting

	// imagess := []googlescraper.Image{
	// 	{ImagePath: "https://www.repstatic.it/content/contenthub/img/2022/03/03/093747590-e8cf7f34-63cf-4739-9fb6-75a6012a1421.png", Width: 200, Height: 200},
	// 	{ImagePath: "https://steemitimages.com/640x0/https://teorico.net/images/test-dgt-1.png", Width: 627, Height: 474},
	// 	{ImagePath: "https://global.discourse-cdn.com/business7/uploads/anki2/original/2X/6/6cce0bd15955f5b074262f6e76bf48c17c34c8b4.jpeg", Width: 605, Height: 405},
	// 	{ImagePath: "https://preview.redd.it/anyone-else-just-get-a-hello-test1-from-seank-notification-v0-e37e29jypjqa1.jpg?width=1080&crop=smart&auto=webp&s=180912e5884b274774cf8923a1ea6969b7b34a87", Width: 1080, Height: 795},
	// 	{ImagePath: "https://novadiagnostics.com.sg/wp-content/gallery/alifax-test1-2-0/Test1-2.0-Machine-Image-Front-View.jpg", Width: 1997, Height: 2500},
	// 	{ImagePath: "https://ultrafino.com/cdn/shop/files/product1_2000x.jpg?v=1711131062", Width: 2000, Height: 3000},
	// 	{ImagePath: "https://www.voromotors.com/cdn/shop/files/c8975a61-ae8b-40f3-9d6a-4fa02bf54663_9c6431e1-5f55-4392-87fa-b55b872d9a08.jpg?v=1712709946", Width: 1000, Height: 1000},
	// 	{ImagePath: "https://www.indumed.be/wp-content/uploads/2023/01/ResizedImage600551-new-test1-logo.png", Width: 600, Height: 551},
	// 	{ImagePath: "https://healthmanagement.org/uploads/product_image/mp_img_15022.jpg", Width: 592, Height: 552},
	// 	{ImagePath: "https://www.crazymuscle.com/cdn/shop/files/Test1-90Capsules_720x720.png?v=1682958076", Width: 720, Height: 720},
	// 	{ImagePath: "https://img.medicalexpo.com/images_me/photo-mg/67562-19159140.jpg", Width: 500, Height: 500},
	// 	{ImagePath: "https://image.slidesharecdn.com/test1-111101160429-phpapp02/85/Test1-1-320.jpg", Width: 320, Height: 240},
	// 	{ImagePath: "https://preview.redd.it/starbucks-intern-hard-at-work-v0-fohfusmcmjqa1.jpg?width=640&crop=smart&auto=webp&s=7b8b7058225694a89db6cc090f02c18e63bc5722", Width: 640, Height: 456},
	// 	{ImagePath: "https://m.media-amazon.com/images/I/41D2xPjVnFL._AC_UF1000,1000_QL80_.jpg", Width: 666, Height: 1000},
	// 	{ImagePath: "https://avatars.githubusercontent.com/u/22581?v=4", Width: 420, Height: 420},
	// 	{ImagePath: "https://i.ytimg.com/vi/DQpleciYOys/hq720.jpg?sqp=-oaymwEXCK4FEIIDSFryq4qpAwkIARUAAIhCGAE=&rs=AOn4CLCPSjnarzkVmmK_TJGg2rVcr5e2gg", Width: 686, Height: 386},
	// 	{ImagePath: "https://missoulaavalanche.org/wp-content/uploads/Snowmobile_Mountain_Riding_Lab_Photo_Shared_.jpeg", Width: 1368, Height: 1364},
	// 	{ImagePath: "https://global.discourse-cdn.com/business7/uploads/plot/original/3X/6/1/618cd92e235d8cb56d398064ae8376e320789e72.gif", Width: 690, Height: 470},
	// 	{ImagePath: "https://www.liftt.com/wp-content/uploads/2023/09/test1-1.png.webp", Width: 450, Height: 149},
	// 	{ImagePath: "https://opengraph.githubassets.com/6461d93afccdca6a7cf9d56c38c9c09d146d5bb2b413ba9e50e34e16db52b6c9/katjasrz/deepstream-test1-usb-people-count", Width: 1200, Height: 600},
	// 	{ImagePath: "https://i.ytimg.com/vi/Rsg7kpm-UK0/sddefault.jpg", Width: 640, Height: 480},
	// 	{ImagePath: "https://image.slidesharecdn.com/vocabularytest1-160221194007/85/Vocabulary-test1-3-320.jpg", Width: 320, Height: 453},
	// 	{ImagePath: "https://belreamed.com/assets/resized/538-538-crop-t/uploads/2020-03/2_0.jpg", Width: 538, Height: 538},
	// 	{ImagePath: "https://forum.obsidian.md/uploads/default/original/3X/6/5/6577c1f69507b84882ebe1afc9d16890479f6634.png", Width: 376, Height: 340},
	// 	{ImagePath: "https://image.spreadshirtmedia.net/image-server/v1/mp/products/T1459A839PA4459PT28D191467568W8121H10000/views/1,width=1200,height=630,appearanceId=839,backgroundColor=F2F2F2/haha-very-funny-test1-sticker.jpg", Width: 1200, Height: 630},
	// 	{ImagePath: "https://imgv2-1-f.scribdassets.com/img/document/652201752/original/98d571ab51/1715657008?v=1", Width: 768, Height: 1024},
	// 	{ImagePath: "https://opengraph.githubassets.com/cd864c10b37ff55e84fa8d05c6c1b96cf734e59899dc6212176aa1528e447fa7/openmrs/openmrs-test-test1", Width: 1200, Height: 600},
	// 	{ImagePath: "https://static.islcollective.com/storage/preview/201206/766x1084/final-test1-easy-grammar-drills-tests_26553_1.jpg", Width: 766, Height: 1084},
	// 	{ImagePath: "https://upload.wikimedia.org/wikipedia/en/thumb/8/80/GAA_logo-test1.png/640px-GAA_logo-test1.png", Width: 640, Height: 636},
	// 	{ImagePath: "https://d3tvd1u91rr79.cloudfront.net/96944be8ef135734141c1930c40e6cd5/html/bg1.png?Policy=eyJTdGF0ZW1lbnQiOlt7IlJlc291cmNlIjoiaHR0cHM6XC9cL2QzdHZkMXU5MXJyNzkuY2xvdWRmcm9udC5uZXRcLzk2OTQ0YmU4ZWYxMzU3MzQxNDFjMTkzMGM0MGU2Y2Q1XC9odG1sXC8qLnBuZyIsIkNvbmRpdGlvbiI6eyJEYXRlTGVzc1RoYW4iOnsiQVdTOkVwb2NoVGltZSI6MTcxNTkxMTQ3Mn19fV19&Signature=Jgssx5EkDpXT8FiagXSzD-WCey5yCheSZ0F7UqviHOF9LIJxqXsZrv4LEpKcQ2xnuvdVfuc9ePym4skPgKuWWbBlpOy0l1qW1hcGhXdrTQlleYaGVE1AAB9TWhwFMeYUJhDL6~hphNiGLIWEpqguGsvXCN6ZZrr7nv5LDFKD68A8g6vBMtJz27W6FjX6iPyA6IctbzBkgKNLaIbF5ZbF~TckEZxXW4sbET~-8X2AkUknE3xRLGhV-Qr0aCCjzO9DKks-Q0H2-KWztcBFxNWFPmI6VhSoRM3rOCtAtziP5whX-ZmQz1NQ3mBIdMzTvGh4m8ecyO7zDpbFbXNrruDzQ__&Key-Pair-Id=APKAIMKXCRHN6VGBBSZA", Width: 1200, Height: 628},
	// 	{ImagePath: "https://opengraph.githubassets.com/7dc8d7b74f70ad55b7746eea0980eb32e85b08ef97f7a345ad269f4fdc434b80/lancedikson/test1", Width: 1200, Height: 600},
	// }

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

	// similarity_scores := []similarityResponse{
	// 	{1.1283223628997803, "/app/go_backend/images/TEST1/0.png", "test it logo"},
	// 	{1.3097389936447144, "/app/go_backend/images/TEST1/1.png", "a close up of a green sign with a white number on it"},
	// 	{1.1162952184677124, "/app/go_backend/images/TEST1/2.jpeg", "a blue sign with the word test on it"},
	// 	{1.3658424615859985, "/app/go_backend/images/TEST1/3.jpg", "there is a photo of a starbucks sign on a cell phone"},
	// 	{1.3785594701766968, "/app/go_backend/images/TEST1/4.jpg", "there is a microwave that is on a table with a remote"},
	// 	{1.325546383857727, "/app/go_backend/images/TEST1/5.jpg", "arafed hat with black band and white hat with black band"},
	// 	{1.375824213027954, "/app/go_backend/images/TEST1/6.jpg", "a close up of a scooter with a handlebar and a seat"},
	// 	{1.3381946086883545, "/app/go_backend/images/TEST1/7.png", "a close up of a machine with a blue button on the side"},
	// 	{1.33420729637146, "/app/go_backend/images/TEST1/8.jpg", "a close up of a machine with a button on the front"},
	// 	{1.3754242658615112, "/app/go_backend/images/TEST1/9.png", "a close up of a nutrition label on a black background"},
	// }

	image_scores := getImageScores(similarity_scores, images)
	simWeight, qualityWeight := config.GetWeight()
	log.Println("simWeight = ", simWeight, "qualityWeight = ", qualityWeight)

	//todo the weighting function needs implementation
	sort.Slice(image_scores, func(i, j int) bool {
		s1, s2 := image_scores[i], image_scores[j]
		// weighted sum per image (slice)
		weight1 := simWeight*s1.SimilarityScore + qualityWeight*s1.QualityScore
		weight2 := simWeight*s2.SimilarityScore + qualityWeight*s2.QualityScore

		// HIGHER weighted sum is considered favorable
		return weight1 > weight2
	})
	log.Println("[Weighted] image_scores ", image_scores)

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

var (
	taskMap         = make(map[string][]imageScore)
	taskMapKeyWords = make(map[string]docprep.KeywordResponse)
	mutex           sync.RWMutex
)

func ProcessDoc(c *gin.Context) {
	/*
	   request format in params:
	   {
	       file_name: doc file
	   }

	   response: returns list of keywords and sentence
	*/
	taskID := uuid.New().String()

	go func(id string) {
		file_name, key := c.GetQuery("file_name")
		if !key {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error: ": "Invalid request parameter for processing doc",
			})
			c.Abort()
		}
		log.Println("[*] Processing file = ", file_name)
		log.Println("[*] Goroutine TaskID = ", id)

		keywords, err := docprep.ExtractKeywords(file_name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Error while extracting keywords",
				"error":  err,
			})
			c.Abort()

		}

		log.Println("keywords = ", keywords)
		mutex.Lock()
		defer mutex.Unlock()
		taskMapKeyWords[id] = keywords
	}(taskID)

	c.JSON(http.StatusOK, gin.H{
		"go_task_id": taskID,
		"endpoint":   "process_doc",
		"status":     "PROCESSING",
	})

}

func GetImages(c *gin.Context) {
	// Image Selection Process
	taskID := uuid.New().String()

	// goroutine <-- 1 query
	go func(id string) {
		query, key := c.GetQuery("query")

		log.Println("[*] Processing query = ", query)
		log.Println("[*] Goroutine TaskID = ", id)

		if !key {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error: ": "Invalid request parameter ",
				"State: ": "Getting images",
			})
			c.Abort()
		}

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

		// Request ML service -- VLNML
		// var images []*googlescraper.Image
		a, b := config.GetWeight()
		if a == 0 && b == 0 {
			fmt.Println("ERROR while getting Weighted values --  exiting the GO backend")
			os.Exit(1)
		}
		image_scores, err := filterImages(images, query)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"Status: ": "Error while filtering downloaded images",
				"Error":    err,
			})
			c.Abort()
		}

		// Update task status with the generated image list
		mutex.Lock()
		defer mutex.Unlock()
		taskMap[id] = image_scores
	}(taskID)

	c.JSON(http.StatusOK, gin.H{
		"go_task_id": taskID,
		"endpoint":   "get_images",
		"status":     "PROCESSING",
	})

}

func GetKeywordStatus(c *gin.Context) {
	// get taskID from query param
	taskID, key := c.GetQuery("task_id")
	// log.Println(" Current taskMap = ", taskMap)
	if !key {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error: ": "Invalid request parameter",
			"State: ": "Retrieve task status",
		})
		c.Abort()
	}

	mutex.RLock()
	defer mutex.RUnlock()
	log.Println("taskMapKeyWords = ", taskMapKeyWords[taskID])
	if keywords, ok := taskMapKeyWords[taskID]; ok {
		// log.Printf("COMPLETED task_id = %s", taskID)
		c.JSON(http.StatusOK, gin.H{
			"task_id":  taskID,
			"status":   "SUCCESS",
			"keywords": keywords,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"task_id": taskID,
			"status":  "PENDING",
		})
	}
}

func GetImageStatus(c *gin.Context) {
	// get taskID from query param
	taskID, key := c.GetQuery("task_id")
	// log.Println(" Current taskMap = ", taskMap)
	if !key {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error: ": "Invalid request parameter",
			"State: ": "Retrieve task status",
		})
		c.Abort()
	}

	mutex.RLock()
	defer mutex.RUnlock()

	if image_scores, ok := taskMap[taskID]; ok {
		// log.Printf("COMPLETED task_id = %s", taskID)
		c.JSON(http.StatusOK, gin.H{
			"task_id":      taskID,
			"status":       "SUCCESS",
			"image_scores": image_scores,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"task_id": taskID,
			"status":  "PENDING",
		})
	}
}
