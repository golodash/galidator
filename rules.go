package galidator

import (
	"fmt"

	rules "github.com/golodash/galidator/internal"
)

type (
	// A map with data recorded to be used in returning error messages
	option map[string]string

	// A map of different rules as key with their own option
	options map[string]option

	// A map full of validators which is assigned for a single key in a validator struct
	Validators map[string]func(interface{}) bool

	// A struct to implement ruleSet interface
	ruleSetS struct {
		// Used to validate user's data
		validators Validators
		// Used in returning error messages
		options options
		// If isOptional is true, if empty is sent, all errors will be ignored
		isOptional bool
		// Holds data for more complex structures, like:
		//
		// map, slice or struct
		deepValidator validator
		// Defines type of elements of a slice
		childrenRule ruleSet
	}

	// An interface with some functions to satisfy validation purpose
	ruleSet interface {
		// Validates all validators defined
		validate(interface{}) []string

		// Checks if input (can be)/is int
		Int() ruleSet
		// Checks if input (can be)/is float
		Float() ruleSet
		// Checks if input acts like: input >= min or len(input) >= min
		//
		// Note: If min > 0, field will be required
		Min(min float64) ruleSet
		// Checks if input acts like: input <= max or len(input) <= max
		Max(max float64) ruleSet
		// Checks if input acts like: len(input) >= from && len(input) <= to
		//
		// If from == -1, no check on from will happen
		// If to == -1, no check on to will happen
		//
		// Note: If from > 0, field will be required
		LenRange(from int, to int) ruleSet
		// Checks if input acts like: len(input) == length
		//
		// Note: If length > 0, field will be required
		Len(length int) ruleSet
		// Checks if input is not zero(0, "", '') or nil or empty
		//
		// Note: Field will be required
		Required() ruleSet
		// Checks if input is not zero(0, "", '')
		//
		// Note: Field will be required
		NonZero() ruleSet
		// Checks if input is not nil
		//
		// Note: Field will be required
		NonNil() ruleSet
		// Checks if input has items inside it
		//
		// Note: Field will be required
		NonEmpty() ruleSet
		// Checks if input is a valid email address
		Email() ruleSet
		// Validates inputs with passed pattern
		//
		// Note: If pattern does not pass empty string, it will be required
		Regex(pattern string) ruleSet
		// Checks if input is a valid phone number
		Phone() ruleSet
		// Adds custom validators
		Custom(validators Validators) ruleSet
		// Checks if input is a map
		Map() ruleSet
		// Checks if input is a slice
		Slice() ruleSet
		// Checks if input is a struct
		Struct() ruleSet
		// Adds another deeper layer to validation structure
		//
		// Can check struct and map
		Complex(validator validator) ruleSet
		// If children of a slice is not struct or map, use this function and otherwise use Complex function after Slice function
		Children(ruleSet ruleSet) ruleSet
		// Checks if input is at least 8 characters long, has one lowercase, one uppercase and one number character
		//
		// Note: Field will be required
		Password() ruleSet

		// Returns option of the passed ruleKey
		getOption(ruleKey string) option
		// Adds a new subKey with a value associated with it to option of passed ruleKey
		addOption(ruleKey string, subKey string, value string)
		// Makes the field optional
		optional()
		// Makes the field required
		required()
		// Returns true if the ruleSet has to pass all validators
		//
		// Returns false if the ruleSet can be empty, nil or zero and is allowed to not pass any validations
		isRequired() bool
		// Returns true if deepValidator is not nil
		hasDeepValidator() bool
		// Validates deepValidator
		validateDeepValidator(input interface{}) interface{}
		// Returns true if children is not nil
		hasChildrenRule() bool
		// Returns children ruleSet
		getChildrenRule() ruleSet
	}
)

func (o *ruleSetS) Int() ruleSet {
	functionName := "int"
	o.validators[functionName] = rules.Int
	return o
}

func (o *ruleSetS) Float() ruleSet {
	functionName := "float"
	o.validators[functionName] = rules.Float
	return o
}

func (o *ruleSetS) Min(min float64) ruleSet {
	functionName := "min"
	o.validators[functionName] = rules.Min(min)
	precision := determinePrecision(min)
	o.addOption(functionName, "min", fmt.Sprintf("%."+precision+"f", min))
	if min > 0 {
		o.required()
	}
	return o
}

