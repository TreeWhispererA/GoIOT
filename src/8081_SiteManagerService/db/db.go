package db

import (
	"context"
	"log"

	"github.com/fatih/color"
	mongo "github.com/helios/go-sdk/proxy-libs/heliosmongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	middlewares "tracio.com/sitemanagerservice/handlers"
)

var client *mongo.Client

// Dbconnect -> connects mongo
func Dbconnect() *mongo.Client {

	// serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	// opts := options.Client().ApplyURI(middlewares.DotEnvVariable("MONGO_ATLAS_URL")).SetServerAPIOptions(serverAPI).SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	// // // Create a new client and connect to the server
	// client, err := mongo.Connect(context.TODO(), opts)

	clientOptions := options.Client().ApplyURI(middlewares.DotEnvVariable("MONGO_URL"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("⛒ Connection Failed to Database")
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("⛒ Connection Failed to Database")
		log.Fatal(err)
	}
	color.Green("⛁ Connected to Database")
	return client
}
