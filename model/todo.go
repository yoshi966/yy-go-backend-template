package model

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
	User *User  `json:"user"`
}

type CreateTodoInput struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

type UpdateTodoInput struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Done   bool   `json:"done"`
	UserID string `json:"userId"`
}

type DeleteTodoInput struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
}