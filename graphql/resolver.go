package graphql

//
//import (
//	"context"
//	"messenger/db"
//)
//
//type Resolver2 struct {
//	messages      []*Message
//	users         []*User
//	Database      db.Database
//	MessageServer MessageServer
//}
//
//type MessageServer interface {
//	LoadAll(ctx context.Context, db db.Database) ([]*Message, error)
//	Create(ctx context.Context, db db.Database, message *Message) error
//}
//
//// // foo
//func (r *mutationResolver) CreateMessage(ctx context.Context, input NewMessage) (*Message, error) {
//	panic("not implemented")
//}
//
//// // foo
//func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (string, error) {
//	panic("not implemented")
//}
//
//// // foo
//func (r *mutationResolver) Login(ctx context.Context, input Login) (string, error) {
//	panic("not implemented")
//}
//
//// // foo
//func (r *queryResolver) Messages(ctx context.Context) ([]*Message, error) {
//	msg, err := r.Database.LoadAllMessages()
//	return msg, err
//}
//
//// // foo
//func (r *queryResolver) Users(ctx context.Context) ([]*User, error) {
//	panic("not implemented")
//}
//
//// Mutation returns MutationResolver implementation.
//func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
//
//// Query returns QueryResolver implementation.
//func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }
