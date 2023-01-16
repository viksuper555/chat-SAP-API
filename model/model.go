package model

import (
	"messenger/graph/customTypes"
	"time"
)

type base struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Message struct {
	base
	ID     int       `gorm:"primaryKey" json:"id" `
	Type   string    `gorm:"-" json:"type"`
	Text   string    `gorm:"type:varchar(255)" json:"text,omitempty"`
	Date   time.Time `gorm:"not null" json:"date,omitempty" `
	UserID int
	User   *User `json:"user" gorm:"foreignkey:id"`
	//Timestamp int64  `gorm:"not null" json:"timestamp,omitempty" `
}

func (m *Message) ToGraph() *customTypes.Message {
	return &customTypes.Message{
		ID:     m.ID,
		Text:   m.Text,
		Date:   m.Date,
		UserID: m.UserID,
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
	//Ch       chan Message    `json:"-" gorm:"-"`
	//Ws       *websocket.Conn `json:"-" gorm:"-"`
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
