package dto

import (
	"time"

	"github.com/Yoshioka9709/yy-go-backend-template/model"
)

// User DTO
type User struct {
	PK        model.PK
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
