package dto_model

import "messenger/graphql"

type MessageBody struct {
	Id         uint          `json:"id"`
	Type       string        `json:"type"`
	Sender     string        `json:"sender_id,omitempty"`
	Message    string        `json:"message,omitempty"`
	Timestamp  int64         `json:"timestamp,omitempty"`
	User       *graphql.User `json:"user"`
	Recipients []string      `json:"recipients,omitempty"`
}

type UserBody struct {
	graphql.User `json:"dto_models,omitempty"`
	Type         string `json:"type,omitempty"`
}

type BroadcastOnline struct {
	Type  string   `json:"type"`
	Users []string `json:"users"`
}
