﻿package test

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

func generateImageList() []int {
	// Simulate long processing time to generate the image list
	time.Sleep(10 * time.Second)

	// Generate the image list
	return []int{1, 2, 3} // Replace this with your actual implementation
}

func StartTask(c *gin.Context) {
	// Generate a unique task ID
	taskID := uuid.New().String()

	// Start processing the task (simulate long processing time)
	go func() {
		imageList := generateImageList()

		// Update task status with the generated image list
		mutex.Lock()
		defer mutex.Unlock()
		taskMap[taskID] = imageList
	}()

	// Return a success response along with the task ID
	c.JSON(http.StatusOK, gin.H{
		"task_id": taskID,
		"message": "Task processing started",
	})
}

func GetTaskStatus(c *gin.Context) {
	// Retrieve the task ID from the request query parameters
	taskID := c.Query("task_id")

	mutex.RLock()
	defer mutex.RUnlock()

	// Check if the task ID exists in the taskMap
	if imageList, ok := taskMap[taskID]; ok {
		// Task ID found, return task status with the generated image list
		c.JSON(http.StatusOK, gin.H{
			"task_id":   taskID,
			"status":    "COMPLETED",
			"imageList": imageList,
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
