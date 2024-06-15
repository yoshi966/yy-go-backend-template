package repository

import (
	"context"

	"github.com/Yoshioka9709/yy-go-backend-template/dto"
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

	Create(ctx context.Context, user *model.User, input *model.CreateTodoInput) (*model.Todo, error)
	Update(ctx context.Context, input *model.UpdateTodoInput) (*model.Todo, error)
	Delete(ctx context.Context, input *model.DeleteTodoInput) (*model.Todo, error)
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

// GetOne Todo情報の取得
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

// Find Todo情報の検索
func (t *todo) Find(ctx context.Context) ([]*model.Todo, error) {
	return nil, nil
}

// Create Todoを作成
func (t *todo) Create(ctx context.Context, user *model.User, input *model.CreateTodoInput) (*model.Todo, error) {
	todosTable := t.dynamoDB.Table(infra.TableName(model.TodosTablePrefix))

	todo := model.NewTodo(user, input.Text)
	todoDTO := dto.Todo{
		PK:        model.PKTodo,
		ID:        todo.ID,
		Text:      todo.Text,
		Done:      todo.Done,
		UserID:    todo.User.ID,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}

	if err := todosTable.Put(todoDTO).RunWithContext(ctx); err != nil {
		return nil, err
	}
	return todo, nil
}

// Update Todoーを更新
func (t *todo) Update(ctx context.Context, input *model.UpdateTodoInput) (*model.Todo, error) {
	todosTable := t.dynamoDB.Table(infra.TableName(model.TodosTablePrefix))

	todo, err := t.GetOne(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	todo.Text = input.Text
	todo.Done = input.Done

	todoDTO := dto.Todo{
		PK:   model.PKTodo,
		ID:   todo.ID,
		Text: todo.Text,
		Done: todo.Done,
	}

	if err := todosTable.Put(todoDTO).RunWithContext(ctx); err != nil {
		return nil, err
	}
	return todo, nil
}

// Delete Todoーを削除
func (t *todo) Delete(ctx context.Context, input *model.DeleteTodoInput) (*model.Todo, error) {
	todosTable := t.dynamoDB.Table(infra.TableName(model.TodosTablePrefix))

	todo, err := t.GetOne(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	todoDTO := dto.Todo{
		PK: model.PKTodo,
		ID: todo.ID,
	}

	if err := todosTable.Delete("PK", string(todoDTO.PK)).Range("ID", todoDTO.ID).RunWithContext(ctx); err != nil {
		return nil, err
	}
	return todo, nil
}
