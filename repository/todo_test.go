package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Yoshioka9709/yy-go-backend-template/infra"
	"github.com/Yoshioka9709/yy-go-backend-template/model"
	"github.com/Yoshioka9709/yy-go-backend-template/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_todo_GetOne(t *testing.T) {
	testutil.InitDynamoDB(t)

	type args struct {
		ID string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Todo
		wantErr bool
	}{
		{
			"success",
			args{
				ID: "dummytodoid1001",
			},
			&model.Todo{
				ID:        "dummytodoid1001",
				Text:      "これはテストの投稿です",
				Done:      false,
				CreatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			d := infra.NewDynamoDBClient()
			td := NewTodo(d)

			got, err := td.GetOne(ctx, tt.args.ID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
