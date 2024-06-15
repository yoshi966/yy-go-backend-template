package factory

import (
	"github.com/Yoshioka9709/yy-go-backend-template/service"
)

// Service はサービスレジストリ
type Service interface {
	// TODO: RepositoryFactory がここにいるのは謎なのでどうにかする
	RepositoryFactory() *RepositoryFactory

	NewUser() service.User
	NewTodo() service.Todo
}

// serviceFactory はサービスレジストリの実装
type serviceFactory struct {
	repository *RepositoryFactory
}

// NewService インフラ層の依存情報を初期化時に注入する
func NewService(setting *ServiceFactorySettings) Service {
	return &serviceFactory{
		repository: &RepositoryFactory{
			setting: setting,
		},
	}
}

// RepositoryFactory リポジトリを返す
func (s *serviceFactory) RepositoryFactory() *RepositoryFactory {
	return s.repository
}

// NewUser ユーザサービスを返す
func (s *serviceFactory) NewUser() service.User {
	userRepo := s.repository.NewUser()
	return service.NewUsers(userRepo)
}

// NewTodo Todoサービスを返す
func (s *serviceFactory) NewTodo() service.Todo {
	todoRepo := s.repository.NewTodo()
	userRepo := s.repository.NewUser()
	return service.NewTodos(todoRepo, userRepo)
}
