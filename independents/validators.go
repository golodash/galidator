package independents

import (
	"reflect"
	"strconv"
)

var DefaultValidatorErrorMessages = map[string]string{
	"int":      "{field} is not integer",
	"float":    "{field} is not float",
	"min":      "{field} must be higher equal to {min}",
	"max":      "{field} must be lower equal to {max}",
	"len":      "{field}'s length must be between {from} to {to} characters long",
	"required": "{field} can not be empty",
}

func Int(input interface{}) bool {
	inputValue := reflect.ValueOf(input)
	switch inputValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		return true
	case reflect.String:
		_, err := strconv.Atoi(input.(string))
		if err != nil {
			return false
		}
		return true
	default:
		return false
	}
}

func Float(input interface{}) bool {
	inputValue := reflect.ValueOf(input)
	switch inputValue.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		_, err := strconv.ParseFloat(input.(string), 64)
		if err != nil {
			return false
		}
		return true
	default:
		return false
	}
}

func Min(min float64) func(interface{}) bool {
	return func(input interface{}) bool {
		inputValue := reflect.ValueOf(input)
		switch inputValue.Kind() {
		case reflect.String:
			inputFloat, err := strconv.ParseFloat(input.(string), 64)
			if err != nil {
				return false
			}
			return inputFloat >= min
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			return inputValue.Convert(reflect.TypeOf(1.0)).Float() >= min
		default:
			return false
		}
	}
}

func Max(max float64) func(interface{}) bool {
	return func(input interface{}) bool {
		inputValue := reflect.ValueOf(input)
		switch inputValue.Kind() {
		case reflect.String:
			inputFloat, err := strconv.ParseFloat(input.(string), 64)
			if err != nil {
				return false
			}
			return inputFloat <= max
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			return inputValue.Convert(reflect.TypeOf(1.0)).Float() <= max
		default:
			return false
		}
	}
}

func Len(from, to int) func(interface{}) bool {
	return func(input interface{}) bool {
		inputValue := reflect.ValueOf(input)
		if inputValue.Kind() == reflect.String {
			inputStr := input.(string)
			if from != -1 && len(inputStr) < from {
				return false
			} else if to != -1 && len(inputStr) > to {
				return false
			}
		} else {
			return false
		}
		return true
	}
}

func Required(input interface{}) bool {
	inputValue := reflect.ValueOf(input)
	return !inputValue.IsZero() && !isNil(input)
}
