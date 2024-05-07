package main

import (
	"VLN-backend/config"
	"VLN-backend/routers"
)

func main() {
	config.Load()
	r := routers.InitRouter()
	r.Run(":9888")
}
