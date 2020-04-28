package client

import (
	"fmt"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestNewPool(t *testing.T) {
	pool := NewPool()
	assert.NotNil(t, pool)
}

func TestPool_NewClientSuccessfull(t *testing.T) {
	pool := &Pool{
		Broadcast: make(chan interface{}),
		Clients:   []*Instance{},
		Mutex:     sync.Mutex{},
	}
	err := pool.NewClient(&Instance{})
	assert.NotNil(t, pool)
	assert.Len(t, pool.Clients, 1)
	assert.Nil(t, err)
}

func TestPool_NewClientAlreadyExists(t *testing.T) {
	t.Run("Adding a duplicate client should throw an error of type *ClientError{}", func(t *testing.T) {
		pool := &Pool{
			Broadcast: make(chan interface{}),
			Clients:   []*Instance{},
			Mutex:     sync.Mutex{},
		}
		id := "testID"
		pool.NewClient(&Instance{ID: id})
		err := pool.NewClient(&Instance{ID: id})
		assert.NotNil(t, err)
		assert.IsType(t, &ClientError{}, err)
		assert.EqualValues(t, fmt.Sprintf("An error occured with client id %s: %s", id, "Client already exixts"), err.Error())
	})

}

func TestPool_Start(t *testing.T) {
	pool := &Pool{
		Broadcast: make(chan interface{}),
		Clients:   []*Instance{},
		Mutex:     sync.Mutex{},
	}
	client := &Instance{
		ID:           "test123",
		conn:         &websocket.Conn{},
		message:      make(chan interface{}),
		Mutex:        sync.Mutex{},
		LastMessages: map[string]string{},
	}
	pool.NewClient(client)
	go pool.Start()
	newMessage := "a new message"
	pool.Broadcast <- newMessage
	assert.EqualValues(t, newMessage, <-client.message)
}
