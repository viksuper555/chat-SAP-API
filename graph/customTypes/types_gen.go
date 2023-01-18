// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package customTypes

import (
	"time"
)

type Message struct {
	ID     int       `json:"id"`
	Text   string    `json:"text"`
	Date   time.Time `json:"date"`
	User   *User     `json:"user"`
	UserID int       `json:"user_id"`
	RoomID string    `json:"room_id"`
}

type NewMessage struct {
	UserID int    `json:"userId"`
	Text   string `json:"text"`
	RoomID string `json:"roomId"`
}

type NewRoom struct {
	ID      string `json:"id"`
	UserIds []int  `json:"user_ids"`
}

type Room struct {
	ID      string  `json:"id"`
	Users   []*User `json:"users"`
	UserIds []int   `json:"user_ids"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserPass struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
