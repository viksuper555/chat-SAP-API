package hub

import (
	uuid "github.com/satori/go.uuid"
)

var MainHub = newHub()

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

	Broadcast chan interface{}
	// Register requests from the Clients.
	register chan *Client

	// Unregister requests from Clients.
	unregister chan *Client
}

func newHub() *Hub {
	g := newRoom()
	g.uuid = "global"
	vp := newRoom()
	vp.uuid = "vp"
	return &Hub{
		Broadcast:  make(chan interface{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Rooms:      map[string]*Room{g.uuid: g, vp.uuid: vp},
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
	for i := range h.Rooms {
		go h.Rooms[i].Run()
	}
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

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.Clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
			}
		case message := <-r.Broadcast:
			for client := range r.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.Clients, client)
				}
			}
		}
	}
}
