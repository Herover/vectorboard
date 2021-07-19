package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type Board struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name    string             `json:"name"`
	Hidden  bool               `json:"hidden"`
	Deleted bool               `json:"deleted"`
}

type NewBoard struct {
	Name string
}

const boardsCollectionName = "boards"
