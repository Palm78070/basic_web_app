package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Chat struct {
	upgrader *websocket.Upgrader
	clients map[*websocket.Conn]bool
	broadcast chan Message
}

type Message struct {
	// `` is a struct tag used when marshaling and unmarshaling Go structs to JSON
	Username string `json:"username"`
	Message string `json:"message"`
}

func (app *App) HandleConnections(w http.ResponseWriter, r *http.Request) {
	//Upgrade() return *Conn, error
	conn, err := app.chat.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close() //Close connection when it's not needed anymore

	app.chat.clients[conn] = true //assign value of *Conn to true to indicate that connection is open
	fmt.Println("Connect to websocket")

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(app.chat.clients, conn)
			return
		}
		app.chat.broadcast <- msg
	}
}

func (app *App) HandleMessages() {
	for {
		msg := <-app.chat.broadcast
		fmt.Println("user:", msg.Username, "message: ", msg.Message)

		for client := range app.chat.clients {
			err := client.WriteJSON(msg) //send a JSON-encoded message over a WebSocket
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(app.chat.clients, client)
			}
		}
	}
}
