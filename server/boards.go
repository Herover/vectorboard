package main

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const boardsCollectionName = "boards"

type BoardObject interface {
	IsBoardObject()
}

type BoardText struct {
	Type   string `json:"type"`
	String string `json:"string"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

func (BoardText) IsBoardObject() {}

type BoardLine struct {
	Type string `json:"type"`
	X1   int    `json:"x1"`
	Y1   int    `json:"y1"`
	X2   int    `json:"x2"`
	Y2   int    `json:"y2"`
}

func (BoardLine) IsBoardObject() {}

type Board struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name    string             `json:"name"`
	Hidden  bool               `json:"hidden"`
	Deleted bool               `json:"deleted"`
	Content []BoardText        `json:"content"`
}

type NewBoard struct {
	Name string
}

type BoardMessage struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Action string             `json:"action"`
	Data   MessageContent     `json:"data"`
}

type MessageContent struct {
	Msg BoardObject
}

func (d *MessageContent) UnmarshalJSON(data []byte) error {
	var meta struct {
		Type string
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return err
	}
	switch meta.Type {
	case "text":
		d.Msg = &BoardText{}
	case "line":
		d.Msg = &BoardLine{}
	default:
		return fmt.Errorf("%q is an invalid item type", meta.Type)
	}

	return json.Unmarshal(data, d.Msg)
}
