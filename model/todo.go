package model

import (
	"time"

	"github.com/Yoshioka9709/yy-go-backend-template/util"
)

type Todo struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Done      bool      `json:"done"`
	User      *User     `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateTodoInput struct {
	Text   string `json:"text" validate:"required"`
	UserID string `json:"userId" validate:"required"`
}

type UpdateTodoInput struct {
	ID     string `json:"id" validate:"required"`
	Text   string `json:"text" validate:"required"`
	Done   bool   `json:"done"`
	UserID string `json:"userId" validate:"required"`
}

type DeleteTodoInput struct {
	ID     string `json:"id" validate:"required"`
	UserID string `json:"userId" validate:"required"`
}

type FindTodoFilter struct {
	Paging *DataPage
}

type TodoConnection struct {
	Edges      []*TodoEdge `json:"edges"`
	PageInfo   *PageInfo   `json:"pageInfo"`
	TotalCount int         `json:"totalCount"`
}

type TodoEdge struct {
	Cursor Cursor `json:"cursor"`
	Node   *Todo  `json:"node"`
}

// Set エッジに値を設定
func (e *TodoEdge) Set(cursor Cursor, node any) {
	e.Cursor = cursor
	e.Node = node.(*Todo)
}

// NewTodo Todo初期化
func NewTodo(user *User, text string, now time.Time) *Todo {
	return &Todo{
		ID:        util.GetULIDString(),
		Text:      text,
		Done:      false,
		User:      user,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
