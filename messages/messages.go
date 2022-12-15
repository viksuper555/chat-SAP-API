package messages

import (
	"context"
	"messenger/db"
	"messenger/graphql"
)

type Messages struct {
}

func (m Messages) Create(ctx context.Context, db db.Database, message *graphql.Message) error {
	return db.Create(message)
}

func LoadAll(ctx context.Context, db db.Database, userId string) (*[]graphql.Message, error) {
	messages, err := db.LoadUserMessages(userId)
	return messages, err
}
