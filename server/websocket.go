package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type wsRoom struct {
	Connections map[int]*websocket.Conn
}

func (room wsRoom) broadcast(msg []byte) error {
	for _, conn := range room.Connections {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			//return err
		}
	}
	return nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var rooms = map[string]wsRoom{}

func wsMessageHandler(c chan []byte, ws *websocket.Conn) {
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		c <- p
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()

	vars := mux.Vars(r)
	roomID := vars["room_id"]
	connectionIndex := -1
	for i := 0; i < 1000; i++ {
		if _, ok := rooms[roomID].Connections[i]; !ok {
			connectionIndex = i
			break
		}
	}
	if connectionIndex == -1 {
		log.Print("Could not find a open index for client")
		return
	}
	if _, ok := rooms[roomID]; !ok {
		rooms[roomID] = wsRoom{
			Connections: map[int]*websocket.Conn{},
		}
	}
	rooms[roomID].Connections[connectionIndex] = ws
	room := rooms[roomID]
	fmt.Printf("%s %d\n", roomID, connectionIndex)
	defer func() {
		delete(rooms[roomID].Connections, connectionIndex)
	}()

	messageChan := make(chan []byte)
	closeChan := make(chan int)

	go wsMessageHandler(messageChan, ws)

	// Make sure to clean up if connection closes
	oldCloseHandler := ws.CloseHandler()
	ws.SetCloseHandler(func(code int, text string) error {
		if err := oldCloseHandler(code, text); err != nil {
			log.Print("error", err)
			closeChan <- code
			return err
		}
		closeChan <- code
		return nil
	})
	for {
		select {
		case msgData := <-messageChan:

			if err := ws.WriteMessage(websocket.TextMessage, msgData); err != nil {
				log.Println(err)
				return
			}

			room.broadcast(msgData)

		case <-closeChan:
			return
		}
	}
}
