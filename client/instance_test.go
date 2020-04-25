package client

import (
	// "reflect"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	// "github.com/posener/wstest"
)

func TestNewConn(t *testing.T) {
	conn := NewConn(&websocket.Conn{}, &Pool{}, "testID")
	assert.ObjectsAreEqual(&Instance{}, conn)
}

func TestInstance_Read(t *testing.T) {
	conn := &websocket.Conn{}
	instance := &Instance{
		ID:           "id123",
		conn:         conn,
		message:      make(chan interface{}),
		LastMessages: map[string]string{},
	}
	go instance.Read()
	
}

func TestInstance_closeConnection(t *testing.T) {
	type fields struct {
		ID           string
		conn         *websocket.Conn
		message      chan interface{}
		Mutex        sync.Mutex
		LastMessages map[string]string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Instance{
				ID:           tt.fields.ID,
				conn:         tt.fields.conn,
				message:      tt.fields.message,
				Mutex:        tt.fields.Mutex,
				LastMessages: tt.fields.LastMessages,
			}
			i.closeConnection()
		})
	}
}

func TestInstance_processMessage(t *testing.T) {
	type fields struct {
		ID           string
		conn         *websocket.Conn
		message      chan interface{}
		Mutex        sync.Mutex
		LastMessages map[string]string
	}
	type args struct {
		message []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Instance{
				ID:           tt.fields.ID,
				conn:         tt.fields.conn,
				message:      tt.fields.message,
				Mutex:        tt.fields.Mutex,
				LastMessages: tt.fields.LastMessages,
			}
			i.processMessage(tt.args.message)
		})
	}
}
