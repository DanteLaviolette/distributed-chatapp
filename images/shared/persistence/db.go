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
var databaseName string

/*
Initialize the database connection.
dbUrl: Database url to connect to
dbName: Database name to connect to
*/
func InitializeDBConnection(dbUrl string, dbName string) {
	// Only initialize client if not set (singleton pattern)
	if client == nil {
		databaseName = dbName
		// Create context
		ctx, ctxCancel = context.WithTimeout(context.Background(), 10*time.Second)
		// Connect to mongoDB
		var err error
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(dbUrl))
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
Returns the user collection along with its context
*/
func GetUserCollection() (mongo.Collection, context.Context) {
	return *(client.Database(databaseName).
		Collection("Users")), ctx
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
	col, _ := GetUserCollection()
	_, err := col.Indexes().CreateOne(ctx, indexModel)
	// Fail if error occurred
	if err != nil {
		log.Fatal(err)
	}
}
