package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
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
		writeError(w, err)
		panic(err)
	}
	json.NewEncoder(w).Encode(httpResponse{
		Data: boards,
	})
}

func getBoardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	boardID := vars["id"]

	var board Board
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	docID, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		writeError(w, err)
		return
	}
	result := boardsCollection.FindOne(ctx, bson.M{"$and": bson.A{
		bson.M{"_id": docID},
		bson.M{"deleted": false},
	}})
	if result.Err() != nil {
		writeError(w, result.Err())
		return
	}
	if err := result.Decode(&board); err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(httpResponse{
		Data: board,
	})
}

func newBoardHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var newBoard Board

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newBoard)

	newBoard.Content = make([]BoardText, 0)

	result, err := boardsCollection.InsertOne(ctx, newBoard)
	if err != nil {
		writeError(w, err)
		return
	}

	json.NewEncoder(w).Encode(httpResponse{
		Data: result.InsertedID,
	})
}

func deleteBoardHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	boardID := vars["id"]
	docID, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		writeError(w, err)
		return
	}

	result, err := boardsCollection.DeleteOne(ctx, bson.M{"_id": docID})
	if err != nil {
		writeError(w, err)
		return
	}
	if result.DeletedCount != 1 {
		writeError(w, fmt.Errorf("nothing deleted"))
		return
	}

	json.NewEncoder(w).Encode(httpResponse{
		Data: "OK",
	})
}

func updateBoardHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var newBoard Board

	vars := mux.Vars(r)
	boardID := vars["id"]
	docID, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		writeError(w, err)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newBoard)

	// Only allow updating some fields
	result, err := boardsCollection.UpdateByID(ctx, docID, bson.M{
		"$set": bson.M{
			"hidden": newBoard.Hidden,
		},
	})
	if err != nil {
		writeError(w, err)
		return
	}
	if result.ModifiedCount != 1 {
		writeError(w, fmt.Errorf("nothing updated"))
		return
	}

	json.NewEncoder(w).Encode(httpResponse{
		Data: "OK",
	})
}

func newBoardContentHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	boardID := vars["id"]

	var msg BoardMessage
	reqBody, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(reqBody, &msg); err != nil {
		writeError(w, err)
		return
	}

	/* if o, ok := d.Msg.(BoardText); ok {

	} */

	docID, err := primitive.ObjectIDFromHex(boardID)
	if err != nil {
		writeError(w, err)
		return
	}

	if msg.Action == "post" {
		result, err := boardsCollection.UpdateByID(ctx, docID, bson.M{"$push": bson.M{"content": msg.Data.Msg}})
		if err != nil {
			writeError(w, err)
			return
		}
		if result.ModifiedCount != 1 {
			writeError(w, fmt.Errorf("did not add new content"))
			return
		}
	} /* else if msg.Action == "update" {
		result, err := boardsCollection.UpdateByID(ctx, docID, bson.M{
			"content": bson.M{"$elemMatch": bson.M{"_id": msg.Data.Msg} msg.Data.Msg},
		})
		if err != nil {
			writeError(w, err)
			return
		}
		if result.ModifiedCount != 1 {
			writeError(w, fmt.Errorf("did not add new content"))
			return
		}
	} */

	json.NewEncoder(w).Encode(httpResponse{
		Data: "OK",
	})
}
