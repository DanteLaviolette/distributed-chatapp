package persistence

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client = nil
var ctx context.Context
var ctxCancel context.CancelFunc

/*
Initialize the database connection.
*/
func InitializeDBConnection() {
	// Only initialize client if not set (singleton pattern)
	if client == nil {
		// Create context
		ctx, ctxCancel = context.WithTimeout(context.Background(), 10*time.Second)
		// Connect to mongoDB
		var err error
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
		// Fail if connection error occurred
		if err != nil {
			log.Fatal(err)
		}
		// Verify connection was made via pinging db
		err = client.Ping(ctx, nil)
		// Fail if connection error occurred
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Connected to MongoDB")
		// Handle cleanup on exit
		cleanupDBConnectionOnExit()
		// Setup indices
		setupIndices()
	}
}

/*
Runs a background function to cleanup the database connection on
exit.
*/
func cleanupDBConnectionOnExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// Run exit func in background waiting for exit signal
	go func() {
		<-sigs
		// Disconnect client
		client.Disconnect(ctx)
		// Cleanup context
		ctxCancel()
		os.Exit(0)
	}()
}

/*
Sets up the indices on the table if they don't yet exist
- db.users.createIndex( { "email": 1 }, { unique: true } )
*/
func setupIndices() {
	// Define index
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"email": 1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}
	// Add index to user collection
	col, _ := getUserCollection()
	_, err := col.Indexes().CreateOne(ctx, indexModel)
	// Fail if error occurred
	if err != nil {
		log.Fatal(err)
	}
}

/*
Returns the user collection along with its context
*/
func getUserCollection() (mongo.Collection, context.Context) {
	return *(client.Database(os.Getenv("DB_NAME")).
		Collection(os.Getenv("USER_COLLECTION_NAME"))), ctx
}
