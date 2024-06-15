package util

import (
	"time"
)

// getTimeNowFunc is current time function.
// Default function is `time.Now`.
var getTimeNowFunc = time.Now

// GetTimeNow は現在時刻を返す。
// `getTimeNowFunc` を上書きすることで、テストなどの用途で任意の日時を取得するように調整できる。
func GetTimeNow() time.Time {
	now := getTimeNowFunc().In(time.UTC)
	return now
}
