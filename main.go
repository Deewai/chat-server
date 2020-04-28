package main

import (
	"log"
	"net/http"

	"github.com/Deewai/chat-server/client"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// Create a new pool of clients
var pool = client.NewPool()

func main() {
	go pool.Start()
	router := mux.NewRouter()
	// Handle function /:id where id represents client's identification
	router.HandleFunc("/{id}", handleSocketConnections)
	log.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleSocketConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		handleWebsocketError(err)
	}
	defer ws.Close()
	// Get the client id from the request
	params := mux.Vars(r)
	clientID := params["id"]
	if clientID == "" {
		clientID = uuid.New().String()
	}
	// Get a new instance of a client
	conn := client.NewConn(ws, pool, clientID)
	// Add client to client pool
	pool.NewClient(conn)
	// Start go routine to send message to clients
	go conn.SendMessage()
	// Read new messages
	conn.Read()

}

func handleWebsocketError(err error) {
	if e, ok := err.(*websocket.CloseError); ok && (e.Code == websocket.CloseNormalClosure || e.Code == websocket.CloseNoStatusReceived) {
		log.Println("Websocket connection closed")
	}

}
