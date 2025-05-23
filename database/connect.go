package database

import (
	"CrispyBot/variables"
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	// Singleton instance of the database connection
	instance *DB
	once     sync.Once
)

// DB represents a MongoDB client and database
type DB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// GetDB returns the singleton database instance
func DBInit() *DB {
	once.Do(func() {
		instance = connectToDB()
	})
	return instance
}

// connectToDB establishes connection to MongoDB and returns a DB
func connectToDB() *DB {
	fmt.Println("Connecting to MongoDB...")

	// Set up MongoDB connection options
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(variables.Mongodb_uri).
		SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
	}

	// Verify the connection
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(fmt.Sprintf("Failed to ping MongoDB: %v", err))
	}

	fmt.Println("Successfully connected to MongoDB database:", variables.Db_name)

	// Get the database
	database := client.Database(variables.Db_name)

	return &DB{
		Client:   client,
		Database: database,
	}
}

// Close closes the database connection
func (db *DB) Close() {
	if db != nil && db.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := db.Client.Disconnect(ctx); err != nil {
			fmt.Printf("Error disconnecting from MongoDB: %v\n", err)
		} else {
			fmt.Println("Successfully disconnected from MongoDB")
		}
	}
}

// GetCollection returns a handle to the specified collection
func (db *DB) GetCollection(name string) *mongo.Collection {
	if db == nil || db.Database == nil {
		panic("Database connection not initialized")
	}
	return db.Database.Collection(name)
}
