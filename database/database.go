package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// MongoInstance contains the Mongo client and database objects
type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var Mg MongoInstance

// Database settings (later we will transfer these info to config)
const dbName = "blog"
const mongoURI = "mongodb://localhost:27017/" + dbName

// Connect function configures the MongoDB client and establishes db connection.
func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		return err
	}

	Mg = MongoInstance{
		Client: client,
		Db:     db,
	}

	return nil
}
