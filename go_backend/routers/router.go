package routers

import (
	"VLN-backend/routers/health"
	"VLN-backend/routers/inference"
	"VLN-backend/routers/video_gen"

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
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 24 * time.Hour,
	}))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Gin API Endpoints
	r.GET("/health", health.GetHealth)

	r.POST("/process_doc", scraper.ProcessDoc)
	r.GET("/keyword_status", scraper.GetKeywordStatus)

	r.POST("/get_images", scraper.GetImages)
	r.GET("/image_status", scraper.GetImageStatus)

	r.POST("/tts", inference.TextToSpeech)
	r.GET("/tts_status", inference.GetTtsStatus)

	r.POST("/get_ffmpeg", video_gen.CreateVideo)
	r.GET("/video_status", video_gen.GetVideotatus)

	// Testing endpoints
	r.POST("/test_start_task", test.StartTask)
	r.GET("/test_status", test.GetTaskStatus)

	return r
}
