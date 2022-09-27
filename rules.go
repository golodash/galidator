package galidator

import (
	"net/mail"
	"reflect"

	"github.com/dlclark/regexp2"
	"github.com/golodash/galidator/internal"
)

// A map which with rule or require's key will provide the default error message of that rule or require's key
var defaultValidatorErrorMessages = map[string]string{
	// Rules
	"int":       "not an integer value",
	"float":     "not a float value",
	"min":       "$field's length must be higher equal to $min",
	"max":       "$field's length must be lower equal to $max",
	"len_range": "$field's length must be between $from to $to characters long",
	"len":       "$field's length must be equal to $length",
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
	"password":  "$field must be at least 8 characters long and contain one lowercase, one uppercase, one special and one number character",
	"or":        "ruleSets in $field did not pass based on or logic",
	"xor":       "ruleSets in $field did not pass based on xor logic",
	"choices":   "$value does not include in allowed choices: $choices",
	"string":    "not a string",
	"type":      "not a $type",

	// Requires
	"when_exist_one": "$field is required because at least one of $choices fields are not nil, empty or zero(0, \"\", '')",
	"when_exist_all": "$field is required because all of $choices fields are not nil, empty or zero(0, \"\", '')",
}

// Returns true if input is int
func intRule(input interface{}) bool {
	inputValue := reflect.ValueOf(input)
	switch inputValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

// Returns true if input is float
func floatRule(input interface{}) bool {
	inputValue := reflect.ValueOf(input)
	switch inputValue.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// Returns true if input type is equal to defined type
func typeRule(t string) func(interface{}) bool {
	return func(input interface{}) bool {
		return reflect.TypeOf(input).String() == t
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
	return !reflect.ValueOf(input).IsValid() || !reflect.ValueOf(input).IsZero()
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

// If at least one of the passed ruleSets pass, this rule will pass
func orRule(ruleSets ...ruleSet) func(interface{}) bool {
	return func(input interface{}) bool {
		output := false

		if len(ruleSets) == 0 {
			return true
		}

		for _, ruleSet := range ruleSets {
			output = len(ruleSet.validate(input)) == 0 || output
		}

		return output
	}
}

// If Xor of the passed ruleSets pass, this rule will pass
func xorRule(ruleSets ...ruleSet) func(interface{}) bool {
	return func(input interface{}) bool {
		output := false

		if len(ruleSets) == 0 {
			return true
		}

		for _, ruleSet := range ruleSets {
			output = len(ruleSet.validate(input)) == 0 != output
		}

		return output
	}
}

// If passed data is not one of choices in choices variable, it fails
func choicesRule(choices ...interface{}) func(interface{}) bool {
	return func(input interface{}) bool {
		for i := 0; i < len(choices); i++ {
			element := choices[i]
			if internal.Same(element, input) {
				return true
			}
		}

		return false
	}
}

// Returns true if input is string
func stringRule(input interface{}) bool {
	return reflect.TypeOf(input).Kind() == reflect.String
}
