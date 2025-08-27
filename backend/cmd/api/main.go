package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/Rjrajnish99/photo-gallery/backend/internal/db"
)

func main() {
	_ = godotenv.Load()
	db.InitMongo()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := SetupRouter()
	log.Println("API on :" + port)
	_ = r.Run(":" + port)
}
