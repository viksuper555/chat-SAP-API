package graphql

import (
	"context"
	"messenger/db"
	"messenger/graphql/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type MessageServer interface {
	LoadAll(ctx context.Context, db db.Database) (*[]Message, error)
	Create(ctx context.Context, db db.Database, message *Message) error
}

type Resolver struct {
	messages      []*model.Message
	users         []*model.User
	Database      db.Database
	MessageServer MessageServer
}
