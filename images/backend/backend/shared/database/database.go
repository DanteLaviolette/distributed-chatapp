package database

import (
	"context"
	"log"
	"time"

	"go.violettedev.com/eecs4222/shared/constants"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client = nil

const databaseInitializationTimeout = 10 * time.Second
const databaseName = "main"
const userCollectionName = "Users"
const refreshTokenCollectionName = "RefreshTokens"
const messagesCollectionName = "Messages"

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
Returns the messages collection
*/
func GetMessageCollection() mongo.Collection {
	return *(client.Database(databaseName).
		Collection(messagesCollectionName))
}

/*
Sets up the indices on the table if they don't yet exist
- db.users.createIndex( { "email": 1 }, { unique: true } )
- db.RefreshTokens.createIndex( expireAfterSeconds: constants.RefreshTokenExpirySeconds } )
*/
func setupIndices() {
	createUserIndices()
	createRefreshTokenIndices()
	createMessageIndex()
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

/*
Creates message indexes:
- db.messages.createIndex( { "ts": 1 } )
*/
func createMessageIndex() {
	// Define index
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"ts": 1, // index in ascending order
		},
	}
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(),
		databaseInitializationTimeout)
	defer cancel()
	// Add index to user collection
	col := GetMessageCollection()
	_, err := col.Indexes().CreateOne(ctx, indexModel)
	// Fail if error occurred
	if err != nil {
		log.Fatal(err)
	}
}
