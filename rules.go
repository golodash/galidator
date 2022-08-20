package galidator

import (
	"fmt"

	filters "github.com/golodash/galidator/internal"
)

type (
	// A map with data recorded to be used in returning error messages
	option map[string]string

	// A map of different rules as key with their own option
	options map[string]option

	// A map full of validators which is assigned for a single key in a validator struct
	Validators map[string]func(interface{}) bool

	// A struct to implement rule interface
	ruleS struct {
		// Used to validate user's data
		validators Validators
		// Used in returning error messages
		options options
		// If isOptional is true, if empty is sent, all errors will be ignored
		isOptional bool
	}

	// An interface with some functions to satisfy validation purpose
	rule interface {
		// Validates all validators defined
		validate(interface{}) []string

		// Checks if input (can be)/is int
		Int() rule
		// Checks if input (can be)/is float
		Float() rule
		// Checks if input acts like: input >= min or len(input) >= min
		//
		// Note: If min > 0, field will be required
		Min(min float64) rule
		// Checks if input acts like: input <= max or len(input) <= max
		Max(max float64) rule
		// Checks if input acts like: len(input) >= from && len(input) <= to
		//
		// If from == -1, no check on from will happen
		// If to == -1, no check on to will happen
		//
		// Note: If from > 0, field will be required
		LenRange(from int, to int) rule
		// Checks if input acts like: len(input) == length
		//
		// Note: If length > 0, field will be required
		Len(length int) rule
		// Checks if input is not zero(0, "", '') or nil or empty
		//
		// Note: Field will be required
		Required() rule
		// Checks if input is not zero(0, "", '')
		//
		// Note: Field will be required
		NonZero() rule
		// Checks if input is not nil
		//
		// Note: Field will be required
		NonNil() rule
		// Checks if input has items inside it
		//
		// Note: Field will be required
		NonEmpty() rule
		// Checks if input is a valid email address
		Email() rule
		// Validates inputs with passed pattern
		Regex(pattern string) rule
		// Checks if input is a valid phone number
		Phone() rule
		// Adds custom validators
		Custom(validators Validators) rule

		// Returns option of the passed rule key
		getOption(key string) option
		// Adds a new subKey with a value associated with it to option of passed rule key
		addOption(key string, subKey string, value string)
		// Makes the field optional
		optional()
		// Makes the field required
		required()
		// Returns true if the rule has to pass all validators
		//
		// Returns false if the rule can be empty, nil or zero and is allowed to not pass any validations
		isRequired() bool
	}
)

func (o *ruleS) Int() rule {
	functionName := "int"
	o.validators[functionName] = filters.Int
	return o
}

func (o *ruleS) Float() rule {
	functionName := "float"
	o.validators[functionName] = filters.Float
	return o
}

func (o *ruleS) Min(min float64) rule {
	functionName := "min"
	o.validators[functionName] = filters.Min(min)
	precision := determinePrecision(min)
	o.addOption(functionName, "min", fmt.Sprintf("%."+precision+"f", min))
	if min > 0 {
		o.required()
	}
	return o
}

func (o *ruleS) Max(max float64) rule {
	functionName := "max"
	o.validators[functionName] = filters.Max(max)
	precision := determinePrecision(max)
	o.addOption(functionName, "max", fmt.Sprintf("%."+precision+"f", max))
	return o
}

func (o *ruleS) LenRange(from, to int) rule {
	functionName := "len_range"
	o.validators[functionName] = filters.LenRange(from, to)
	o.addOption(functionName, "from", fmt.Sprintf("%d", from))
	o.addOption(functionName, "to", fmt.Sprintf("%d", to))
	if from > 0 {
		o.required()
	}
	return o
}

func (o *ruleS) Len(length int) rule {
	functionName := "len"
	o.validators[functionName] = filters.Len(length)
	o.addOption(functionName, "length", fmt.Sprint(length))
	if length > 0 {
		o.required()
	}
	return o
}

func (o *ruleS) Required() rule {
	functionName := "required"
	o.validators[functionName] = filters.Required
	o.required()
	return o
}

func (o *ruleS) NonZero() rule {
	functionName := "non_zero"
	o.validators[functionName] = filters.NonZero
	o.required()
	return o
}

func (o *ruleS) NonNil() rule {
	functionName := "non_nil"
	o.validators[functionName] = filters.NonNil
	o.required()
	return o
}

func (o *ruleS) NonEmpty() rule {
	functionName := "non_empty"
	o.validators[functionName] = filters.NonEmpty
	o.required()
	return o
}

func (o *ruleS) Email() rule {
	functionName := "email"
	o.validators[functionName] = filters.Email
	return o
}

func (o *ruleS) Regex(pattern string) rule {
	functionName := "regex"
	o.validators[functionName] = filters.Regex(pattern)
	o.addOption(functionName, "pattern", pattern)
	if !filters.Regex(pattern)("") {
		o.required()
	}
	return o
}

func (o *ruleS) Phone() rule {
	functionName := "phone"
	o.validators[functionName] = filters.Phone
	return o
}

func (o *ruleS) Custom(validators Validators) rule {
	for key, function := range validators {
		if _, ok := o.validators[key]; ok {
			panic(fmt.Sprintf("%s is duplicate and has to be unique", key))
		}
		o.validators[key] = function
	}
	return o
}

func (o *ruleS) validate(input interface{}) []string {
	fails := []string{}
	for key, vFunction := range o.validators {
		if !vFunction(input) {
			fails = append(fails, key)
		}
	}

	return fails
}

func (o *ruleS) getOption(key string) option {
	if option, ok := o.options[key]; ok {
		return option
	}
	return option{}
}

func (o *ruleS) addOption(key string, subKey string, value string) {
	if option, ok := o.options[key]; ok {
		option[subKey] = value
		return
	}
	o.options[key] = option{subKey: value}
}

func (o *ruleS) optional() {
	o.isOptional = true
}

func (o *ruleS) required() {
	o.isOptional = false
}

func (o *ruleS) isRequired() bool {
	return !o.isOptional
}
