package test

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	taskMap = make(map[string][]int)
	mutex   sync.RWMutex
)

func generateimage_scores() []int {
	time.Sleep(30 * time.Second)

	return []int{1, 2, 3}
}

func StartTask(c *gin.Context) {
	taskID := uuid.New().String()

	go func() {
		image_scores := generateimage_scores()

		// Update task status with the generated image list
		mutex.Lock()
		defer mutex.Unlock()
		taskMap[taskID] = image_scores
	}()

	c.JSON(http.StatusOK, gin.H{
		"task_id": taskID,
		"message": "Task processing started",
	})
}

func GetTaskStatus(c *gin.Context) {
	// get taskID from query param
	taskID := c.Query("task_id")
	log.Println("Requested task_id = ", taskID)
	log.Println("taskMap = ", taskMap)

	mutex.RLock()
	defer mutex.RUnlock()

	// Check if the task ID exists in the taskMap
	if image_scores, ok := taskMap[taskID]; ok {

		// Task ID found, return task status with the generated image list
		log.Printf("COMPLETE task_id = %s", taskID)
		c.JSON(http.StatusOK, gin.H{
			"task_id":      taskID,
			"status":       "COMPLETED",
			"image_scores": image_scores,
		})
	} else {
		// Task ID not found, return status as "PENDING"
		c.JSON(http.StatusOK, gin.H{
			"task_id": taskID,
			"status":  "PENDING",
		})
	}
}

// EXPERIMENTING
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
