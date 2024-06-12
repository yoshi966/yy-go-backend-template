package codes

import "net/http"

// Code はエラーコードです
type Code uint32

const (
	// InvalidParameter はパラメーターの不足や不正な値が渡された場合のエラー
	InvalidParameter Code = iota + 1

	// ResourceNotFound は条件に一致するデータが見つからない
	ResourceNotFound

	// Conflict は矛盾する操作を行おうとした
	Conflict

	// InvalidToken は認証トークンが無効
	InvalidToken

	// Unauthorized はトークン由来以外の認証エラー
	Unauthorized

	// Forbidden は権限不足エラー
	Forbidden

	// InternalError はサーバー内部エラー
	//
	// 他のエラーに当てはまらないもの（DBエラーなど）全てが対象
	InternalError

	// ServiceUnavailable は一時的なサービス利用不可
	ServiceUnavailable

	// BadRequest は汎用的なクライアントエラー
	//
	// 他のエラーに当てはまらないクライアント起因エラーが対象
	BadRequest

	// Unknown は不明なエラーコード
	Unknown Code = 0xffffffff
)

var codeText = map[Code]string{
	InvalidParameter:   "InvalidParameter",
	ResourceNotFound:   "ResourceNotFound",
	Conflict:           "Conflict",
	InvalidToken:       "InvalidToken",
	Unauthorized:       "Unauthorized",
	Forbidden:          "Forbidden",
	InternalError:      "InternalError",
	ServiceUnavailable: "ServiceUnavailable",
	BadRequest:         "BadRequest",
	Unknown:            "Unknown",
}

var codeDefaultMessage = map[Code]string{
	InvalidParameter:   "invalid parameter",
	ResourceNotFound:   "resource not found",
	Conflict:           "conflict",
	InvalidToken:       "invalid token",
	Unauthorized:       "unauthorized",
	Forbidden:          "access denied",
	InternalError:      "internal error",
	ServiceUnavailable: "service unavailable",
	BadRequest:         "bad request",
	Unknown:            "unknown",
}

var codeError = map[Code]string{
	InvalidParameter:   "invalid parameter",
	ResourceNotFound:   "resource not found",
	Conflict:           "conflict",
	InvalidToken:       "invalid token",
	Unauthorized:       "unauthorized",
	Forbidden:          "access denied",
	InternalError:      "internal error",
	ServiceUnavailable: "service unavailable",
	BadRequest:         "bad request",
	Unknown:            "unknown",
}

var codeHTTPStatus = map[Code]int{
	InvalidParameter:   http.StatusBadRequest,
	ResourceNotFound:   http.StatusBadRequest,
	Conflict:           http.StatusConflict,
	InvalidToken:       http.StatusUnauthorized,
	Unauthorized:       http.StatusUnauthorized,
	Forbidden:          http.StatusForbidden,
	InternalError:      http.StatusInternalServerError,
	ServiceUnavailable: http.StatusServiceUnavailable,
	BadRequest:         http.StatusBadRequest,
	Unknown:            http.StatusInternalServerError,
}

// String 文字列で取得
func (c Code) String() string {
	if text, ok := codeText[c]; ok {
		return text
	}
	return codeText[Unknown]
}

// Error はエラーコードの説明文を返す
func (c Code) Error() string {
	if text, ok := codeError[c]; ok {
		return text
	}
	return codeError[Unknown]
}

// DefaultMessage はエラーコードのデフォルト説明文を返す
func (c Code) DefaultMessage() string {
	if text, ok := codeDefaultMessage[c]; ok {
		return text
	}
	return codeDefaultMessage[Unknown]
}

// HTTPStatus はエラーコードが示すHTTPステータスコードを返す
func (c Code) HTTPStatus() int {
	if status, ok := codeHTTPStatus[c]; ok {
		return status
	}
	return codeHTTPStatus[Unknown]
}
