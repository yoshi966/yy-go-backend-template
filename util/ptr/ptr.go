package ptr

// True はtrue(bool型)のポインタを返します。
func True() *bool {
	return Ptr(true)
}

// False はfalse(bool型)のポインタを返します。
func False() *bool {
	return Ptr(false)
}

// Ptr は任意の値を受け取り、そのポインタを返します。
func Ptr[T any](v T) *T {
	return &v
}
