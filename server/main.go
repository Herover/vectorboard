package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// Replace the uri string with your MongoDB deployment's connection string.
	uri := os.Getenv("MONGO_STR")
	// export MONGO_STR=mongodb+srv://doadmin:19i4UAm8CWK0325D@vectorboard-data-47dcec53.mongo.ondigitalocean.com/admin?authSource=admin\&tls=true\&tlsCAFile=./ca-certificate.crt
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
}
