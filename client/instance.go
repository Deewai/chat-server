package client

import (
	"sync"

	"log"

	"github.com/gorilla/websocket"
)

// interface to be implemented by the pool passed to a client
type pool interface {
	DeleteClient(string) error
	GetClients() map[string]*Instance
}

// Instance type which is a client
type Instance struct {
	ID      string
	conn    *websocket.Conn
	message chan []byte
	close   chan interface{}
	sync.Mutex
	LastMessages map[string]string
	Pool         pool
}

// Create a new connection/client
func NewConn(ws *websocket.Conn, pool pool, ID string) *Instance {
	return &Instance{
		ID:           ID,
		conn:         ws,
		message:      make(chan []byte),
		close:        make(chan interface{}),
		LastMessages: map[string]string{},
		Pool:         pool,
	}
}

func (i *Instance) Read() {
	defer func() {
		i.closeConnection()
	}()

	for {

		// Read in a new message byte array
		_, message, err := i.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			if _, ok := err.(*websocket.CloseError); ok {
				i.close <- 1
				err = i.Pool.DeleteClient(i.ID)
				if err != nil {
					log.Println(err)
				}
			}
			// break
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
	// Close client SendMessage routine
	i.close <- 1
	// Delete client from pool
	err = i.Pool.DeleteClient(i.ID)
	if err != nil {
		log.Println(err)
	}

}

func (i *Instance) processMessage(message []byte) {
	// Send to all clients
	for _, client := range i.Pool.GetClients() {
		client.message <- message
	}
}

func (i *Instance) SendMessage() {
	for {
		select {
		case msg, _ := <-i.message:
			// Grab the next message from the message channel
			log.Println(msg)
			// Send it out to instance client that is currently connected
			i.conn.WriteMessage(websocket.TextMessage, msg)
		case <-i.close:
			return
		}

	}
}
