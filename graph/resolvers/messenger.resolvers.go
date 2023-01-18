package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.22

import (
	"context"
	"messenger/graph/customTypes"
	"messenger/graph/generated"
	"messenger/model"
	"time"

	"github.com/satori/go.uuid"
)

// CreateMessage is the resolver for the createMessage field.
func (r *mutationResolver) CreateMessage(ctx context.Context, input customTypes.NewMessage) (*customTypes.Message, error) {
	msg := &model.Message{
		Text:   input.Text,
		UserID: input.UserID,
		RoomID: input.RoomID,
		Date:   time.Now(),
	}
	err := r.DB.Create(&msg).Error
	if err != nil {
		return nil, err
	}

	res := &customTypes.Message{
		ID:   msg.ID,
		Text: input.Text,
		Date: time.Now(),
	}
	return res, nil
}

// CreateRoom is the resolver for the createRoom field.
func (r *mutationResolver) CreateRoom(ctx context.Context, input customTypes.NewRoom) (*customTypes.Room, error) {
	var users []*model.User
	err := r.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}

	if input.ID == "" {
		input.ID = uuid.NewV4().String()
	}
	rm := &model.Room{
		ID:    input.ID,
		Users: users,
	}

	err = r.DB.Create(&rm).Error
	if err != nil {
		return nil, err
	}

	res := &customTypes.Room{
		ID:    rm.ID,
		Users: model.UsersToGraph(rm.Users),
	}
	return res, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input customTypes.UserPass) (*customTypes.User, error) {
	//context := common.GetContext(ctx)
	user := &model.User{
		Username: input.Username,
		Password: input.Password,
	}
	err := r.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	res := &customTypes.User{
		ID:   user.ID,
		Name: user.Username,
	}
	return res, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input customTypes.UserPass) (bool, error) {
	var user model.User
	var exists bool
	err := r.DB.Model(&user).
		Select("count(*) > 0").
		Where("name = ? AND password = ?", input.Username, input.Password).
		Find(&exists).
		Error

	if err != nil {
		return exists, err
	}
	if !exists {
		return exists, err
	}
	return exists, nil
}

// GetMessages is the resolver for the getMessages field.
func (r *queryResolver) GetMessages(ctx context.Context) ([]*customTypes.Message, error) {
	var messages []*model.Message
	err := r.DB.Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return model.MessagesToGraph(messages), nil
}

// GetUserMessages is the resolver for the getUserMessages field.
func (r *queryResolver) GetUserMessages(ctx context.Context, userID int) ([]*customTypes.Message, error) {
	var messages []*model.Message
	err := r.DB.Where(&model.Message{UserID: userID}).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return model.MessagesToGraph(messages), nil
}

// GetMessage is the resolver for the getMessage field.
func (r *queryResolver) GetMessage(ctx context.Context, messageID int) (*customTypes.Message, error) {
	var msg model.Message
	err := r.DB.Where(&model.Message{ID: messageID}).First(&msg).Error
	if err != nil {
		return nil, err
	}
	return msg.ToGraph(), nil
}

// GetUsers is the resolver for the getUsers field.
func (r *queryResolver) GetUsers(ctx context.Context) ([]*customTypes.User, error) {
	var users []*model.User
	err := r.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return model.UsersToGraph(users), nil
}

// GetUser is the resolver for the getUser field.
func (r *queryResolver) GetUser(ctx context.Context, userID int) (*customTypes.User, error) {
	var u model.User
	err := r.DB.Where(&model.User{ID: userID}).First(&u).Error
	if err != nil {
		return nil, err
	}
	return u.ToGraph(), nil
}

// GetRooms is the resolver for the getRooms field.
func (r *queryResolver) GetRooms(ctx context.Context, userID int) ([]*customTypes.Room, error) {
	var rooms []*model.Room

	err := r.DB.Preload("Users").Where("id IN (SELECT room_id FROM user_room WHERE user_id = ?)", userID).Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return model.RoomsToGraph(rooms), nil
}

// GetRoom is the resolver for the getRoom field.
func (r *queryResolver) GetRoom(ctx context.Context, roomID string) (*customTypes.Room, error) {
	var room model.Room
	err := r.DB.Where(&model.Room{ID: roomID}).Preload("Users").First(&room).Error
	if err != nil {
		return nil, err
	}
	return room.ToGraph(), nil
}

// GetRoomMessages is the resolver for the getRoomMessages field.
func (r *queryResolver) GetRoomMessages(ctx context.Context, roomID string) ([]*customTypes.Message, error) {
	var messages []*model.Message
	err := r.DB.Where(&model.Message{RoomID: roomID}).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return model.MessagesToGraph(messages), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
