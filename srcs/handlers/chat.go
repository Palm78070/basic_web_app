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
	clients map[string]*websocket.Conn
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

	defer func() {
		fmt.Println("Closing connection")
		conn.Close()
		delete(app.chat.clients, app.currentUser.username)
	}()

	if app.currentUser.username == "" {
		fmt.Println("User is not logged in")
		return
	}

	fmt.Println("Connect to websocket: ", app.currentUser.username)

	app.mutex.Lock()
	if _, ok := app.chat.clients[app.currentUser.username]; !ok {
		app.chat.clients[app.currentUser.username] = conn
	}
	app.mutex.Unlock()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading msg: ", err)
			return
		}
		app.chat.broadcast <- msg
	}
}

func (app *App) HandleMessages() {
	app.wg.Add(1)
	defer app.wg.Done()
	for {
		fmt.Println("HandleMessages")
		msg, ok := <-app.chat.broadcast
		if !ok {
			fmt.Println("Broadcast channel is closed")
			return
		}

		users := []string{msg.Sender, msg.Recipient}
		sort.Strings(users)
		room_name := users[0] + "_" + users[1]

		app.mutex.Lock()
		if app.chat.rooms != nil {
			if _, ok := app.chat.rooms[room_name]; !ok {
				app.chat.rooms[room_name] = users
			}
		}
		app.mutex.Unlock()

		for _, user := range app.chat.rooms[room_name] {
			if app.chat.rooms == nil {
				fmt.Println("room map is nil")
				break
			}
			app.mutex.Lock()
			client, ok := app.chat.clients[user]
			app.mutex.Unlock()
			if !ok {
				fmt.Println("User is not connected to websocket")
				continue
			}
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println("xxx")
				fmt.Println(err)
				client.Close()
			}
		}
	}
}
