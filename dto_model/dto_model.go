package dto_model

import (
	"messenger/model"
)

type MessageBody struct {
	Id         int    `json:"id"`
	Type       string `json:"type"`
	SenderId   int    `json:"sender_id,omitempty"`
	SenderName string `json:"sender_name,omitempty"`
	Message    string `json:"message,omitempty"`
	Timestamp  int64  `json:"timestamp,omitempty"`
	//User       *model.User `json:"user"`
	Recipients []int `json:"recipients,omitempty"`
}

type UserBody struct {
	model.User `json:"user,omitempty"`
	Type       string `json:"type,omitempty"`
}

type BroadcastOnline struct {
	Type  string   `json:"type"`
	Users []string `json:"users"`
}
