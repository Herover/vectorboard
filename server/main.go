package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client
var boardsCollection *mongo.Collection

func main() {
	// Replace the uri string with your MongoDB deployment's connection string.
	uri := os.Getenv("MONGO_STR")
	// export MONGO_STR=mongodb+srv://doadmin:19i4UAm8CWK0325D@vectorboard-data-47dcec53.mongo.ondigitalocean.com/?authSource=admin\&tls=true\&tlsCAFile=./ca-certificate.crt
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
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

	database := client.Database("admin")

	// Confirm collections exist, or create new
	collectionNames, err := database.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	boardsExists := false
	for _, collectionName := range collectionNames {
		if collectionName == boardsCollectionName {
			boardsExists = true
		}
	}
	if !boardsExists {
		if err = database.CreateCollection(ctx, boardsCollectionName); err != nil {
			panic(err)
		}
		log.Print("Created boards collection")
	}
	boardsCollection = database.Collection(boardsCollectionName)

	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)
	r.HandleFunc("/boards", getBoardsHandler).Methods(http.MethodGet)
	r.HandleFunc("/boards", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/boards", newBoardHandler).Methods(http.MethodPost)
	r.HandleFunc("/boards/{id}", getBoardHandler).Methods(http.MethodGet)
	r.HandleFunc("/boards/{id}", updateBoardHandler).Methods(http.MethodPut)
	r.HandleFunc("/boards/{id}", deleteBoardHandler).Methods(http.MethodDelete)
	r.HandleFunc("/boards/{id}", newBoardContentHandler).Methods(http.MethodPost)
	r.HandleFunc("/boards/{id}", optionsHandler).Methods(http.MethodOptions)
	r.HandleFunc("/ws/{room_id}", websocketHandler)
	var corsMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		})
	}
	r.Use(corsMiddleware)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Listening for connections")

	log.Fatal(srv.ListenAndServe())
}
