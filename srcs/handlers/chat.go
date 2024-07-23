package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Chat struct {
	upgrader *websocket.Upgrader
	clients map[*websocket.Conn]bool
	dummy map[string]*websocket.Conn
	broadcast chan Message
	rooms map[string][]string
}

type Message struct {
	// `` is a struct tag used when marshaling and unmarshaling Go structs to JSON
	Sender string `json:"sender"`
	Recipient string `json:"recipient"`
	Message string `json:"message"`
}

func (app *App) HandleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In HandleConnections")
	session, err := app.SessionStore.Get(r, "session-name")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !app.session_exist(session) {
		fmt.Println("User is not logged in")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	app.session_map_user(session)
	fmt.Printf("Session values: %v\n", session.Values["IsLogin"])
	fmt.Printf("Current user state: %+v\n", app.currentUser)


	path_var := mux.Vars(r)
	room_type := path_var["room_type"]

	if room_type != "private" && room_type != "group" {
		fmt.Println("Invalid room type")
		http.Error(w, "Invalid form submission", http.StatusBadRequest)
		return
	}
	//Upgrade() return *Conn, error
	conn, err := app.chat.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer conn.Close() //Close connection when it's not needed anymore

	defer func() {
		fmt.Println("Closing connection")
		conn.Close()
	}()

	app.chat.clients[conn] = true //assign value of *Conn to true to indicate that connection is open

	if app.currentUser.username == "" {
		fmt.Println("User is not logged in")
		return
	}

	fmt.Println("Connect to websocket: ", app.currentUser.username)
	// var msg_dummy Message
	// conn.ReadJSON(&msg_dummy)
	// fmt.Println("Connect to websocket: ", msg_dummy)

	if _, ok := app.chat.dummy[app.currentUser.username]; !ok {
		app.chat.dummy[app.currentUser.username] = conn
	}
	for {
		fmt.Println("In HandleConnections for loop")
		var msg Message
		err := conn.ReadJSON(&msg)
		fmt.Println("msg is ", msg)
		if msg == (Message{}) {
			fmt.Println("Empty message")
			return
		}
		if err != nil {
			fmt.Println("Error reading msg: ", err)
			delete(app.chat.clients, conn)
			return
		}
		app.chat.broadcast <- msg
	}
}

func (app *App) HandleMessages() {
	fmt.Println("In HandleMessages")
	for {
		fmt.Println("In HandleMessages for loop")
		msg := <-app.chat.broadcast
		fmt.Printf("Sender: %v\n", msg.Sender)
		fmt.Printf("Recipient: %v\n", msg.Recipient)
		fmt.Printf("Message: %v\n", msg.Message)
		users := []string{msg.Sender, msg.Recipient}
		sort.Strings(users)
		fmt.Println("sort users is ", users)
		room_name := users[0] + "_" + users[1]
		fmt.Println("room_name is ", room_name)

		if _, ok := app.chat.rooms[room_name]; !ok {
			app.chat.rooms[room_name] = users
		}

		for _, user := range app.chat.rooms[room_name] {
			fmt.Println("User is ", user)
			if _, ok := app.chat.dummy[user]; !ok {
				fmt.Println("User is not connected to websocket")
				continue
			}
			client := app.chat.dummy[user]
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(app.chat.dummy, user)
			}
		}

		// for client := range app.chat.clients {
		// 	err := client.WriteJSON(msg) //send a JSON-encoded message over a WebSocket
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		client.Close()
		// 		delete(app.chat.clients, client)
		// 	}
		// }
	}
}
