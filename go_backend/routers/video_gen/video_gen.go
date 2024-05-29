package video_gen

import (
	ffmpeg "VLN-backend/internal/ffmpeg"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	taskMapVideo = make(map[string]string)
	mutex        sync.RWMutex
)

func CreateVideo(c *gin.Context) {
	/*
		   request format in Params:
		   {
		       audio_path: ["au1.mp3","au2.mp3"]
			   image_path: []
		   }

		   response: returns list of keywords and sentence
	*/

	taskID := uuid.New().String()

	go func(id string) {
		audio, audio_key := c.GetQuery("audio_path")
		images, images_key := c.GetQuery("image_path")

		if !audio_key || !images_key {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error: ": "Invalid request parameter for TTS",
			})
			c.Abort()
		}

		log.Println("[*] Processing video")
		log.Println("[*] Goroutine TaskID = ")

		log.Println("[*] audio = ", audio)
		log.Println("[*] images = ", images)

		res_video, err := ffmpeg.Test(audio, images)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Error while Creating Video",
				"error":  err,
			})
			c.Abort()
		}

		log.Println("res_video = ", res_video)
		mutex.Lock()
		defer mutex.Unlock()
		taskMapVideo[id] = res_video

	}(taskID)

	c.JSON(http.StatusOK, gin.H{
		"go_task_id": taskID,
		"endpoint":   "get_ffmpeg",
		"status":     "PROCESSING",
	})

}

func GetVideotatus(c *gin.Context) {
	taskID, key := c.GetQuery("task_id")
	if !key {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error: ": "Invalid request parameter -- Expecting ?=task_id (go_task_id)",
			"State: ": "Retrieve task status",
		})
		c.Abort()
	}

	mutex.RLock()
	defer mutex.RUnlock()

	if video, ok := taskMapVideo[taskID]; ok {
		c.JSON(http.StatusOK, gin.H{
			"go_task_id": taskID,
			"status":     "SUCCESS",
			"video":      video,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"go_task_id": taskID,
			"status":     "PENDING",
		})
	}
}
