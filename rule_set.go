package galidator

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	// A map with data recorded to be used in returning error messages
	option map[string]string

	// A map of different rules as key with their own option
	options map[string]option

	// A map full of validators which is assigned for a single key in a validator struct
	Validators map[string]func(interface{}) bool

	// A map full of field require determining
	requires map[string]func(interface{}) func(interface{}) bool

	// A struct to implement ruleSet interface
	ruleSetS struct {
		// The name that will be shown in output of the errors
		name string
		// Used to validate user's data
		validators Validators
		// Used to determine what needs to be required
		requires requires
		// Used in returning error messages
		options options
		// Sets messages for specific rules in current ruleSet
		specificMessages Messages
		// If isOptional is true, if empty is sent, all errors will be ignored
		isOptional bool
		// Holds data for more complex structures, like:
		//
		// map, slice or struct
		deepValidator validator
		// Defines type of elements of a slice
		childrenValidator validator
	}

	// An interface with some functions to satisfy validation purpose
	ruleSet interface {
		// Validates all validators defined
		validate(interface{}) []string

		// Checks if input is int
		Int() ruleSet
		// Checks if input is float
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
		// Makes this field required, Checks if input is not zero(0, "", ''), nil or empty
		//
		// Note: Field will be required
		Required() ruleSet
		// Makes this field optional
		//
		// Note: Use this after functions that make fields required automatically
		//
		// Note: Field will be optional
		Optional() ruleSet
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
		Complex(rules Rules) ruleSet
		// If children of a slice is not struct or map, use this function and otherwise use Complex function after Slice function
		Children(rule ruleSet) ruleSet
		// Checks if input is a Specific type
		Type(input interface{}) ruleSet
		// Checks if input is at least 8 characters long, has one lowercase, one uppercase and one number character
		//
		// Note: Field will be required
		Password() ruleSet
		// Checks if at least one of passed ruleSets pass its own validation check
		OR(ruleSets ...ruleSet) ruleSet
		// Checks if XOR of all ruleSets passes
		XOR(ruleSets ...ruleSet) ruleSet
		// Gets a list of values and checks if input is one of them
		Choices(choices ...interface{}) ruleSet
		// Makes field required if at least one of passed fields are not empty, nil or zero(0, "", '')
		WhenExistOne(choices ...string) ruleSet
		// Makes field required if all passed fields are not empty, nil or zero(0, "", '')
		WhenExistAll(choices ...string) ruleSet
		// Checks if input is a string
		String() ruleSet
		// Sets messages for specific rules in current ruleSet
		SpecificMessages(specificMessages Messages)
		// Return specificMessages
		getSpecificMessages() Messages

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
		// Returns false if the ruleSet can be empty, nil or zero(0, "", '') and is allowed to not pass any validations
		isRequired() bool
		// Replaces passed validator with existing deepValidator
		setDeepValidator(input validator)
		// Returns deepValidator
		getDeepValidator() validator
		// Returns true if deepValidator is not nil
		hasDeepValidator() bool
		// Validates deepValidator
		validateDeepValidator(input interface{}) interface{}
		// Replaces passed validator with existing childrenValidator
		setChildrenValidator(input validator)
		// Returns childrenValidator
		getChildrenValidator() validator
		// Returns true if children is not nil
		hasChildrenValidator() bool
		// Validates childrenValidator
		validateChildrenValidator(input interface{}) interface{}
		// Returns requires
		getRequires() requires
		// Returns name
		getName() string
	}
)

func (o *ruleSetS) Int() ruleSet {
	functionName := "int"
	o.validators[functionName] = intRule
	return o
}

func (o *ruleSetS) Float() ruleSet {
	functionName := "float"
	o.validators[functionName] = floatRule
	return o
}

func (o *ruleSetS) Min(min float64) ruleSet {
	functionName := "min"
	o.validators[functionName] = minRule(min)
	precision := determinePrecision(min)
	o.addOption(functionName, "min", fmt.Sprintf("%."+precision+"f", min))
	if min > 0 {
		o.required()
	}
	return o
}

func (o *ruleSetS) Max(max float64) ruleSet {
	functionName := "max"
	o.validators[functionName] = maxRule(max)
	precision := determinePrecision(max)
	o.addOption(functionName, "max", fmt.Sprintf("%."+precision+"f", max))
	return o
}

func (o *ruleSetS) LenRange(from, to int) ruleSet {
	functionName := "len_range"
	o.validators[functionName] = lenRangeRule(from, to)
	o.addOption(functionName, "from", fmt.Sprintf("%d", from))
	o.addOption(functionName, "to", fmt.Sprintf("%d", to))
	if from > 0 {
		o.required()
	}
	return o
}

func (o *ruleSetS) Len(length int) ruleSet {
	functionName := "len"
	o.validators[functionName] = lenRule(length)
	o.addOption(functionName, "length", fmt.Sprint(length))
	if length > 0 {
		o.required()
	}
	return o
}

func (o *ruleSetS) Required() ruleSet {
	functionName := "required"
	o.validators[functionName] = requiredRule
	o.required()
	return o
}

