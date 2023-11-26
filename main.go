package main

import (
	"booking-api/auth"
	"booking-api/database/driver"
	"booking-api/database/repo"
	"booking-api/structs"
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

// add Booking
func addBooking(c *gin.Context, client *structs.MongoClient) {
	var newBooking structs.Booking

	err := c.BindJSON(&newBooking)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	result, err := repo.AddBooking(&newBooking, client)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)

}

// Collects all bookings
func getAllBookings(c *gin.Context, client *structs.MongoClient) {
	cursor, err := repo.CollectAllBookings(client)
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

func getBookingById(c *gin.Context, client *structs.MongoClient) {
	// Get movie ID from URL-Params
	idStr := c.Param("id")

	id, errConvertId := primitive.ObjectIDFromHex(idStr) // Mongo ObjectID convert from String
	if errConvertId != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errConvertId.Error()})
		return
	}

	result, err := repo.GetBookingById(id, client)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func verifyToken(c *gin.Context, firebaseClient *firebase.App) {
	// Get the Firebase Auth client
	client, err := firebaseClient.Auth(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Firebase Auth client"})
		return
	}

	// Get the token from the request header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
		return
	}

	// verify token the token
	decodedToken, err := client.VerifyIDToken(c, token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims := decodedToken.Claims

	c.JSON(http.StatusOK, claims)
}

func main() {
	appRouter := gin.Default()

	// Set base URL for the API
	api := appRouter.Group("/api")

	// init mongo client
	mongoClient := driver.ConnectToMongo()

	// init firebase app
	firebaseClient := auth.InitFireBaseApp()

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{ // H is a shorthand for MAP - new map instance with the given values
			"message": "Hello World",
		})
	})

	api.GET("/bookings", func(c *gin.Context) {
		getAllBookings(c, mongoClient)
	})

	api.GET("/verify-token", func(c *gin.Context) {
		verifyToken(c, firebaseClient)
	})

	api.GET("/bookings/:id", func(c *gin.Context) {
		getBookingById(c, mongoClient)
	})

	api.POST("/bookings/add", func(c *gin.Context) {
		addBooking(c, mongoClient)
	})

	log.Println("Mongo DB Connected")
	err := appRouter.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}
