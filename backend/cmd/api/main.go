package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	// ✅ apply CORS before adding routes
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// add routes
	SetupRouter(r)

	log.Println("API on :" + port)
	_ = r.Run(":" + port)
}
