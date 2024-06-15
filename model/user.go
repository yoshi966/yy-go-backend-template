package model

import "github.com/Yoshioka9709/yy-go-backend-template/util"

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateUserInput struct {
	Name string `json:"name" validate:"required"`
}

type UpdateUserInput struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type DeleteUserInput struct {
	ID string `json:"id" validate:"required"`
}

type FindUserFilter struct {
	Paging *DataPage
}

type UserConnection struct {
	Edges      []*UserEdge `json:"edges"`
	PageInfo   *PageInfo   `json:"pageInfo"`
	TotalCount int         `json:"totalCount"`
}

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
