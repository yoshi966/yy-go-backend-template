package errs

import (
	"runtime"
	"strings"

	"github.com/Yoshioka9709/yy-go-backend-template/util/errs/codes"

	"golang.org/x/xerrors"
)

// NewPanicError は recover の戻り値からエラーを生成します
func NewPanicError(e any) CodeError {
	// panic が発生した箇所を無理やり割り出す
	skip := 1
	rt := false
	for ; ; skip++ {
		_, file, _, ok := runtime.Caller(skip)
		if !ok {
			// 見つからなければ呼び出し元に設定
			skip = 1
			break
		}

		// go の runtime パッケージが最後に出てきた場所を panic 発生位置とみなす
		if strings.Contains(file, "/src/runtime/") {
			rt = true
		} else if rt {
			break
		}
	}

	return &codeError{
		error: xerrors.Errorf("panic: %v", e),
		code:  codes.InternalError,
		frame: xerrors.Caller(skip),
	}
}
