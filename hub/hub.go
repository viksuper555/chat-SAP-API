package hub

import (
	uuid "github.com/satori/go.uuid"
)

var Hub1 = newHub()

// Room maintains the set of active Clients and broadcasts messages to the
// Clients.
type Room struct {
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

type Hub struct {
	// Registered Clients.
	Clients map[*Client]bool
	// Open Rooms.
	Rooms map[string]*Room

	// Register requests from the Clients.
	register chan *Client

	// Unregister requests from Clients.
	unregister chan *Client
}

func newHub() *Hub {
	room := newRoom()
	room.uuid = "global"
	return &Hub{
		register:   room.register,
		unregister: room.unregister,
		Clients:    room.Clients,
		Rooms:      map[string]*Room{room.uuid: room},
	}
}

func newRoom() *Room {
	return &Room{
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
		case message := <-h.Rooms["global"].Broadcast:
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
