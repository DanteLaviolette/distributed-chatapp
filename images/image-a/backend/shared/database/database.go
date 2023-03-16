package database

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.violettedev.com/eecs4222/constants"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client = nil

const databaseInitializationTimeout = 10 * time.Second
const databaseName = "main"
const userCollectionName = "Users"
const refreshTokenCollectionName = "RefreshTokens"

/*
Initialize the database connection.
dbUrl: Database url to connect to
*/
func InitializeDBConnection(dbUrl string) {
	// Only initialize client if not set (singleton pattern)
	if client == nil {
		// Create context
		ctx, cancel := context.WithTimeout(context.Background(),
			databaseInitializationTimeout)
		defer cancel()
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
Returns the user collection
*/
func GetUserCollection() mongo.Collection {
	return *(client.Database(databaseName).
		Collection(userCollectionName))
}

/*
Returns the refresh token collection
*/
func GetRefreshTokenCollection() mongo.Collection {
	return *(client.Database(databaseName).
		Collection(refreshTokenCollectionName))
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
		// Create context
		ctx, cancel := context.WithTimeout(context.Background(),
			databaseInitializationTimeout)
		// Disconnect client
		client.Disconnect(ctx)
		// Cleanup context
		cancel()
		os.Exit(0)
	}()
}

/*
Sets up the indices on the table if they don't yet exist
- db.users.createIndex( { "email": 1 }, { unique: true } )
- db.RefreshTokens.createIndex( expireAfterSeconds: constants.RefreshTokenExpirySeconds } )
*/
func setupIndices() {
	createUserIndices()
	createRefreshTokenIndices()
}

/*
Creates user index: db.users.createIndex( { "email": 1 }, { unique: true } )
*/
func createUserIndices() {
	// Define index
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"email": 1, // index in ascending order
		},
		Options: options.Index().SetUnique(true),
	}
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(),
		databaseInitializationTimeout)
	defer cancel()
	// Add index to user collection
	col := GetUserCollection()
	_, err := col.Indexes().CreateOne(ctx, indexModel)
	// Fail if error occurred
	if err != nil {
		log.Fatal(err)
	}
}

/*
Creates refresh token index:
- db.RefreshTokens.createIndex( expireAfterSeconds: constants.RefreshTokenExpirySeconds } )
- db.RefreshTokens.createIndex( { "secret": 1 }, { unique: true } )
*/
func createRefreshTokenIndices() {
	// Define indices
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.M{
				"createdAt": 1,
			},
			Options: options.Index().SetExpireAfterSeconds(constants.RefreshTokenExpirySeconds),
		},
		{
			Keys: bson.M{
				"secret": 1, // index in ascending order
			},
			Options: options.Index().SetUnique(true),
		},
	}
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(),
		databaseInitializationTimeout)
	defer cancel()
	// Add index to refresh token collection
	col := GetRefreshTokenCollection()
	_, err := col.Indexes().CreateMany(ctx, indexModels)
	// Fail if error occurred
	if err != nil {
		log.Fatal(err)
	}
}
