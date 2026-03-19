// reflect_compat.go provides a polyfill for reflect.TypeAssert which was
// added in Go 1.26. It is used to make this package compile with Go 1.24.
package json

import "reflect"

// typeAssertVal is a polyfill for reflect.TypeAssert[T](v).
// It extracts the value held by a reflect.Value as type T.
func typeAssertVal[T any](v reflect.Value) (T, bool) {
	i, ok := v.Interface().(T)
	return i, ok
}
