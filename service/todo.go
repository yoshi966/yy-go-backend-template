package service

import (
	"context"

	"github.com/Yoshioka9709/yy-go-backend-template/model"
	"github.com/Yoshioka9709/yy-go-backend-template/repository"
)

// Todo Todoのインターフェース
type Todo interface {
	GetOne(ctx context.Context, id string) (*model.Todo, error)
	Find(ctx context.Context) ([]*model.Todo, error)

	Create(ctx context.Context, input *model.CreateTodoInput) (*model.Todo, error)
	Update(ctx context.Context, input *model.UpdateTodoInput) (*model.Todo, error)
	Delete(ctx context.Context, input *model.DeleteTodoInput) (*model.Todo, error)
}

// todo Todoのサービス実装
type todo struct {
	todoRepo repository.Todo
	userRepo repository.User
}

// NewTodos Todoを初期化
func NewTodos(tr repository.Todo, ur repository.User) Todo {
	t := &todo{
		todoRepo: tr,
		userRepo: ur,
	}
	return t
}

// GetOne Todo情報の取得
func (t *todo) GetOne(ctx context.Context, id string) (*model.Todo, error) {
	return t.todoRepo.GetOne(ctx, id)
}

// Find Todo情報の検索
func (t *todo) Find(ctx context.Context) ([]*model.Todo, error) {
	return t.todoRepo.Find(ctx)
}

// Create Todoを作成
func (t *todo) Create(ctx context.Context, input *model.CreateTodoInput) (*model.Todo, error) {
	user, err := t.userRepo.GetOne(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return t.todoRepo.Create(ctx, user, input)
}

// Update Todoーを更新
func (t *todo) Update(ctx context.Context, input *model.UpdateTodoInput) (*model.Todo, error) {
	return t.todoRepo.Update(ctx, input)
}

// Delete Todoーを削除
func (t *todo) Delete(ctx context.Context, input *model.DeleteTodoInput) (*model.Todo, error) {
	return t.todoRepo.Delete(ctx, input)
}
