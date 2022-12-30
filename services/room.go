package services

import (
	"log"
	"messenger/graphql"
	"sync"
)

var Rm = NewRoom()

type Room struct {
	UMap  map[string]*graphql.User // ID: User
	mutex sync.Mutex
}

func NewRoom() *Room {
	r := new(Room)
	r.UMap = make(map[string]*graphql.User)
	return r
}

func (r *Room) LoginUser(u *graphql.User) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	log.Printf("User logged in: %s\n", u.ID)
	r.UMap[u.ID] = u
}

func (r *Room) LogoutUser(id string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	log.Printf("User offline: %s\n", id)
	delete(r.UMap, id)
}
