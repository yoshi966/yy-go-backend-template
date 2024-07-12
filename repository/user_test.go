package repository

import (
	"context"
	"testing"

	"github.com/Yoshioka9709/yy-go-backend-template/infra"
	"github.com/Yoshioka9709/yy-go-backend-template/model"
	"github.com/Yoshioka9709/yy-go-backend-template/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_user_GetOne(t *testing.T) {
	testutil.InitDynamoDB(t)

	type args struct {
		ID string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			"success",
			args{
				ID: "dummyuserid0001",
			},
			&model.User{
				ID:   "dummyuserid0001",
				Name: "テストユーザー1",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.TODO()
			d := infra.NewDynamoDBClient()
			u := NewUser(d)

			got, err := u.GetOne(ctx, tt.args.ID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
