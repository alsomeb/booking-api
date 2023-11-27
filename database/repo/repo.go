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
func AddBooking(newBooking *structs.Booking, client *structs.MongoClient) (*mongo.InsertOneResult, error) {
	/*
		   When you pass a variable to a function, it is passed by value, meaning the function receives a copy of the variable.
		   If you want the function to modify the original variable (rather than just a copy), you need to pass a pointer to that variable.
		   This is because pointers allow functions to indirectly modify the value they point to.

			The newBooking parameter is a pointer to the original structs.Booking instance. If you modify the fields of newBooking inside the function,
			you will be modifying the fields of the original structs.Booking variable outside the function.

			Similarly, the client parameter is a pointer to the original structs.MongoClient instance.
			If you modify the fields of client inside the function, you will be modifying the fields of the original structs.MongoClient variable outside the function.
	*/

	newBooking.CreatedAt = time.Now().UTC()
	newBooking.UpdatedAt = time.Now().UTC()

	result, err := client.MongoClient.Database("Bookings").Collection("Bookings").InsertOne(context.TODO(), newBooking)

	return result, err
}

// CollectAllBookings - Collects all bookings
func CollectAllBookings(client *structs.MongoClient) (*mongo.Cursor, error) {
	// Find movies
	// bson.D is a type representing a BSON document in Go.
	// {{}} is a composite literal creating an instance of bson.D with an empty BSON document.
	// bson.D{{}} is an empty filter that retrieves all documents when used in a MongoDB query.
	cursor, err := client.MongoClient.Database("Bookings").Collection("Bookings").Find(context.TODO(), bson.D{{}})

	return cursor, err
}

// GetBookingById - Fetch booking by id (ObjectId)
func GetBookingById(id primitive.ObjectID, client *structs.MongoClient) (structs.Booking, error) {
	var bookingResult structs.Booking
	err := client.MongoClient.Database("Bookings").Collection("Bookings").FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&bookingResult)

	return bookingResult, err
}
