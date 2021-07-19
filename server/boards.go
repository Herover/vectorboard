package main

type Board struct {
	ID      string `json:"_id"`
	Name    string `json:"name"`
	Hidden  bool   `json:"hidden"`
	Deleted bool   `json:"deleted"`
}

type NewBoard struct {
	Name string
}

const boardsCollectionName = "boards"
