package main

import (
	 
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rjrajnish/photo_gallery/backend/internal/db"
)

func main() {
	_ = godotenv.Load(".env")
	 
	db.InitMongo()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := SetupRouter()

	log.Println("API on :" + port)
	_ = r.Run(":" + port)
}
