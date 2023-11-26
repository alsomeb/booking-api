package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

// A global variable that will hold a reference to the MongoDB client
var mongoClient *mongo.Client

// Booking - tags help the MongoDB driver understand how to serialize and deserialize data between your Go code and the MongoDB database.
// Not needed but recommended
/*
	ID corresponds to the _id field in MongoDB documents.
	Name corresponds to the name field in MongoDB documents.
	DateAdded corresponds to the dateAdded field in MongoDB documents.
*/
type Booking struct {
	Id    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Added string             `bson:"added"`
}

// Loads current env file and returns requested value by key
func loadEnvVariable(key string) string {
	err := godotenv.Load(".env") // Load will read your env file(s) and load them into ENV for this process.

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

// The init function will run before our main function to establish a connection to MongoDB. If it cannot connect it will fail and the program will exit.
func init() {
	var username = loadEnvVariable("DB_USER")
	var password = loadEnvVariable("DB_PASSWORD")
	var mongoURI = fmt.Sprintf("mongodb+srv://%v:%v@alsomeb.jcl49rx.mongodb.net/Bookings?retryWrites=true&w=majority", username, password)

	if err := connectToMongoDB(mongoURI); err != nil {
		log.Fatal("Could not connect to MongoDB")
	}
}

// Our implementation logic for connecting to MongoDB
func connectToMongoDB(mongoURI string) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	mongoClient = client
	return err
}

// add Booking
func addBooking(c *gin.Context) {
	// Bind the request JSON to a Booking Struct

	/*
	   When you pass a variable to a function, it is passed by value, meaning the function receives a copy of the variable.
	   If you want the function to modify the original variable (rather than just a copy), you need to pass a pointer to that variable.
	   This is because pointers allow functions to indirectly modify the value they point to.

	   In this case, c.BindJSON(&newBooking) modifies the original 'newBooking' struct with the properties received from the payload.
	   Without passing a pointer (&newBooking), BindJSON would receive a copy of 'newBooking', and any modifications it makes
	   would not affect the original variable.

	   Since we receive properties from the payload, 'newBooking' needs to be modified in place.
	*/
	var newBooking Booking
	err := c.BindJSON(&newBooking)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	newBooking.Added = time.Now().UTC().Format(time.RFC3339) // "2006-01-02T15:04:05Z07:00"

	result, err := mongoClient.Database("Bookings").Collection("Bookings").InsertOne(context.TODO(), newBooking)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return the id of the newly inserted document
	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

// Collects all bookings
func getAllBookings(c *gin.Context) {

	// Find movies
	// bson.D is a type representing a BSON document in Go.
	// {{}} is a composite literal creating an instance of bson.D with an empty BSON document.
	// bson.D{{}} is an empty filter that retrieves all documents when used in a MongoDB query.
	cursor, err := mongoClient.Database("Bookings").Collection("Bookings").Find(context.TODO(), bson.D{{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map results to my struct
	results := make([]Booking, 0)
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func getBookingById(c *gin.Context) {
	// Get movie ID from URL
	idStr := c.Param("id")

	id, errConvertId := primitive.ObjectIDFromHex(idStr) // Mongo ObjectID convert from String
	if errConvertId != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errConvertId.Error()})
		return
	}

	var bookingResult Booking
	errFetch := mongoClient.Database("Bookings").Collection("Bookings").FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&bookingResult)
	if errFetch != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": errFetch.Error()})
		return
	}

	c.JSON(http.StatusOK, bookingResult)
}

func main() {
	r := gin.Default()

	// Set base URL for the API
	api := r.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{ // H is a shorthand for MAP - new map instance with the given values
			"message": "Hello World",
		})
	})

	api.GET("/bookings", getAllBookings)
	api.GET("/bookings/:id", getBookingById)
	api.POST("/bookings/add", addBooking)

	log.Println("Mongo DB Connected")
	_ = r.Run()
}
