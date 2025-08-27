package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func InitMongo() {
	_ = godotenv.Load()
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("mongo connect:", err)
	}

	// ✅ Verify connection with Ping
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("mongo ping:", err)
	}

	Client = client
	DB = client.Database(dbName)
	log.Println("✅ Connected to MongoDB:", dbName)
}
