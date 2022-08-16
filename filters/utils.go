package filters

import "reflect"

// Returns true if `input` is nil
func isNil(input interface{}) bool {
	inputValue := reflect.ValueOf(input)
	switch inputValue.Kind() {
	case reflect.Array, reflect.Struct, reflect.Map:
		return inputValue.IsNil()
	default:
		return false
	}
}
