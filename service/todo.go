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
}

// todo Todoのサービス実装
type todo struct {
	todoRepo repository.Todo
}

// NewTodos Todoを初期化
func NewTodos(ur repository.Todo) Todo {
	u := &todo{
		todoRepo: ur,
	}
	return u
}

// GetOne Todo情報の取得
func (u *todo) GetOne(ctx context.Context, id string) (*model.Todo, error) {
	return u.todoRepo.GetOne(ctx, id)
}

// Find Todo情報の検索
func (u *todo) Find(ctx context.Context) ([]*model.Todo, error) {
	return u.todoRepo.Find(ctx)
}