func (o *ruleSetS) Max(max float64) ruleSet {
	functionName := "max"
	o.validators[functionName] = rules.Max(max)
	precision := determinePrecision(max)
	o.addOption(functionName, "max", fmt.Sprintf("%."+precision+"f", max))
	return o
}

func (o *ruleSetS) LenRange(from, to int) ruleSet {
	functionName := "len_range"
	o.validators[functionName] = rules.LenRange(from, to)
	o.addOption(functionName, "from", fmt.Sprintf("%d", from))
	o.addOption(functionName, "to", fmt.Sprintf("%d", to))
	if from > 0 {
		o.required()
	}
	return o
}

func (o *ruleSetS) Len(length int) ruleSet {
	functionName := "len"
	o.validators[functionName] = rules.Len(length)
	o.addOption(functionName, "length", fmt.Sprint(length))
	if length > 0 {
		o.required()
	}
	return o
}

func (o *ruleSetS) Required() ruleSet {
	functionName := "required"
	o.validators[functionName] = rules.Required
	o.required()
	return o
}

func (o *ruleSetS) NonZero() ruleSet {
	functionName := "non_zero"
	o.validators[functionName] = rules.NonZero
	o.required()
	return o
}

func (o *ruleSetS) NonNil() ruleSet {
	functionName := "non_nil"
	o.validators[functionName] = rules.NonNil
	o.required()
	return o
}

func (o *ruleSetS) NonEmpty() ruleSet {
	functionName := "non_empty"
	o.validators[functionName] = rules.NonEmpty
	o.required()
	return o
}

func (o *ruleSetS) Email() ruleSet {
	functionName := "email"
	o.validators[functionName] = rules.Email
	return o
}

func (o *ruleSetS) Regex(pattern string) ruleSet {
	functionName := "regex"
	o.validators[functionName] = rules.Regex(pattern)
	o.addOption(functionName, "pattern", pattern)
	if !rules.Regex(pattern)("") {
		o.required()
	}
	return o
}

func (o *ruleSetS) Phone() ruleSet {
	functionName := "phone"
	o.validators[functionName] = rules.Phone
	return o
}

func (o *ruleSetS) Custom(validators Validators) ruleSet {
	for key, function := range validators {
		if _, ok := o.validators[key]; ok {
			panic(fmt.Sprintf("%s is duplicate and has to be unique", key))
		}
		o.validators[key] = function
	}
	return o
}

func (o *ruleSetS) Map() ruleSet {
	functionName := "map"
	o.validators[functionName] = rules.Map
	return o
}

func (o *ruleSetS) Slice() ruleSet {
	functionName := "slice"
	o.validators[functionName] = rules.Slice
	return o
}

func (o *ruleSetS) Struct() ruleSet {
	functionName := "struct"
	o.validators[functionName] = rules.Struct
	return o
}

func (o *ruleSetS) Complex(validator validator) ruleSet {
	o.deepValidator = validator
	return o
}

func (o *ruleSetS) Children(ruleSet ruleSet) ruleSet {
	o.childrenRule = ruleSet
	return o
}

func (o *ruleSetS) Password() ruleSet {
	functionName := "password"
	o.validators[functionName] = rules.Password
	o.required()
	return o
}

func (o *ruleSetS) hasChildrenRule() bool {
	return o.childrenRule != nil
}

func (o *ruleSetS) getChildrenRule() ruleSet {
	return o.childrenRule
}

func (o *ruleSetS) validate(input interface{}) []string {
	fails := []string{}
	for key, vFunction := range o.validators {
		if !vFunction(input) {
			fails = append(fails, key)
		}
	}

	return fails
}

func (o *ruleSetS) getOption(ruleKey string) option {
	if option, ok := o.options[ruleKey]; ok {
		return option
	}
	return option{}
}

func (o *ruleSetS) addOption(ruleKey string, subKey string, value string) {
	if option, ok := o.options[ruleKey]; ok {
		option[subKey] = value
		return
	}
	o.options[ruleKey] = option{subKey: value}
}

func (o *ruleSetS) optional() {
	o.isOptional = true
}

func (o *ruleSetS) required() {
	o.isOptional = false
}

func (o *ruleSetS) isRequired() bool {
	return !o.isOptional
}

func (o *ruleSetS) hasDeepValidator() bool {
	return o.deepValidator != nil
}

func (o *ruleSetS) validateDeepValidator(input interface{}) interface{} {
	return o.deepValidator.Validate(input)
}
