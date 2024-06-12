package errs

import (
	"errors"
	"fmt"
	"strings"

	"yy-go-backend-template/util/errs/codes"

	"golang.org/x/xerrors"
)

// CodeError はエラーコードとスタックトレースの付いたエラー
type CodeError interface {
	error
	fmt.Formatter
	xerrors.Wrapper
	xerrors.Formatter

	// Code はエラーコードを返す
	Code() codes.Code
	// Stack はスタックトレースを返す
	Stack() string
	// レスポンス出力用のエラーメッセージを返す
	// showInternal が false の場合、内部情報を隠す
	Message(showInternal bool) string
}

type codeError struct {
	error
	code  codes.Code
	frame xerrors.Frame
}

// Code はエラーからエラーコードを返す
func Code(err error) codes.Code {
	if err == nil {
		return 0
	}

	var ce CodeError
	if ok := errors.As(err, &ce); !ok {
		return codes.Unknown
	}
	return ce.Code()
}

// Error はエラーコードからエラーを生成します
func Error(c codes.Code, text string) CodeError {
	return &codeError{
		error: xerrors.New(text),
		code:  c,
		frame: xerrors.Caller(1),
	}
}

// Errorf はエラーコードからエラーを生成します
// 最後の引数が `error` でフォーマットが ": %w" で終わる場合、エラーをラップすることができます。
func Errorf(c codes.Code, format string, v ...any) CodeError {
	return &codeError{
		error: xerrors.Errorf(format, v...),
		code:  c,
		frame: xerrors.Caller(1),
	}
}

// Wrap は任意のエラーにエラーコードに付与します
// Errorf に自分で付けるメッセージが不要な場合はこちらを使用してください。
func Wrap(c codes.Code, err error) CodeError {
	if err != nil {
		err = xerrors.Errorf("%s: %w", c.DefaultMessage(), err)
	} else {
		err = xerrors.New(c.DefaultMessage())
	}

	return &codeError{
		error: err,
		code:  c,
		frame: xerrors.Caller(1),
	}
}

func (e *codeError) Error() string {
	return e.error.Error()
}

func (e *codeError) Unwrap() error {
	return errors.Unwrap(e.error)
}

func (e *codeError) Format(s fmt.State, v rune) {
	xerrors.FormatError(e, s, v)
}

func (e *codeError) FormatError(p xerrors.Printer) (next error) {
	p.Print(e.Error())
	e.frame.Format(p)
	return e.Unwrap()
}

func (e *codeError) Code() codes.Code {
	return e.code
}

func (e *codeError) Stack() string {
	return fmt.Sprintf("%+v", e)
}

func (e *codeError) Message(showInternal bool) string {
	s := e.Error()

	if !showInternal {
		// 内部の情報をそのまま出すのは好ましくないので、コロンより後ろを切る
		if i := strings.Index(s, ": "); i > 0 {
			s = s[:i]
		}
	}
	return s
}
