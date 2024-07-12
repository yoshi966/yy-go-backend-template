package model

import "github.com/Yoshioka9709/yy-go-backend-template/util"

// User ユーザー
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateUserInput ユーザー作成のインプット
type CreateUserInput struct {
	Name string `json:"name" validate:"required"`
}

// UpdateUserInput ユーザー更新のインプット
type UpdateUserInput struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

// DeleteUserInput ユーザー削除のインプット
type DeleteUserInput struct {
	ID string `json:"id" validate:"required"`
}

// FindUserFilter ユーザー検索のフィルター
type FindUserFilter struct {
	Paging *DataPage
}

// UserConnection ユーザー検索結果
type UserConnection struct {
	Edges      []*UserEdge `json:"edges"`
	PageInfo   *PageInfo   `json:"pageInfo"`
	TotalCount int         `json:"totalCount"`
}

// UserEdge ユーザー検索エッジ
type UserEdge struct {
	Cursor Cursor `json:"cursor"`
	Node   *User  `json:"node"`
}

// Set エッジに値を設定
func (e *UserEdge) Set(cursor Cursor, node any) {
	e.Cursor = cursor
	e.Node = node.(*User)
}

// NewUser ユーザー初期化
func NewUser(name string) *User {
	return &User{
		ID:   util.GetULIDString(),
		Name: name,
	}
}
