package config

import (
	"log"
	"os"
	"strconv"

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

func GetWeight() (float64, float64) {
	simWeight, _ := strconv.ParseFloat(os.Getenv("SIM_WEIGHT"), 64)
	qualityWeight, _ := strconv.ParseFloat(os.Getenv("QUALITY_WEIGHT"), 64)
	return simWeight, qualityWeight
}
