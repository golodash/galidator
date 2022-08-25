package galidator

import (
	"net/mail"
	"reflect"
	"strconv"

	"github.com/dlclark/regexp2"
)

// A map which with rule's key will provide the default error message of that rule's key
var defaultValidatorErrorMessages = map[string]string{
	"int":       "not an integer value",
	"float":     "not a float value",
	"min":       "$fieldS's length must be higher equal to $min",
	"max":       "$fieldS's length must be lower equal to $max",
	"len_range": "$fieldS's length must be between $from to $to characters long",
	"len":       "$fieldS's length must be equal to $length",
	"required":  "required",
	"non_zero":  "can not be 0",
	"non_nil":   "can not be nil",
	"non_empty": "can not be empty",
	"email":     "not a valid email address",
	"regex":     "$value does not pass /$pattern/ pattern",
	"phone":     "$value is not a valid phone number",
	"map":       "not a map",
	"struct":    "not a struct",
	"slice":     "not a slice",
	"password":  "$fieldS must be at least 8 characters long and contain one lowercase, one uppercase, one special and one number character",
}

// Returns true if input (can be)/is int
func intRule(input interface{}) bool {
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

// Returns true if input (can be)/is float
func floatRule(input interface{}) bool {
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

// Returns true if: input >= min or len(input) >= min
func minRule(min float64) func(interface{}) bool {
	return func(input interface{}) bool {
		inputValue := reflect.ValueOf(input)
		switch inputValue.Kind() {
		case reflect.String, reflect.Map, reflect.Slice:
			return inputValue.Len() >= int(min)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			return inputValue.Convert(reflect.TypeOf(1.0)).Float() >= min
		default:
			return false
		}
	}
}

// Returns true if: input <= max or len(input) <= max
func maxRule(max float64) func(interface{}) bool {
	return func(input interface{}) bool {
		inputValue := reflect.ValueOf(input)
		switch inputValue.Kind() {
		case reflect.String, reflect.Map, reflect.Slice:
			return inputValue.Len() <= int(max)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
			reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			return inputValue.Convert(reflect.TypeOf(1.0)).Float() <= max
		default:
			return false
		}
	}
}

// Returns true if len(input) >= from && len(input) <= to
//
// If from == -1, no check on from config will happen
// If to == -1, no check on to config will happen
func lenRangeRule(from, to int) func(interface{}) bool {
	return func(input interface{}) bool {
		inputValue := reflect.ValueOf(input)
		switch inputValue.Kind() {
		case reflect.String, reflect.Map, reflect.Slice:
			if from != -1 && inputValue.Len() < from {
				return false
			} else if to != -1 && inputValue.Len() > to {
				return false
			}
			return true
		default:
			return false
		}
	}
}

// Returns true if len(input) is equal to passed length
func lenRule(length int) func(interface{}) bool {
	return func(input interface{}) bool {
		inputValue := reflect.ValueOf(input)
		switch inputValue.Kind() {
		case reflect.String, reflect.Map, reflect.Slice:
			return inputValue.Len() == length
		default:
			return false
		}
	}
}

// Returns true if input is not 0, "", ”, nil and empty
func requiredRule(input interface{}) bool {
	inputValue := reflect.ValueOf(input)
	return !inputValue.IsZero() && !isNil(input) && !hasZeroItems(input)
}

// Returns true if input is not zero(0, "", ”)
func nonZeroRule(input interface{}) bool {
	return !reflect.ValueOf(input).IsZero()
}

// Returns true if input is not nil
func nonNilRule(input interface{}) bool {
	return !isNil(input)
}

// Returns true if input has items
func nonEmptyRule(input interface{}) bool {
	return !hasZeroItems(input)
}

// Returns true if input is a valid email
func emailRule(input interface{}) bool {
	switch reflect.ValueOf(input).Kind() {
	case reflect.String:
		_, err := mail.ParseAddress(input.(string))
		if err != nil {
			return false
		}
		return true
	default:
		return false
	}
}

// Returns true if input is a valid phone number
func phoneRule(input interface{}) bool {
	return regexRule(`((\+|\(|0)?\d{1,3})?((\s|\)|\-))?(\d{10})$`)(input)
}

// Returns true if input matches the passed pattern
func regexRule(pattern string) func(interface{}) bool {
	regex := regexp2.MustCompile(pattern, regexp2.None)
	return func(input interface{}) bool {
		inputValue := reflect.ValueOf(input)
		switch inputValue.Kind() {
		case reflect.String:
			output, _ := regex.MatchString(input.(string))
			return output
		default:
			return false
		}
	}
}

// Returns true if input is a map
func mapRule(input interface{}) bool {
	return reflect.TypeOf(input).Kind() == reflect.Map
}

// Returns true if input is a struct
func structRule(input interface{}) bool {
	return reflect.TypeOf(input).Kind() == reflect.Struct
}

// Returns true if input is a slice
func sliceRule(input interface{}) bool {
	return reflect.TypeOf(input).Kind() == reflect.Slice
}

// Returns true if input is at least 8 characters long, has one lowercase, one uppercase, one special and one number character
func passwordRule(input interface{}) bool {
	return regexRule(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`)(input)
}

// func whenExistRule(validator validator, options option) func(interface{}) bool {
// 	return func(input interface{}) bool {
// 		validator.GetRules()validator
// 	}
// }
