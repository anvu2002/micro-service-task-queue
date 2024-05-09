package test

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	imageList []int

	mutex sync.RWMutex
)

func generateImageList() {
	time.Sleep(10 * time.Second)

	mutex.Lock()
	defer mutex.Unlock()
	imageList = []int{1, 2, 3}
}

func StartTask(c *gin.Context) {
	go func() {
		generateImageList()

	}()
	c.JSON(http.StatusOK, gin.H{
		"message": "Task processing started",
	})
}

func GetTaskStatus(c *gin.Context) {
	mutex.RLock()
	defer mutex.RUnlock()
	status := imageList

	c.JSON(http.StatusOK, gin.H{
		"imageList": status,
	})
}

type imageScore struct {
	SimilarityScore float64 `json:"similarity_score"`
	QualityScore    float64 `json:"quality_score"`
	ImagePath       string  `json:"image_path"`
}

func TestFeature(c *gin.Context) {
	query, key1 := c.GetQuery("query")
	task_id, key2 := c.GetQuery("task_id")
	if !key1 || !key2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid request param for playing ping pong with GO",
		})
		c.Abort()
	}

	image_scores := make([]imageScore, 0)

	if task_id == "" {
		log.Println("Processing Similarity")
		log.Println("Requested data: ", query)
		for i := 0; i < 10; i++ {
			image_scores = append(image_scores, imageScore{
				SimilarityScore: 10,
				QualityScore:    12,
				ImagePath:       "/" + strconv.Itoa(i) + ".jpg",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"image_scores": image_scores,
			"status":       http.StatusOK,
			"task_id":      query,
		})
	} else {
		log.Println("Need to implement query task status through task_id")
	}

}
