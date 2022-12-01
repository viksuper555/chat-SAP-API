package main

import (
	"golang.org/x/net/websocket"
	"sync"
)

type Message struct {
	Type      string `json:"type"`
	Sender    string `json:"sender_id,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type MessageBody struct {
	Message
	Recipients []string `json:"recipients,omitempty"`
}

type User struct {
	uuid string
	name string
	ch   chan Message
	ws   *websocket.Conn
}

type RegisterBody struct {
	Uuid     string `json:"uuid,omitempty"`
	Username string `json:"username,omitempty"`
	Type     string `json:"type,omitempty"`
}

type Room struct {
	uMap  map[string]*User // uuid: User
	mutex sync.Mutex
}

type BroadcastOnline struct {
	Type  string   `json:"type"`
	Users []string `json:"users"`
}
