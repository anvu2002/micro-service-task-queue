package routers

import (
	"VLN-backend/routers/inference"
	"VLN-backend/routers/scraper"

	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.POST("/tts", inference.TextToSpeech)
	r.POST("/get_images", scraper.GetImages)
	// r.GET("/status")
	return r
}
