package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
		log.Fatal("Error loading .env file")
	}
}

func GetMLService() string {
	return os.Getenv("VLNML_SERVER_DEV")
}

func GetSerapAPIKey() string {
	return os.Getenv("SERPAPI_KEY")
}
