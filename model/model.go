package model

import (
	"messenger/graph/customTypes"
	"time"
)

//type base struct {
//	CreatedAt time.Time  `gorm:"not null" json:"created_at,omitempty" `
//	UpdatedAt time.Time  `gorm:"not null" json:"updated_at,omitempty"`
//	DeletedAt *time.Time `gorm:"not null" json:"deleted_at,omitempty" `
//}

// Message belongs to `User`, `UserID` is the foreign key
type Message struct {
	ID     int       `gorm:"primaryKey" json:"id" `
	Text   string    `gorm:"type:varchar(255)" json:"text,omitempty"`
	Date   time.Time `gorm:"not null" json:"date,omitempty" `
	UserID int
	User   User `json:"user"`
	RoomID string
	Room   Room `gorm:"constraint:OnDelete:CASCADE;" json:"room,omitempty"`
}

type User struct {
	ID       int    `json:"id,omitempty" gorm:"primaryKey"`
	Username string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
	Rooms    []Room `gorm:"many2many:user_room; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"rooms"`
}

type Room struct {
	ID       string     `json:"id" gorm:"type:varchar(255); primaryKey"`
	Name     string     `json:"name,omitempty"`
	Users    []*User    `gorm:"many2many:user_room; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"users"`
	Messages []*Message `gorm:"foreignKey:RoomID; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"messages"`
}

func (m *Message) ToGraph() *customTypes.Message {
	return &customTypes.Message{
		ID:     m.ID,
		Text:   m.Text,
		Date:   m.Date,
		UserID: m.UserID,
		RoomID: m.RoomID,
	}
}

func MessagesToGraph(m []*Message) []*customTypes.Message {
	g := make([]*customTypes.Message, len(m))
	for i := range m {
		g[i] = m[i].ToGraph()
	}
	return g
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

func (r *Room) ToGraph() *customTypes.Room {
	uIds := make([]int, len(r.Users))
	for i := range r.Users {
		uIds[i] = r.Users[i].ID
	}
	return &customTypes.Room{
		ID:       r.ID,
		UserIds:  uIds,
		Users:    UsersToGraph(r.Users),
		Name:     r.Name,
		Messages: MessagesToGraph(r.Messages),
	}
}

func RoomsToGraph(m []*Room) []*customTypes.Room {
	g := make([]*customTypes.Room, len(m))
	for i := range m {
		g[i] = m[i].ToGraph()
	}
	return g
}
