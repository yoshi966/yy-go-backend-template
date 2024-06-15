package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.48

import (
	"context"

	"github.com/Yoshioka9709/yy-go-backend-template/graphql/generated"
	"github.com/Yoshioka9709/yy-go-backend-template/model"
)

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return getService(ctx).NewUser().GetOne(ctx, id)
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, paging *model.DataPage) (*model.UserConnection, error) {
	return getService(ctx).NewUser().Find(ctx, &model.FindUserFilter{
		Paging: paging,
	})
}

// Todo is the resolver for the todo field.
func (r *queryResolver) Todo(ctx context.Context, id string) (*model.Todo, error) {
	return getService(ctx).NewTodo().GetOne(ctx, id)
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context, paging *model.DataPage) (*model.TodoConnection, error) {
	return getService(ctx).NewTodo().Find(ctx, &model.FindTodoFilter{
		Paging: paging,
	})
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
