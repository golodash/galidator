package galidator

import (
	"context"
	"net/mail"
	"reflect"

	"github.com/dlclark/regexp2"
	"github.com/golodash/godash/generals"
	"github.com/nyaruka/phonenumbers"
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
	"phone":     "$value is not a valid international phone number format",
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
	"when_exist_one":     "$field is required because at least one of $choices fields are not nil, empty or zero(0, \"\", '')",
	"when_exist_all":     "$field is required because all of $choices fields are not nil, empty or zero(0, \"\", '')",
	"when_not_exist_one": "$field is required because at least one of $choices fields are nil, empty or zero(0, \"\", '')",
	"when_not_exist_all": "$field is required because all of $choices fields are nil, empty or zero(0, \"\", '')",
}

func isValid(input interface{}) bool {
	return reflect.ValueOf(input).IsValid()
}

// Returns true if input is int
func intRule(ctx context.Context, input interface{}) bool {
	if !isValid(input) {
		return false
	}
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
func floatRule(ctx context.Context, input interface{}) bool {
	if !isValid(input) {
		return false
	}
	inputValue := reflect.ValueOf(input)
	switch inputValue.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// Returns true if input type is equal to defined type
func typeRule(t string) func(context.Context, interface{}) bool {
	return func(ctx context.Context, input interface{}) bool {
		return isValid(input) && reflect.TypeOf(input).String() == t
	}
}

// Returns true if: input >= min or len(input) >= min
func minRule(min float64) func(context.Context, interface{}) bool {
	return func(ctx context.Context, input interface{}) bool {
		if !isValid(input) {
			return false
		}
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
func maxRule(max float64) func(context.Context, interface{}) bool {
	return func(ctx context.Context, input interface{}) bool {
		if !isValid(input) {
			return false
		}
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
func lenRangeRule(from, to int) func(context.Context, interface{}) bool {
	return func(ctx context.Context, input interface{}) bool {
		if !isValid(input) {
			return false
		}
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
func lenRule(length int) func(context.Context, interface{}) bool {
	return func(ctx context.Context, input interface{}) bool {
		if !isValid(input) {
			return false
		}
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
func requiredRule(ctx context.Context, input interface{}) bool {
	return !isNil(input) && !reflect.ValueOf(input).IsZero() && !hasZeroItems(input)
}

// Returns true if input is not zero(0, "", ”)
func nonZeroRule(ctx context.Context, input interface{}) bool {
	return !isValid(input) || !reflect.ValueOf(input).IsZero()
}

// Returns true if input is not nil
func nonNilRule(ctx context.Context, input interface{}) bool {
	return !isNil(input)
}

// Returns true if input has items
func nonEmptyRule(ctx context.Context, input interface{}) bool {
	return !hasZeroItems(input)
}

// Returns true if input is a valid email
func emailRule(ctx context.Context, input interface{}) bool {
	if !isValid(input) {
		return false
	}
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
func phoneRule(ctx context.Context, input interface{}) bool {
	if !isValid(input) {
		return false
	}
	InternationalPhoneRegex := regexp2.MustCompile(`^\+\d+$`, regexp2.None)
	if ok, err := InternationalPhoneRegex.MatchString(input.(string)); !ok || err != nil {
		return false
	}
	parsedNumber, err := phonenumbers.Parse(input.(string), "")
	return err == nil && phonenumbers.IsValidNumber(parsedNumber)
}

// Returns true if input matches the passed pattern
func regexRule(pattern string) func(context.Context, interface{}) bool {
	regex := regexp2.MustCompile(pattern, regexp2.None)
	return func(ctx context.Context, input interface{}) bool {
		if !isValid(input) {
			return false
		}
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
func mapRule(ctx context.Context, input interface{}) bool {
	if !isValid(input) {
		return false
	}
	return reflect.TypeOf(input).Kind() == reflect.Map
}

// Returns true if input is a struct
func structRule(ctx context.Context, input interface{}) bool {
	if !isValid(input) {
		return false
	}
	return reflect.TypeOf(input).Kind() == reflect.Struct
}

// Returns true if input is a slice
func sliceRule(ctx context.Context, input interface{}) bool {
	if !isValid(input) {
		return false
	}
	return reflect.TypeOf(input).Kind() == reflect.Slice
}

// Returns true if input is at least 8 characters long, has one lowercase, one uppercase, one special and one number character
func passwordRule(ctx context.Context, input interface{}) bool {
	if !isValid(input) {
		return false
	}
	return regexRule(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`)(ctx, input)
}

// If at least one of the passed ruleSets pass, this rule will pass
func orRule(ruleSets ...ruleSet) func(context.Context, interface{}) bool {
	return func(ctx context.Context, input interface{}) bool {
		output := false

		if len(ruleSets) == 0 {
			return true
		}

		for _, ruleSet := range ruleSets {
			output = len(ruleSet.validate(ctx, input)) == 0 || output
		}

		return output
	}
}

// If Xor of the passed ruleSets pass, this rule will pass
func xorRule(ruleSets ...ruleSet) func(context.Context, interface{}) bool {
	return func(ctx context.Context, input interface{}) bool {
		output := false

		if len(ruleSets) == 0 {
			return true
		}

		for _, ruleSet := range ruleSets {
			output = len(ruleSet.validate(ctx, input)) == 0 != output
		}

		return output
	}
}

// If passed data is not one of choices in choices variable, it fails
func choicesRule(choices ...interface{}) func(context.Context, interface{}) bool {
	return func(ctx context.Context, input interface{}) bool {
		for i := 0; i < len(choices); i++ {
			element := choices[i]
			if generals.Same(element, input) {
				return true
			}
		}

		return false
	}
}

// Returns true if input is string
func stringRule(ctx context.Context, input interface{}) bool {
	if !isValid(input) {
		return false
	}
	return reflect.TypeOf(input).Kind() == reflect.String
}
