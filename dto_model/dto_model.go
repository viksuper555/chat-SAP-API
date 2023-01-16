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
	HubId string `json:"hub_id,omitempty"`
}

type UserBody struct {
	model.User `json:"user,omitempty"`
	Type       string `json:"type,omitempty"`
}

type DataOnLogin struct {
	model.User  `json:"user,omitempty"`
	Type        string `json:"type,omitempty"`
	OnlineUsers []int  `json:"online_user_ids,omitempty"`
}

type ActiveUsersUpdate struct {
	Type      string `json:"type"`
	UserId    int    `json:"user_id"`
	Connected bool   `json:"connected"`
}
