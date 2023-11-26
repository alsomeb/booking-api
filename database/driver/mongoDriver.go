package driver

import (
	"booking-api/structs"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func ConnectToMongo() *structs.MongoClient {
	var username = loadEnvVariable("DB_USER")
	var password = loadEnvVariable("DB_PASSWORD")
	var mongoURI = fmt.Sprintf("mongodb+srv://%v:%v@alsomeb.jcl49rx.mongodb.net/Bookings?retryWrites=true&w=majority", username, password)
	log.Println("Loaded .env")

	client := setUpMongo(mongoURI)

	return client
}

// Loads current env file and returns requested value by key
func loadEnvVariable(key string) string {
	err := godotenv.Load(".env") // Load will read your env file(s) and load them into ENV for this process.

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

func setUpMongo(mongoURI string) *structs.MongoClient {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatalf("failed to ping MongoDB: %v", err)
	}

	mongoClient := structs.MongoClient{MongoClient: client}

	return &mongoClient
}
