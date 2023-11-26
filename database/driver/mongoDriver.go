package driver

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func ConnectToMongo() *mongo.Client {
	var username = loadEnvVariable("DB_USER")
	var password = loadEnvVariable("DB_PASSWORD")
	var mongoURI = fmt.Sprintf("mongodb+srv://%v:%v@alsomeb.jcl49rx.mongodb.net/Bookings?retryWrites=true&w=majority", username, password)

	err, client := setUpMongo(mongoURI)

	if err != nil {
		panic(err)
	}

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

func setUpMongo(mongoURI string) (error, *mongo.Client) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)

	return err, client
}
