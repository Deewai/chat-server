package client

import (
	// "flag"
	// "fmt"
	// "github.com/gorilla/mux"
	"encoding/json"
	"sync"

	"log"

	"github.com/gorilla/websocket"
	// "net/http"
)

type Instance struct {
	ID       string
	conn     *websocket.Conn
	message  chan interface{}
	sync.Mutex
	LastMessages map[string]string
	Pool *Pool
}

func NewConn(ws *websocket.Conn, pool *Pool, ID string) *Instance {
	return &Instance{
		ID:           ID,
		conn:         ws,
		message:      make(chan interface{}),
		LastMessages: map[string]string{},
		Pool: pool,
	}
}

func (i *Instance) Read() {
	defer func() {
		i.closeConnection()
	}()

	for {

		// Read in a new message as JSON and map it to a Message object
		_, message, err := i.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		} else {
			i.processMessage(message)
		}
	}
}

func (i *Instance) closeConnection() {
	err := i.conn.Close()
	if err != nil {
		log.Printf("Error closing instance connection: %v", err)
	}
}

func (i *Instance) processMessage(message []byte) {
	var request map[string]interface{}
	err := json.Unmarshal(message, &request)
	if err != nil {
		//process message as a custom message
		log.Println(err)
	} else {
		switch action,_ := request["action"].(string); action {
		case "register":
			//register client
			log.Printf("action is %s", action)
			log.Printf("client id is %s", i.ID)
			log.Printf("client remote id is %s", i.ID)
		case "message":
			//get recipient from request payload
			recipient := request["recipient"].(string)
			client, err := i.Pool.GetClient(recipient)
			if err != nil{
				log.Printf("An error occured while getting client: %v", err)
			} else {
				//send message to receiver
				client.message <- request
			}
		
		case "broadcast":
			//send message to receiver
			i.Pool.Broadcast <- request
		default:
			// send to all clients
			for _, client := range i.Pool.Clients{
				client.message <- request
			}
		}
	}
}

func (i *Instance) SendMessage(){
	for {
		select {
		case msg, _ := <-i.message:
			// Grab the next message from the broadcast channel
			log.Println(msg)
			// Send it out to instance client that is currently connected
			i.conn.WriteJSON(msg)
		}
	}
}