package graph

import (
	"context"
	"poc/graph/common"
	"poc/graph/customTypes"
	"poc/graph/generated"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, text string) (*customTypes.Todo, error) {
	context := common.GetContext(ctx)
	todo := &customTypes.Todo{
		Text: text,
		Done: false,
	}
	err := context.Database.Create(&todo).Error
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, input customTypes.TodoInput) (*customTypes.Todo, error) {
	context := common.GetContext(ctx)
	todo := &customTypes.Todo{
		ID:   input.ID,
		Text: input.Text,
		Done: input.Done,
	}
	err := context.Database.Save(&todo).Error
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, todoID int) (*customTypes.Todo, error) {
	context := common.GetContext(ctx)
	var todo *customTypes.Todo
	err := context.Database.Where("id = ?", todoID).Delete(&todo).Error
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *queryResolver) GetTodos(ctx context.Context) ([]*customTypes.Todo, error) {
	context := common.GetContext(ctx)
	var todos []*customTypes.Todo
	err := context.Database.Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *queryResolver) GetTodo(ctx context.Context, todoID int) (*customTypes.Todo, error) {
	context := common.GetContext(ctx)
	var todo *customTypes.Todo
	err := context.Database.Where("id = ?", todoID).Find(&todo).Error
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
