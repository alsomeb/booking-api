package main

import (
	"booking-api/database/driver"
	"booking-api/database/repo"
	"booking-api/structs"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

// A global variable that will hold a reference to the MongoDB client
var mongoClient *mongo.Client

// The init function will run before our main function to establish a connection to MongoDB. If it cannot connect it will fail and the program will exit.
func init() {
	mongoClient = driver.ConnectToMongo()
}

// add Booking
func addBooking(c *gin.Context) {
	var newBooking structs.Booking

	err := c.BindJSON(&newBooking)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	result, err := repo.AddBooking(&newBooking, mongoClient)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)

}

// Collects all bookings
func getAllBookings(c *gin.Context) {
	cursor, err := repo.GetAllBookings(mongoClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map results to my struct
	results := make([]structs.Booking, 0)
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func getBookingById(c *gin.Context) {
	// Get movie ID from URL-Params
	idStr := c.Param("id")

	id, errConvertId := primitive.ObjectIDFromHex(idStr) // Mongo ObjectID convert from String
	if errConvertId != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errConvertId.Error()})
		return
	}

	result, err := repo.GetBookingById(id, mongoClient)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func main() {
	appRouter := gin.Default()

	// Set base URL for the API
	api := appRouter.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{ // H is a shorthand for MAP - new map instance with the given values
			"message": "Hello World",
		})
	})

	api.GET("/bookings", getAllBookings)
	api.GET("/bookings/:id", getBookingById)
	api.POST("/bookings/add", addBooking)

	log.Println("Mongo DB Connected")
	err := appRouter.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
