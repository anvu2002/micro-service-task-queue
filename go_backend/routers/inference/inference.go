package inference

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"VLN-backend/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ttsRequest struct {
	Text     string `json:"text"`
	SavePath string `json:"save_path"`
}

type ttsResponse struct {
	AudioPath string `json:"save_path"`
	Text      string `json:"text"`
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

var (
	taskMapTTS = make(map[string]ttsResponse)
	mutex      sync.RWMutex
)

func tts(text string, path string) (ttsResponse, error) {

	ttsAPI := config.GetMLService() + "/api/tts"

	request_body, _ := json.Marshal((ttsRequest{
		Text:     text,
		SavePath: path,
	}))

	resp, err := http.Post(ttsAPI, "application/json", bytes.NewBuffer(request_body))

	if err != nil {
		log.Print("TTS Request Failed\n")
		return ttsResponse{}, err
	}
	defer resp.Body.Close()

	log.Print("TTS Request Success\n")
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ttsResponse{}, err
	}

	var tts_res []VLNMLResponse

	tts_resp_format_err := json.Unmarshal([]byte(string(respBody)), &tts_res)

	if tts_resp_format_err != nil {
		log.Println("ERROR when trying to map tts response", err)
		return ttsResponse{}, err
	}

	for _, res := range tts_res {
		task_id := res.TaskID
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			status_api := config.GetMLService() + "/api/status/" + task_id
			statusResp, err := http.Get(status_api)
			if err != nil {
				return ttsResponse{}, err
			}
			defer statusResp.Body.Close()

			statusBody, err := io.ReadAll(statusResp.Body)
			if err != nil {
				return ttsResponse{}, err
			}

			var statusResponse StatusResponse
			if err := json.Unmarshal(statusBody, &statusResponse); err != nil {
				return ttsResponse{}, err
			}

			if statusResponse.Status == "SUCCESS" {
				result_api := config.GetMLService() + "/api/result/" + task_id
				resultResp, err := http.Get(result_api)
				if err != nil {
					return ttsResponse{}, err
				}
				defer resultResp.Body.Close()

				resultBody, err := io.ReadAll(resultResp.Body)

				var tts_res ttsResponse

				json_format_err := json.Unmarshal([]byte(string(resultBody)), &tts_res)
				if json_format_err != nil {
					log.Println("ERROR while maping the result value", json_format_err)
					return ttsResponse{}, err
				}
				return tts_res, nil
			}
			if statusResponse.Status == "FAIL" {
				return ttsResponse{}, errors.New("task failed")
			}

		}
	}
	return ttsResponse{}, nil
}

func TextToSpeech(c *gin.Context) {
	/*
		   request format in Params:
		   {
		       text: "bla bla",
			   save_path:
		   }

		   response: returns list of keywords and sentence
	*/

	taskID := uuid.New().String()

	go func(id string) {
		text, text_key := c.GetQuery("text")
		path, path_key := c.GetQuery("save_path")

		if !text_key || !path_key {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error: ": "Invalid request parameter for TTS",
			})
			c.Abort()
		}

		log.Println("[*] Processing tts = ", text)
		log.Println("[*] Goroutine TaskID = ", id)

		speech, err := tts(text, path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Status": "Error while converting text to speech",
				"error":  err,
			})
			c.Abort()
		}

		log.Println("speech = ", speech)
		mutex.Lock()
		defer mutex.Unlock()
		taskMapTTS[id] = speech

	}(taskID)

	c.JSON(http.StatusOK, gin.H{
		"go_task_id": taskID,
		"endpoint":   "tts",
		"status":     "PROCESSING",
	})

}

func GetTtsStatus(c *gin.Context) {
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
	log.Println("taskMapTTS = ", taskMapTTS[taskID])
	if speech, ok := taskMapTTS[taskID]; ok {
		// log.Printf("COMPLETED task_id = %s", taskID)
		c.JSON(http.StatusOK, gin.H{
			"go_task_id": taskID,
			"status":     "SUCCESS",
			"speech":     speech,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"go_task_id": taskID,
			"status":     "PENDING",
		})
	}
}
