package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoService struct {
	MongoClient *mongo.Client
	Database    *string
}

var databaseName = "numbernama"

func (mc *MongoService) DisconnectMongoClient(ctx context.Context) error {
	fmt.Print("mongo-client disconnect has been called,")
	if err := mc.MongoClient.Disconnect(ctx); err != nil {
		return err
	}
	fmt.Println(" and closed")
	return nil
}

func ConnectMongoDB(ctx context.Context) (*MongoService, error) {
	// create a context with timeOut
	cCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// fmt.Println("mongo connect", *config.MongoURI, *secrets.MongodbUsername, *secrets.MongodbPassword)

	// // Use the SetServerAPIOptions() method to set the Stable API version to 1
	credentials := options.Credential{
		AuthMechanism: "SCRAM-SHA-1",
		Username:      "user-admin",
		Password:      "pass-secrets",
	}

	opts := options.
		Client().
		ApplyURI("mongodb://localhost:27017").
		SetAuth(credentials)

	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		fmt.Println("\nerror while mongo connect", "Error", err.Error())
		return nil, err
	}

	// // create a Client with Encryption
	// client, err := mongo.NewClientEncryption(client, opts)

	// Send a ping to confirm a successful connection
	if err := client.
		Database(databaseName).
		RunCommand(cCtx, bson.D{{Key: "ping", Value: 1}}).
		Err(); err != nil {
		fmt.Println("\n", "Error", err.Error())
		return nil, err
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return &MongoService{
		MongoClient: client,
		Database:    &databaseName,
	}, nil
}
