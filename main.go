package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// ############# web socket ################
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { // allow cors
		return true
	},
}

var (
	user     []*websocket.Conn
	userLock sync.Mutex // Mutex to synchronize access to the 'user' slice
)

func addUser(conn *websocket.Conn) {
	userLock.Lock()
	defer userLock.Unlock()
	user = append(user, conn)
}

func removeUser(conn *websocket.Conn) {
	userLock.Lock()
	defer userLock.Unlock()
	for i, u := range user {
		if u == conn {
			user = append(user[:i], user[i+1:]...)
			return
		}
	}
}

func reader(conn *websocket.Conn) {
	addUser(conn)
	defer func() {
		removeUser(conn)
		conn.Close()
	}()

	for {
		messageType, p, err := conn.ReadMessage() // read the message sent
		if err != nil {
			log.Println(err)
			return
		}

		// broadcast the message
		userLock.Lock()
		for _, u := range user {
			if u != conn {
				err := u.WriteMessage(messageType, p) // send the messages
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
		userLock.Unlock()
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Client Connected")
	reader(ws)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func setUpRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("server is running")
	setUpRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
