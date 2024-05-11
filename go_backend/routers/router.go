package routers

import (
	"VLN-backend/routers/health"
	"VLN-backend/routers/inference"
	"VLN-backend/routers/scraper"
	"VLN-backend/routers/test"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Gin Services
	r.POST("/health", health.GetHealth)
	r.POST("/tts", inference.TextToSpeech)
	r.POST("/get_images", scraper.GetImages)
	r.GET("/status", scraper.GetTaskStatus)

	r.POST("/test_start_task", test.StartTask)
	r.GET("/test_status", test.GetTaskStatus)

	return r
}
