package client

import (
	"sync"
	"fmt"
	"errors"
)

var(
	ErrClientAlreadyExists = errors.New("")
)

type ClientError struct{
	ClientID string
	Message string
}

func(ce *ClientError) Error() string {
	return fmt.Sprintf("An error occured with client id %s: %s", ce.ClientID, ce.Message)
}

type Pool struct {
	Broadcast chan interface{}
	Clients   []*Instance
	sync.Mutex
}

func NewPool() *Pool {
	return &Pool{
		Broadcast: make(chan interface{}),
		Clients:   []*Instance{},
	}
}

func (c *Pool) NewClient(instance *Instance) error{
	if c.ClientExists(instance.ID){
		return &ClientError{ClientID:instance.ID, Message:"Client already exists"}
	}
	c.Mutex.Lock()
	c.Clients = append(c.Clients, instance)
	c.Mutex.Unlock()
	return nil
}

func (c *Pool) Start() {
	for {
		select {
		case message := <-c.Broadcast:
			for _, client := range c.Clients {
				client.message <- message
			}
		}
	}
}

func (c *Pool) GetClient(ID string) (*Instance, error){
	for _,instance := range c.Clients{
		if instance.ID == ID{
			return instance, nil
		}
	}
	return nil, &ClientError{ClientID: ID, Message: "Client not found"}
}

func (c *Pool) ClientExists(instanceID string) bool{
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	for _,client := range c.Clients{
		if client.ID == instanceID{
			return true
		}
	}
	return false
}
