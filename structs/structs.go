package structs

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// Booking - tags help the MongoDB driver understand how to serialize and deserialize data between your Go code and the MongoDB database.
// Not needed but recommended
/*
	ID corresponds to the _id field in MongoDB documents.
	Name corresponds to the name field in MongoDB documents.
	DateAdded corresponds to the dateAdded field in MongoDB documents.
*/
type Booking struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"` // Json kommer serialize som "id" ist för "_id"
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// MongoClient - holds a pointer to the MongoDB Client
// When you create an instance of this struct and pass it around, you are essentially passing a reference to the underlying mongo.Client instance
type MongoClient struct {
	MongoClient *mongo.Client
}

/*
	In Go, when you pass a struct to a function, a copy of the struct is made unless you explicitly pass it as a POINTER.
	When you pass a struct that contains a pointer to another object (in this case, a MongoDB client),
	you are still passing a copy of the struct, but since the pointer inside the struct points to the same underlying object, modifications to the object are visible to all references.

	When you pass *structs.MongoClient to functions like addBooking, getAllBookings, and getBookingById, you are passing the actual client and not the address of the client (&variable).
	The * in this context is used to dereference the pointer inside the MongoClient struct, giving you the actual mongo.Client object.

*/
