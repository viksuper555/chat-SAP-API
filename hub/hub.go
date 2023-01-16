package hub

import (
	uuid "github.com/satori/go.uuid"
)

var Hub1 = newHub()
var Hubs = make(map[string]Hub)

// Hub maintains the set of active Clients and broadcasts messages to the
// Clients.
type Hub struct {
	uuid string
	// Registered Clients.
	Clients map[*Client]bool

	// Inbound messages from the Clients.
	Broadcast chan interface{}

	// Register requests from the Clients.
	register chan *Client

	// Unregister requests from Clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		uuid:       uuid.NewV4().String(),
		Broadcast:  make(chan interface{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.Clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
