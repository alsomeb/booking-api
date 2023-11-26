package repo

import (
	"booking-api/structs"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// AddBooking - add new Booking
func AddBooking(newBooking *structs.Booking, mongoClient *mongo.Client) (*mongo.InsertOneResult, error) {
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

	newBooking.CreatedAt = time.Now().UTC()
	newBooking.UpdatedAt = time.Now().UTC()

	result, err := mongoClient.Database("Bookings").Collection("Bookings").InsertOne(context.TODO(), newBooking)

	return result, err
}

// GetAllBookings - Collects all bookings
func GetAllBookings(mongoClient *mongo.Client) (*mongo.Cursor, error) {
	// Find movies
	// bson.D is a type representing a BSON document in Go.
	// {{}} is a composite literal creating an instance of bson.D with an empty BSON document.
	// bson.D{{}} is an empty filter that retrieves all documents when used in a MongoDB query.
	cursor, err := mongoClient.Database("Bookings").Collection("Bookings").Find(context.TODO(), bson.D{{}})

	return cursor, err
}

// GetBookingById - Fetch booking by id (ObjectId)
func GetBookingById(id primitive.ObjectID, mongoClient *mongo.Client) (structs.Booking, error) {
	var bookingResult structs.Booking
	err := mongoClient.Database("Bookings").Collection("Bookings").FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&bookingResult)

	return bookingResult, err
}
