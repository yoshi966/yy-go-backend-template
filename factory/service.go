package factory

// Service はサービスレジストリ
type Service interface {
	// TODO: RepositoryFactory がここにいるのは謎なのでどうにかする
	RepositoryFactory() *RepositoryFactory
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
