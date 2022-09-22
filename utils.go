package galidator

import (
	"fmt"
	"math"
	"reflect"
)

// Determines the precision of a float number for print
func determinePrecision(number float64) string {
	for i := 0; ; i++ {
		ten := math.Pow10(i)
		if math.Floor(ten*number) == ten*number {
			return fmt.Sprint(i)
		}
	}
}

// Returns true if input is nil
func isNil(input interface{}) bool {
	inputValue := reflect.ValueOf(input)
	switch inputValue.Kind() {
	case reflect.Slice, reflect.Map:
		return inputValue.IsNil()
	default:
		return false
	}
}

// Returns true if input is map or slice and has 0 elements
//
// Returns false if input is not map or slice
func hasZeroItems(input interface{}) bool {
	inputValue := reflect.ValueOf(input)
	switch inputValue.Kind() {
	case reflect.Slice, reflect.Map:
		return inputValue.Len() == 0
	default:
		return false
	}
}

// A helper function to use inside code
func isEmptyNilZero(input interface{}) bool {
	return !requiredRule(input)
}

// Determines if passed variables are exactly the same
func same(value1 interface{}, value2 interface{}) (condition bool) {
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
			condition = same(v1.Index(i).Interface(), v2.Index(i).Interface())
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

			if condition = same(value1.Interface(), value2.Interface()); !condition {
				condition = false
				return
			}
		}
	case reflect.Struct:
		condition = value1 == value2
	case reflect.Ptr:
		condition = same(v1.Elem().Interface(), v2.Elem().Interface())
	default:
		if v1.Interface() != v2.Interface() {
			condition = false
		}
	}

	return
}

// Returns values of passed fields from passed struct or map
func getValues(all interface{}, fields ...string) []interface{} {
	fieldsValues := []interface{}{}
	allValue := reflect.ValueOf(all)

	if allValue.Kind() == reflect.Map {
		for _, key := range fields {
			element := allValue.MapIndex(reflect.ValueOf(key))
			if !element.IsValid() {
				element = allValue.MapIndex(reflect.ValueOf(key))
			}

			if !element.IsValid() {
				panic(fmt.Sprintf("value on %s is not valid", key))
			}

			fieldsValues = append(fieldsValues, element.Interface())
		}
	} else if allValue.Kind() == reflect.Struct {
		for _, key := range fields {
			element := allValue.FieldByName(key)
			if !element.IsValid() {
				panic(fmt.Sprintf("value on %s is not valid", key))
			}

			fieldsValues = append(fieldsValues, element.Interface())
		}
	}

	return fieldsValues
}

// Returns a list of keys for requires which determine not required and a bool which determines if we need to validate or not
func determineRequires(all interface{}, input interface{}, requires requires) (map[string]interface{}, bool) {
	output := map[string]interface{}{}
	for key, req := range requires {
		if !req(all)(input) {
			output[key] = 1
		}
	}

	return output, len(output) == len(requires)
}

func addTypeCheck(r ruleSet, kind reflect.Kind) {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		r.Int()
	case reflect.Float32, reflect.Float64:
		r.Float()
	case reflect.Slice:
		r.Slice()
	case reflect.String:
		r.String()
	}
}
