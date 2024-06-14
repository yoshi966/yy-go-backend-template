package str

import (
	"unicode"
	"unicode/utf8"
)

// UpperFirst support ASCII only
func UpperFirst(s string) string {
	if len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if r != utf8.RuneError {
			upper := unicode.ToUpper(r)
			if upper != r {
				s = string(upper) + s[size:]
			}
		}
	}
	return s
}
