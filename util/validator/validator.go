package validator

import (
	"context"

	"github.com/Yoshioka9709/yy-go-backend-template/util/errs"
	"github.com/Yoshioka9709/yy-go-backend-template/util/errs/codes"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Struct は struct のバリデートを行い、エラーコード付きのエラーを返す
func Struct(ctx context.Context, s any) error {
	err := validate.StructCtx(ctx, s)
	if err != nil {
		return errs.Wrap(codes.InvalidParameter, err)
	}
	return nil
}
