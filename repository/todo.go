package repository

import (
	"context"

	"github.com/Yoshioka9709/yy-go-backend-template/infra"
	"github.com/Yoshioka9709/yy-go-backend-template/model"
	"github.com/Yoshioka9709/yy-go-backend-template/util/errs"
	"github.com/Yoshioka9709/yy-go-backend-template/util/errs/codes"
	"github.com/guregu/dynamo"
)

// Todo Todoリポジトリのインターフェイス
type Todo interface {
	GetOne(ctx context.Context, id string) (*model.Todo, error)
	Find(ctx context.Context) ([]*model.Todo, error)
}

// todo Todoリポジトリの実装
type todo struct {
	dynamoDB *dynamo.DB
}

// NewTodo Todoリポジトリの初期化
func NewTodo(dynamoDB *dynamo.DB) Todo {
	return &todo{
		dynamoDB: dynamoDB,
	}
}

// GetOne ユーザ情報の取得
func (t *todo) GetOne(ctx context.Context, id string) (*model.Todo, error) {
	todosTable := t.dynamoDB.Table(infra.TableName(model.TodosTablePrefix))
	var Todo model.Todo
	err := todosTable.
		Get("PK", model.PKTodo).
		Range("ID", dynamo.Equal, id).
		OneWithContext(ctx, &Todo)
	if err != nil {
		return nil, errs.Wrap(codes.InternalError, err)
	}
	return &Todo, nil
}

// Find ユーザ情報の検索
func (t *todo) Find(ctx context.Context) ([]*model.Todo, error) {
	return nil, nil
}
