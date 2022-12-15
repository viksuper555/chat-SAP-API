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
	User   *User     `json:"user"`
	//Timestamp int64  `gorm:"not null" json:"timestamp,omitempty" `
}

type User struct {
	Id       string          `json:"id,omitempty" gorm:"primaryKey"`
	Name     string          `json:"name,omitempty"`
	Password string          `json:"password,omitempty"`
	Ch       chan Message    `json:"-"`
	Ws       *websocket.Conn `json:"-"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type NewMessage struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
