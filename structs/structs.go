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

type MongoClient struct {
	MongoClient *mongo.Client
}