func (o *ruleSetS) Optional() ruleSet {
	o.optional()
	return o
}

func (o *ruleSetS) NonZero() ruleSet {
	functionName := "non_zero"
	o.validators[functionName] = nonZeroRule
	o.required()
	return o
}

func (o *ruleSetS) NonNil() ruleSet {
	functionName := "non_nil"
	o.validators[functionName] = nonNilRule
	o.required()
	return o
}

func (o *ruleSetS) NonEmpty() ruleSet {
	functionName := "non_empty"
	o.validators[functionName] = nonEmptyRule
	o.required()
	return o
}

func (o *ruleSetS) Email() ruleSet {
	functionName := "email"
	o.validators[functionName] = emailRule
	return o
}

func (o *ruleSetS) Regex(pattern string) ruleSet {
	functionName := "regex"
	o.validators[functionName] = regexRule(pattern)
	o.addOption(functionName, "pattern", pattern)
	if !regexRule(pattern)("") {
		o.required()
	}
	return o
}

func (o *ruleSetS) Phone() ruleSet {
	functionName := "phone"
	o.validators[functionName] = phoneRule
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
	o.validators[functionName] = mapRule
	return o
}

func (o *ruleSetS) Slice() ruleSet {
	functionName := "slice"
	o.validators[functionName] = sliceRule
	return o
}

func (o *ruleSetS) Struct() ruleSet {
	functionName := "struct"
	o.validators[functionName] = structRule
	return o
}

func (o *ruleSetS) Complex(rules Rules) ruleSet {
	v := &validatorS{rule: nil, rules: rules}
	o.setDeepValidator(v)
	return o
}

func (o *ruleSetS) Children(rule ruleSet) ruleSet {
	v := &validatorS{rule: rule, rules: nil}
	o.childrenValidator = v
	return o
}

func (o *ruleSetS) Type(input interface{}) ruleSet {
	functionName := "type"
	switch v := input.(type) {
	case reflect.Type:
		o.addOption(functionName, "type", input.(reflect.Type).String())
		o.validators[functionName] = typeRule(v.String())
	default:
		o.addOption(functionName, "type", reflect.TypeOf(input).String())
		o.validators[functionName] = typeRule(reflect.TypeOf(input).String())
	}
	return o
}

func (o *ruleSetS) Password() ruleSet {
	functionName := "password"
	o.validators[functionName] = passwordRule
	o.required()
	return o
}

func (o *ruleSetS) OR(ruleSets ...ruleSet) ruleSet {
	functionName := "or"
	o.validators[functionName] = orRule(ruleSets...)
	return o
}

func (o *ruleSetS) XOR(ruleSets ...ruleSet) ruleSet {
	functionName := "xor"
	o.validators[functionName] = xorRule(ruleSets...)
	return o
}

func (o *ruleSetS) Choices(choices ...interface{}) ruleSet {
	functionName := "choices"
	o.validators[functionName] = choicesRule(choices...)
	o.addOption(functionName, "choices", strings.ReplaceAll(fmt.Sprint(choices), " ", ", "))
	return o
}

func (o *ruleSetS) WhenExistOne(choices ...string) ruleSet {
	functionName := "when_exist_one"
	o.requires[functionName] = whenExistOneRequireRule(choices...)
	o.validators[functionName] = requiredRule
	o.addOption(functionName, "choices", strings.ReplaceAll(fmt.Sprint(choices), " ", ", "))
	return o
}

func (o *ruleSetS) WhenExistAll(choices ...string) ruleSet {
	functionName := "when_exist_all"
	o.requires[functionName] = whenExistAllRequireRule(choices...)
	o.validators[functionName] = requiredRule
	o.addOption(functionName, "choices", strings.ReplaceAll(fmt.Sprint(choices), " ", ", "))
	return o
}

func (o *ruleSetS) String() ruleSet {
	functionName := "string"
	o.validators[functionName] = stringRule
	return o
}

func (o *ruleSetS) SpecificMessages(specificMessages Messages) {
	o.specificMessages = specificMessages
}

func (o *ruleSetS) getSpecificMessages() Messages {
	return o.specificMessages
}

func (o *ruleSetS) setChildrenValidator(input validator) {
	o.childrenValidator = input
}

func (o *ruleSetS) getChildrenValidator() validator {
	return o.childrenValidator
}

func (o *ruleSetS) hasChildrenValidator() bool {
	return o.childrenValidator != nil
}

func (o *ruleSetS) validateChildrenValidator(input interface{}) interface{} {
	return o.childrenValidator.Validate(input)
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

func (o *ruleSetS) setDeepValidator(input validator) {
	o.deepValidator = input
}

func (o *ruleSetS) getDeepValidator() validator {
	return o.deepValidator
}

func (o *ruleSetS) hasDeepValidator() bool {
	return o.deepValidator != nil
}

func (o *ruleSetS) validateDeepValidator(input interface{}) interface{} {
	return o.deepValidator.Validate(input)
}

func (o *ruleSetS) getRequires() requires {
	return o.requires
}

func (o *ruleSetS) getName() string {
	return o.name
}
