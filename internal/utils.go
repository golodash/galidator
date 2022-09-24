package internal

import "reflect"

// Determines if passed variables are exactly the same
func Same(value1 interface{}, value2 interface{}) (condition bool) {
	condition = true
	v1 := reflect.ValueOf(value1)
	v2 := reflect.ValueOf(value2)

	// Check for nil and "" and other zero values
	if (!v1.IsValid() && !v2.IsValid()) && (v1.Kind() == v2.Kind()) {
		return
	}

	if v1.Kind() != v2.Kind() {
		if v1.Kind() == reflect.Interface && v2.Kind() != reflect.Interface {
			if v1.CanConvert(v2.Type()) {
				v1 = v1.Convert(v2.Type())
			}
		} else if v2.Kind() == reflect.Interface && v1.Kind() != reflect.Interface {
			if v2.CanConvert(v1.Type()) {
				v2 = v2.Convert(v1.Type())
			}
		}
		if v1.Kind() == v2.Kind() && v1.Interface() == v2.Interface() {
			return
		}
		condition = false
		return
	}

	defer func() {
		if r := recover(); r != nil {
			condition = false
		}
	}()

	switch v1.Kind() {
	case reflect.Array, reflect.Slice:
		if v1.Len() != v2.Len() {
			condition = false
			return
		}
		for i := 0; i < v1.Len(); i++ {
			condition = Same(v1.Index(i).Interface(), v2.Index(i).Interface())
			if !condition {
				condition = false
				return
			}
		}
	case reflect.Map:
		if v1.Len() != v2.Len() {
			condition = false
			return
		}

		if len(v1.MapKeys()) != len(v2.MapKeys()) {
			condition = false
			return
		}

		keys := v1.MapKeys()
		for i := 0; i < len(v1.MapKeys()); i = i + 1 {
			value1 := v1.MapIndex(keys[i])
			value2 := v2.MapIndex(keys[i])
			if !value1.IsValid() {
				condition = false
				return
			}
			if !value2.IsValid() {
				condition = false
				return
			}

			if condition = Same(value1.Interface(), value2.Interface()); !condition {
				condition = false
				return
			}
		}
	case reflect.Struct:
		condition = value1 == value2
	case reflect.Ptr:
		condition = Same(v1.Elem().Interface(), v2.Elem().Interface())
	default:
		if v1.Interface() != v2.Interface() {
			condition = false
		}
	}

	return
}
