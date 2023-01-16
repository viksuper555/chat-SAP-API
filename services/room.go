package services

import (
	"log"
	"messenger/model"
	"sync"
)

var Rm = NewRoom()

type Room struct {
	UMap  map[int]*model.User // ID: User
	mutex sync.Mutex
}

func NewRoom() *Room {
	r := new(Room)
	r.UMap = make(map[int]*model.User)
	return r
}

func (r *Room) LoginUser(u *model.User) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	log.Printf("User logged in: %d\n", u.ID)
	r.UMap[u.ID] = u
}

func (r *Room) LogoutUser(id int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	log.Printf("User offline: %d\n", id)
	delete(r.UMap, id)
}
