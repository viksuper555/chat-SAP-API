package model

import (
	"messenger/graph/customTypes"
	"time"
)

type base struct {
	CreatedAt time.Time  `gorm:"not null" json:"created_at,omitempty" `
	UpdatedAt time.Time  `gorm:"not null" json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"not null" json:"deleted_at,omitempty" `
}

type Message struct {
	base
	ID     int       `gorm:"primaryKey" json:"id" `
	Text   string    `gorm:"type:varchar(255)" json:"text,omitempty"`
	Date   time.Time `gorm:"not null" json:"date,omitempty" `
	UserID int
	User   *User `json:"user" gorm:"foreignkey:id"`
	Room   *Room `gorm:"type:varchar(255)" json:"room,omitempty"`
}

func (m *Message) ToGraph() *customTypes.Message {
	return &customTypes.Message{
		ID:     m.ID,
		Text:   m.Text,
		Date:   m.Date,
		UserID: m.UserID,
		RoomId: m.Room.ID,
	}
}

func MessagesToGraph(m []*Message) []*customTypes.Message {
	g := make([]*customTypes.Message, len(m))
	for i := range m {
		g[i] = m[i].ToGraph()
	}
	return g
}

type User struct {
	ID       int    `json:"id,omitempty" gorm:"primaryKey"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (u *User) ToGraph() *customTypes.User {
	return &customTypes.User{
		ID:   u.ID,
		Name: u.Username,
	}
}

func UsersToGraph(u []*User) []*customTypes.User {
	g := make([]*customTypes.User, len(u))
	for i := range u {
		g[i] = u[i].ToGraph()
	}
	return g
}

type Room struct {
	base
	ID    string  `gorm:"type:varchar(255)" json:"text,omitempty"`
	Users []*User `json:"user" gorm:"many2many:id"`
}

func (r *Room) ToGraph() *customTypes.Room {
	g := make([]int, len(r.Users))
	for i := range r.Users {
		g[i] = r.Users[i].ID
	}
	return &customTypes.Room{
		ID:      r.ID,
		UserIds: g,
	}
}

func RoomsToGraph(m []*Room) []*customTypes.Room {
	g := make([]*customTypes.Room, len(m))
	for i := range m {
		g[i] = m[i].ToGraph()
	}
	return g
}
