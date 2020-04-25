package main

import (
	// "flag"
	// "fmt"
	// "encoding/json"
	"log"
	"net/http"

	"github.com/Deewai/chat-server/client"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var pool = client.NewPool()

func main() {
	go pool.Start()
	router := mux.NewRouter()
	router.HandleFunc("/", handleSocketConnections)
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
	params := mux.Vars(r)
	userID := params["user_id"]
	if userID == "" {
		userID = uuid.New().String()
	}
	conn := client.NewConn(ws, pool, userID)
	pool.NewClient(conn)
	go conn.SendMessage()
	conn.Read()

}

func handleWebsocketError(err error) {
	if e, ok := err.(*websocket.CloseError); ok && (e.Code == websocket.CloseNormalClosure || e.Code == websocket.CloseNoStatusReceived) {
		log.Println("Websocket connection closed")
	}

}

// func NewConn(w http.ResponseWriter, r *http.Request) (*Conn, error) {
// 	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		handleWebsocketError(err)
// 	}
// 	defer ws.Close()

// }
