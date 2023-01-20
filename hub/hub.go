package hub

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"messenger/internal/common"
)

var MainHub = newHub()
var Rooms = map[string]*Room{}

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
	Clients map[int]*Client
	// Open Rooms.
	Broadcast chan interface{}
	// Register requests from the Clients.
	register chan *Client

	// Unregister requests from Clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		Broadcast:  make(chan interface{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[int]*Client),
	}
}
func initRooms() {
	dbRooms, err := common.GetRooms(common.Db)
	if err != nil {
		log.Fatal(err)
	}
	for _, dbRoom := range dbRooms {
		err := InitRoom(dbRoom.ID)
		if err != nil {
			log.Printf("Error starting room %s %s", dbRoom.ID, err.Error())
		}
	}
}

func InitRoom(id string) error {
	if id == "" {
		id = uuid.NewV4().String()
	}

	r := &Room{
		uuid:       id,
		Broadcast:  make(chan interface{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
	Rooms[r.uuid] = r
	go r.Run()
	return nil
}
func ConnectUserToRoom(uId int, rId string) error {
	client, ok := MainHub.Clients[uId]
	if !ok {
		return fmt.Errorf("User not connected")
	}
	Rooms[rId].register <- client

	return nil
}

func (h *Hub) Run() {
	go initRooms()
	for {
		select {
		case client := <-h.register:
			h.Clients[client.user.ID] = client
		case client := <-h.unregister:
			if _, ok := h.Clients[client.user.ID]; ok {
				delete(h.Clients, client.user.ID)
				close(client.send)
			}
		case message := <-h.Broadcast:
			for _, client := range h.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.Clients, client.user.ID)
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
