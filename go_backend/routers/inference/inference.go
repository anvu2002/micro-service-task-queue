package inference

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"VLN-backend/config"

	"github.com/gin-gonic/gin"
)

func TextToSpeech(c *gin.Context) {
	tts_api := config.GetMLService() + "/api/tts"
	text, text_key := c.GetQuery("text")
	path, path_key := c.GetQuery("save_path")

	tts_request_string := fmt.Sprintf("{\"text\": \"%s\",\n\"save_path\": \"%s\"}", text, path)

	if !text_key || !path_key {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error: ": "Invalid request parameter for TTS",
		})
		c.Abort()
	}

	body, _ := json.Marshal(tts_request_string)

	resp, err := http.Post(tts_api, "application/json", bytes.NewBuffer(body))

	if err != nil {
		log.Print("TTS Request Failed\n")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		log.Print("TTS Request Success\n")
		resp_body, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusOK, gin.H{
			"message": json.RawMessage(resp_body),
		})
	}
}
