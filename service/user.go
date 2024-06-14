package service

import (
	"context"

	"github.com/Yoshioka9709/yy-go-backend-template/model"
	"github.com/Yoshioka9709/yy-go-backend-template/repository"
)

// User ユーザのインターフェース
type User interface {
	GetOne(ctx context.Context, id string) (*model.User, error)
	Find(ctx context.Context) ([]*model.User, error)
}

// user ユーザのサービス実装
type user struct {
	userRepo repository.User
}

// NewUsers ユーザを初期化
func NewUsers(ur repository.User) User {
	u := &user{
		userRepo: ur,
	}
	return u
}

// GetOne ユーザ情報の取得
func (u *user) GetOne(ctx context.Context, id string) (*model.User, error) {
	return u.userRepo.GetOne(ctx, id)
}

// Find ユーザ情報の検索
func (u *user) Find(ctx context.Context) ([]*model.User, error) {
	return u.userRepo.Find(ctx)
}
