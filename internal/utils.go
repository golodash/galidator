package filters

import (
	"reflect"
)

// Returns true if input is nil
func isNil(input interface{}) bool {
	inputValue := reflect.ValueOf(input)
	switch inputValue.Kind() {
	case reflect.Array, reflect.Slice, reflect.Struct, reflect.Map:
		return inputValue.IsNil()
	default:
		return false
	}
}

func hasZeroItems(input interface{}) bool {
	inputValue := reflect.ValueOf(input)
	switch inputValue.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		return inputValue.Len() == 0
	default:
		return false
	}
}
