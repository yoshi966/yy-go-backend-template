package repository

import (
	"context"

	"github.com/Yoshioka9709/yy-go-backend-template/infra"
	"github.com/Yoshioka9709/yy-go-backend-template/model"
	"github.com/Yoshioka9709/yy-go-backend-template/util/errs"
	"github.com/Yoshioka9709/yy-go-backend-template/util/errs/codes"
	"github.com/guregu/dynamo"
)

// User ユーザーリポジトリのインターフェイス
type User interface {
	GetOne(ctx context.Context, id string) (*model.User, error)
	Find(ctx context.Context) ([]*model.User, error)
}

// user ユーザーリポジトリの実装
type user struct {
	dynamoDB *dynamo.DB
}

// NewUser ユーザーリポジトリの初期化
func NewUser(dynamoDB *dynamo.DB) User {
	return &user{
		dynamoDB: dynamoDB,
	}
}

// GetOne ユーザ情報の取得
func (u *user) GetOne(ctx context.Context, id string) (*model.User, error) {
	usersTable := u.dynamoDB.Table(infra.TableName(model.UsersTablePrefix))
	var user model.User
	err := usersTable.
		Get("PK", model.PKUser).
		Range("ID", dynamo.Equal, id).
		OneWithContext(ctx, &user)
	if err != nil {
		return nil, errs.Wrap(codes.InternalError, err)
	}
	return &user, nil
}

// Find ユーザ情報の検索
func (u *user) Find(ctx context.Context) ([]*model.User, error) {
	return nil, nil
}
