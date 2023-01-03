package model

import (
	"golang.org/x/net/websocket"
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
	Type   string    `gorm:"type:varchar(100)" json:"type"`
	Sender int       `json:"sender_id,omitempty"`
	Text   string    `gorm:"type:varchar(255)" json:"text,omitempty"`
	Date   time.Time `gorm:"not null" json:"date,omitempty" `
	UserID uint
	User   *User `json:"user" gorm:"foreignkey:id"`
	//Timestamp int64  `gorm:"not null" json:"timestamp,omitempty" `
}

func (m *Message) ToGraph() *customTypes.Message {
	return &customTypes.Message{
		ID:     m.ID,
		Text:   m.Text,
		Date:   m.Date,
		UserID: m.Sender,
	}
}

type User struct {
	Id       int             `json:"id,omitempty" gorm:"primaryKey"`
	Name     string          `json:"name,omitempty"`
	Password string          `json:"password,omitempty"`
	Ch       chan Message    `json:"-" gorm:"-"`
	Ws       *websocket.Conn `json:"-" gorm:"-"`
}
