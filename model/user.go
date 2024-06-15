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

// NewUser ユーザー初期化
func NewUser(name string) *User {
	return &User{
		ID:   util.GetULIDString(),
		Name: name,
	}
}
