package Infrastructure

import (
	"context"
	"fmt"
	"log"
	"github.com/joho/godotenv"
	"os"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// initializing the working structure of mongo-db
func MongoDBInit() *mongo.Client {

	mongoURI := DotEnvLoader("MONGODB_URI") //extracts the secret URI from the .env

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}
func DotEnvLoader(identifier string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	result, exists := os.LookupEnv(identifier)

	if !exists {
		log.Fatal(".env entry doesn't exist")
	}

	return result
}