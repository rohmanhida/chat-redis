package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// websocket conn
var Conn *websocket.Conn

// connected clients
var (
	Clients      = make(map[*websocket.Conn]bool)
	ClientsMutex = sync.Mutex{}
)

func GetClientCount(w http.ResponseWriter, r *http.Request) int {

	// count connected clients
	ClientsMutex.Lock()
	count := len(Clients)
	ClientsMutex.Unlock()

	return count
}

func websocketInit(w http.ResponseWriter, r *http.Request) {

	// websocket upgrader
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	// init upgrade
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade error: ", err)
	}
	Conn = conn
}

func addClient() {

	// add client to client mutex
	ClientsMutex.Lock()
	Clients[Conn] = true
	ClientsMutex.Unlock()
}

func removeClient(conn *websocket.Conn) {

	// add client to client mutex
	ClientsMutex.Lock()
	delete(Clients, conn)
	ClientsMutex.Unlock()
}
