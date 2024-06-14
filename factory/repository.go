package factory

import (
	"sync"

	"github.com/Yoshioka9709/yy-go-backend-template/repository"
)

// RepositoryFactory リポジトリの実装
type RepositoryFactory struct {
	setting *ServiceFactorySettings

	userRepo repository.User
	userOnce sync.Once

	todoRepo repository.Todo
	todoOnce sync.Once
}

// NewUser ユーザリポジトリを初期化
func (r *RepositoryFactory) NewUser() repository.User {
	r.userOnce.Do(func() {
		r.userRepo = repository.NewUser(r.setting.DynamoDB)
	})
	return r.userRepo
}

// NewTodo ユーザリポジトリを初期化
func (r *RepositoryFactory) NewTodo() repository.Todo {
	r.todoOnce.Do(func() {
		r.todoRepo = repository.NewTodo(r.setting.DynamoDB)
	})
	return r.todoRepo
}
