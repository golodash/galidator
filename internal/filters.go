package filters

import (
	"net/mail"
	"reflect"
	"regexp"
	"strconv"
)

// A map which with rule's key will provide the default error message of that key
var DefaultValidatorErrorMessages = map[string]string{
	"int":       "$field is not integer",
	"float":     "$field is not float",
	"min":       "$field must be higher equal to $min",
	"max":       "$field must be lower equal to $max",
	"len_range": "$field's length must be between $from to $to characters long",
	"len":       "$field's length must be equal to $length",
	"required":  "$field can not be 0 or \"\" or '' or nil or empty",
	"non_zero":  "$field can not be 0 or \"\" or ''",
	"non_nil":   "$field can not be nil",
	"non_empty": "$field can not be empty",
	"email":     "$field is not a valid email address",
	"regex":     "$field does not pass /$pattern/ pattern",
	"phone":     "$field is not a valid phone number",
	"map":       "$field is not a map",
	"struct":    "$field is not a struct",
	"slice":     "$field is not a slice",
}

// Returns true if input (can be)/is int
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

// Returns true if input (can be)/is float
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

// Returns true if: input >= min or len(input) >= min
func Min(min float64) func(interface{}) bool {
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
func Max(max float64) func(interface{}) bool {
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
func LenRange(from, to int) func(interface{}) bool {
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
func Len(length int) func(interface{}) bool {
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
func Required(input interface{}) bool {
	inputValue := reflect.ValueOf(input)
	return !inputValue.IsZero() && !isNil(input) && !hasZeroItems(input)
}

// Not a filter
//
// A helper function to use inside code
func IsEmptyNilZero(input interface{}) bool {
	return !Required(input)
}

// Returns true if input is not zero(0, "", ”)
func NonZero(input interface{}) bool {
	return !reflect.ValueOf(input).IsZero()
}

// Returns true if input is not nil
func NonNil(input interface{}) bool {
	return !isNil(input)
}

// Returns true if input has items
func NonEmpty(input interface{}) bool {
	return !hasZeroItems(input)
}

// Returns true if input is a valid email
func Email(input interface{}) bool {
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
func Phone(input interface{}) bool {
	return Regex(`((\+|\(|0)?\d{1,3})?((\s|\)|\-))?(\d{10})$`)(input)
}

// Returns true if input matches the passed pattern
func Regex(pattern string) func(interface{}) bool {
	regex := regexp.MustCompile(pattern)
	return func(input interface{}) bool {
		inputValue := reflect.ValueOf(input)
		switch inputValue.Kind() {
		case reflect.String:
			return regex.MatchString(input.(string))
		default:
			return false
		}
	}
}

// Returns true if input is a map
func Map(input interface{}) bool {
	return reflect.TypeOf(input).Kind() == reflect.Map
}

// Returns true if input is a struct
func Struct(input interface{}) bool {
	return reflect.TypeOf(input).Kind() == reflect.Struct
}

// Returns true if input is a slice
func Slice(input interface{}) bool {
	return reflect.TypeOf(input).Kind() == reflect.Slice
}
