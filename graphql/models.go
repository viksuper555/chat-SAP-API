package graphql

import (
	"golang.org/x/net/websocket"
	"time"
)

type base struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Message struct {
	base
	ID     uint      `gorm:"primaryKey" json:"id" `
	Type   string    `gorm:"type:varchar(100)" json:"type"`
	Sender string    `gorm:"type:varchar(100)" json:"sender_id,omitempty"`
	Text   string    `gorm:"type:varchar(255)" json:"text,omitempty"`
	Date   time.Time `gorm:"not null" json:"date,omitempty" `
	UserID int
	User   *User `json:"user"`
	//Timestamp int64  `gorm:"not null" json:"timestamp,omitempty" `
}

type User struct {
	ID       string          `json:"id,omitempty" gorm:"primaryKey"`
	Name     string          `json:"name,omitempty"`
	Password string          `json:"password,omitempty"`
	Ch       chan Message    `json:"-" gorm:"-"`
	Ws       *websocket.Conn `json:"-" gorm:"-"`
}
