package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type health struct {
	Mongo string `json:"mongo"`
}

type Error struct {
	Msg string `json:"msg"`
}

type httpResponse struct {
	Errors []string    `json:"errors"`
	Data   interface{} `json:"data"`
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(Error{
		Msg: err.Error(),
	})
	log.Print(err.Error())
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	mongoStatus := "OK"
	code := 200

	if client == nil {
		mongoStatus = "unconnected"
		code = 500
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			mongoStatus = err.Error()
			code = 500
		}
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(health{
		Mongo: mongoStatus,
	})
}

func getBoardsHandler(w http.ResponseWriter, r *http.Request) {
	var boards []Board
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := boardsCollection.Find(ctx, bson.M{"$and": bson.A{
		bson.M{"deleted": false},
		bson.M{"hidden": false},
	}})
	if err != nil {
		writeError(w, err)
		return
	}
	if err = cursor.All(ctx, &boards); err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(httpResponse{
		Data: boards,
	})
}

func newBoardHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var newBoard Board

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newBoard)

	fmt.Printf("%+v\n", newBoard)

	result, err := boardsCollection.InsertOne(ctx, newBoard)
	if err != nil {
		writeError(w, err)
		return
	}

	json.NewEncoder(w).Encode(httpResponse{
		Data: result.InsertedID,
	})
}
