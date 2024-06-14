package model

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateUserInput struct {
	Name string `json:"name"`
}

type UpdateUserInput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DeleteUserInput struct {
	ID string `json:"id"`
}
