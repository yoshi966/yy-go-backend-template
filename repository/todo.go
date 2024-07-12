package repository

import (
	"context"

	"github.com/Yoshioka9709/yy-go-backend-template/dto"
	"github.com/Yoshioka9709/yy-go-backend-template/infra"
	"github.com/Yoshioka9709/yy-go-backend-template/model"
	"github.com/Yoshioka9709/yy-go-backend-template/util"
	"github.com/Yoshioka9709/yy-go-backend-template/util/errs"
	"github.com/Yoshioka9709/yy-go-backend-template/util/errs/codes"
	"github.com/guregu/dynamo"
)

// Todo Todoリポジトリのインターフェイス
type Todo interface {
	GetOne(ctx context.Context, id string) (*model.Todo, error)
	Find(ctx context.Context, filter *model.FindTodoFilter) (*model.TodoConnection, error)

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
	var todo model.Todo
	err := todosTable.
		Get("PK", model.PKTodo).
		Range("ID", dynamo.Equal, id).
		OneWithContext(ctx, &todo)
	if err != nil {
		return nil, errs.Wrap(codes.InternalError, err)
	}
	return &todo, nil
}

// Find Todo情報の検索
// nolint: dupl
func (t *todo) Find(ctx context.Context, filter *model.FindTodoFilter) (*model.TodoConnection, error) {
	todosTable := t.dynamoDB.Table(infra.TableName(model.TodosTablePrefix))
	query := todosTable.Get("PK", model.PKTodo)

	pager, err := model.NewForwardPager(filter.Paging, func(_, offset int) (any, error) {
		var todos []*model.Todo

		// 検索 最大件数は欲しい値+offset
		err := query.Order(dynamo.Descending).
			Limit(int64(filter.Paging.First+offset)).
			AllWithContext(ctx, &todos)
		if err != nil {
			return nil, errs.Wrap(codes.InvalidParameter, err)
		}
		return todos[offset:], nil
	})
	if err != nil {
		return nil, errs.Wrap(codes.InternalError, err)
	}

	// 件数取得
	totalCount, err := query.CountWithContext(ctx)
	if err != nil {
		return nil, errs.Wrap(codes.InvalidParameter, err)
	}

	result := &model.TodoConnection{
		Edges:      pager.Edges(&model.TodoEdge{}).([]*model.TodoEdge),
		PageInfo:   pager.PageInfo(int(totalCount)),
		TotalCount: int(totalCount),
	}
	return result, nil
}

// Create Todoを作成
func (t *todo) Create(ctx context.Context, user *model.User, input *model.CreateTodoInput) (*model.Todo, error) {
	todosTable := t.dynamoDB.Table(infra.TableName(model.TodosTablePrefix))

	now := util.GetTimeNow()
	todo := model.NewTodo(user, input.Text, now)
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
