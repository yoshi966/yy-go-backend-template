package service

import (
	"context"

	"github.com/Yoshioka9709/yy-go-backend-template/model"
	"github.com/Yoshioka9709/yy-go-backend-template/repository"
)

// User ユーザのインターフェース
type User interface {
	GetOne(ctx context.Context, id string) (*model.User, error)
	Find(ctx context.Context, filter *model.FindUserFilter) (*model.UserConnection, error)

	Create(ctx context.Context, input *model.CreateUserInput) (*model.User, error)
	Update(ctx context.Context, input *model.UpdateUserInput) (*model.User, error)
	Delete(ctx context.Context, input *model.DeleteUserInput) (*model.User, error)
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
func (u *user) Find(ctx context.Context, filter *model.FindUserFilter) (*model.UserConnection, error) {
	return u.userRepo.Find(ctx, filter)
}

// Create ユーザを作成
func (u *user) Create(ctx context.Context, input *model.CreateUserInput) (*model.User, error) {
	return u.userRepo.Create(ctx, input)
}

// Update ユーザーを更新
func (u *user) Update(ctx context.Context, input *model.UpdateUserInput) (*model.User, error) {
	return u.userRepo.Update(ctx, input)
}

// Delete ユーザーを削除
func (u *user) Delete(ctx context.Context, input *model.DeleteUserInput) (*model.User, error) {
	return u.userRepo.Delete(ctx, input)
}
