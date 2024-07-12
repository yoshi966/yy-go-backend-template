package model

import (
	"time"

	"github.com/Yoshioka9709/yy-go-backend-template/util"
)

// Todo Todo
type Todo struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Done      bool      `json:"done"`
	User      *User     `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateTodoInput Todo 作成のインプット
type CreateTodoInput struct {
	Text   string `json:"text" validate:"required"`
	UserID string `json:"userId" validate:"required"`
}

// UpdateTodoInput Todo 更新のインプット
type UpdateTodoInput struct {
	ID     string `json:"id" validate:"required"`
	Text   string `json:"text" validate:"required"`
	Done   bool   `json:"done"`
	UserID string `json:"userId" validate:"required"`
}

// DeleteTodoInput Todo 削除のインプット
type DeleteTodoInput struct {
	ID     string `json:"id" validate:"required"`
	UserID string `json:"userId" validate:"required"`
}

// FindTodoFilter Todo 検索のフィルター
type FindTodoFilter struct {
	Paging *DataPage
}

// TodoConnection Todo 検索結果
type TodoConnection struct {
	Edges      []*TodoEdge `json:"edges"`
	PageInfo   *PageInfo   `json:"pageInfo"`
	TotalCount int         `json:"totalCount"`
}

// TodoEdge Todo 検索エッジ
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
