package dto

import (
	"time"

	"github.com/Yoshioka9709/yy-go-backend-template/model"
)

type Todo struct {
	PK        model.PK
	ID        string
	Text      string
	Done      bool
	UserID    string
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
