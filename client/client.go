package client

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrClientAlreadyExists = errors.New("")
)

type ClientError struct {
	ClientID string
	Message  string
}

func (ce ClientError) Error() string {
	return fmt.Sprintf("An error occured with client id %s: %s", ce.ClientID, ce.Message)
}

// A pool of clients
type Pool struct {
	Broadcast chan []byte
	clients   map[string]*Instance
	sync.Mutex
}

// Create a new pool of clients
func NewPool() *Pool {
	return &Pool{
		Broadcast: make(chan []byte),
		clients:   make(map[string]*Instance),
	}
}

// Add new client to pool
func (c *Pool) NewClient(instance *Instance) error {
	if c.ClientExists(instance.ID) {
		return &ClientError{ClientID: instance.ID, Message: "Client already exists"}
	}
	c.Mutex.Lock()
	c.clients[instance.ID] = instance
	c.Mutex.Unlock()
	return nil
}

// Delete a client from the pool
func (c *Pool) DeleteClient(instanceID string) error {
	if !c.ClientExists(instanceID) {
		return ClientError{ClientID: instanceID, Message: "Client not found"}
	}
	delete(c.clients, instanceID)
	return nil
}

// Start a pool -
func (c *Pool) Start() {
	for {
		select {
		case message := <-c.Broadcast:
			for _, client := range c.clients {
				client.message <- message
			}
		}
	}
}

// Get all clients in a pool
func (c *Pool) GetClients() map[string]*Instance {
	return c.clients
}

// Get a client from the pool using the client id
func (c *Pool) GetClient(instanceID string) (*Instance, error) {
	clients := c.GetClients()
	client, ok := clients[instanceID]
	if ok {
		return client, nil
	}
	return nil, ClientError{ClientID: instanceID, Message: "Client not found"}
}

// Check if a client exists in the pool using the client id
func (c *Pool) ClientExists(instanceID string) bool {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	_, ok := c.clients[instanceID]
	if ok {
		return true
	}
	return false
}
