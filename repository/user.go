package repository

import (
	"context"

	"github.com/Yoshioka9709/yy-go-backend-template/dto"
	"github.com/Yoshioka9709/yy-go-backend-template/infra"
	"github.com/Yoshioka9709/yy-go-backend-template/model"
	"github.com/Yoshioka9709/yy-go-backend-template/util"
	"github.com/Yoshioka9709/yy-go-backend-template/util/errs"
	"github.com/Yoshioka9709/yy-go-backend-template/util/errs/codes"
	"github.com/guregu/dynamo"
)

// User ユーザーリポジトリのインターフェイス
type User interface {
	GetOne(ctx context.Context, id string) (*model.User, error)
	Find(ctx context.Context, filter *model.FindUserFilter) (*model.UserConnection, error)

	Create(ctx context.Context, input *model.CreateUserInput) (*model.User, error)
	Update(ctx context.Context, input *model.UpdateUserInput) (*model.User, error)
	Delete(ctx context.Context, input *model.DeleteUserInput) (*model.User, error)
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
// nolint: dupl
func (u *user) Find(ctx context.Context, filter *model.FindUserFilter) (*model.UserConnection, error) {
	usersTable := u.dynamoDB.Table(infra.TableName(model.UsersTablePrefix))
	query := usersTable.Get("PK", model.PKUser)

	pager, err := model.NewForwardPager(filter.Paging, func(_, offset int) (any, error) {
		var users []*model.User

		// 検索 最大件数は欲しい値+offset
		err := query.Order(dynamo.Descending).
			Limit(int64(filter.Paging.First+offset)).
			AllWithContext(ctx, &users)
		if err != nil {
			return nil, errs.Wrap(codes.InvalidParameter, err)
		}
		return users[offset:], nil
	})
	if err != nil {
		return nil, errs.Wrap(codes.InternalError, err)
	}

	// 件数取得
	totalCount, err := query.CountWithContext(ctx)
	if err != nil {
		return nil, errs.Wrap(codes.InvalidParameter, err)
	}

	result := &model.UserConnection{
		Edges:      pager.Edges(&model.UserEdge{}).([]*model.UserEdge),
		PageInfo:   pager.PageInfo(int(totalCount)),
		TotalCount: int(totalCount),
	}
	return result, nil
}

// Create ユーザを作成
func (u *user) Create(ctx context.Context, input *model.CreateUserInput) (*model.User, error) {
	usersTable := u.dynamoDB.Table(infra.TableName(model.UsersTablePrefix))

	user := model.NewUser(input.Name)

	now := util.GetTimeNow()
	userDTO := &dto.User{
		PK:        model.PKUser,
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := usersTable.Put(userDTO).RunWithContext(ctx); err != nil {
		return nil, err
	}
	return user, nil
}

// Update ユーザーを更新
func (u *user) Update(ctx context.Context, input *model.UpdateUserInput) (*model.User, error) {
	usersTable := u.dynamoDB.Table(infra.TableName(model.UsersTablePrefix))

	user, err := u.GetOne(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	user.Name = input.Name

	userDTO := dto.User{
		PK:   model.PKUser,
		ID:   user.ID,
		Name: user.Name,
	}

	if err := usersTable.Put(userDTO).RunWithContext(ctx); err != nil {
		return nil, err
	}
	return user, nil
}

// Delete ユーザーを削除
func (u *user) Delete(ctx context.Context, input *model.DeleteUserInput) (*model.User, error) {
	usersTable := u.dynamoDB.Table(infra.TableName(model.UsersTablePrefix))

	user, err := u.GetOne(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	userDTO := dto.User{
		PK: model.PKUser,
		ID: user.ID,
	}

	if err := usersTable.Delete("PK", string(userDTO.PK)).Range("ID", userDTO.ID).RunWithContext(ctx); err != nil {
		return nil, err
	}
	return user, nil
}
