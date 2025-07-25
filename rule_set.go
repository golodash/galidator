package galidator

import (
	"context"
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
	Validators map[string]func(ctx context.Context, input interface{}) bool

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
		// map or struct
		deepValidator Validator
		// Defines type of elements of a slice
		childrenValidator Validator
		// Custom validators which is defined in generator
		customValidators *Validators
	}

	// An interface with some functions to satisfy validation purpose
	ruleSet interface {
		// Validates all validators defined
		validate(ctx context.Context, input interface{}) []string

		// Checks if input is int
		Int() ruleSet
		// Checks if input is float
		Float() ruleSet
		// Checks if input acts like: input >= min or len(input) >= min
		Min(min float64) ruleSet
		// Checks if input acts like: input <= max or len(input) <= max
		Max(max float64) ruleSet
		// Checks if input acts like: len(input) >= from && len(input) <= to
		//
		// If from == -1, no check on from will happen
		// If to == -1, no check on to will happen
		LenRange(from int, to int) ruleSet
		// Checks if input acts like: len(input) == length
		Len(length int) ruleSet
		// Makes this field required, Checks if input is not zero(0, "", ''), nil or empty
		Required() ruleSet
		// Makes this field optional
		Optional() ruleSet
		// If data is zero(0, "", ''), nil or empty, all validation rules will be checked anyway
		//
		// Note: Validations will not work by default if the value is zero(0, "", ''), nil or empty
		//
		// # Note: `Required`, `NonZero`, `NonNil` and `NonEmpty` will call this method by default
		//
		// Note: If not used, when an int is 0 and min is 5, min function will not activate
		//
		// Note: Field will be validated anyway from now on but if `Optional` method is called
		// after calling this method in continue of the current ruleSet chain, it has no effect
		AlwaysCheckRules() ruleSet
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
		Regex(pattern string) ruleSet
		// Checks if input is a valid phone number
		Phone() ruleSet
		// Adds custom validators
		Custom(validators Validators) ruleSet
		// Adds one custom validator which is registered before in generator
		RegisteredCustom(validatorKeys ...string) ruleSet
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
		// Makes field required if at least one of passed fields are empty, nil or zero(0, "", '')
		WhenNotExistOne(choices ...string) ruleSet
		// Makes field required if all passed fields are empty, nil or zero(0, "", '')
		WhenNotExistAll(choices ...string) ruleSet
		// Checks if input is a string
		String() ruleSet
		// Returns Validator of current Element (For map and struct elements)
		GetValidator() Validator
		// Returns Validator of children Elements (For slices)
		GetChildrenValidator() Validator

		// Adds a new pair of `key: value` message into existing SpecificMessages variable
		appendSpecificMessages(key string, value string)
		// Sets messages for specific rules in current ruleSet
		SpecificMessages(specificMessages Messages) ruleSet
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
		setDeepValidator(input Validator)
		// Returns deepValidator
		getDeepValidator() Validator
		// Returns true if deepValidator is not nil
		hasDeepValidator() bool
		// Validates deepValidator
		validateDeepValidator(ctx context.Context, input interface{}, translator Translator) interface{}
		// Replaces passed validator with existing childrenValidator
		setChildrenValidator(input Validator)
		// Returns childrenValidator
		getChildrenValidator() Validator
		// Returns true if children is not nil
		hasChildrenValidator() bool
		// Validates childrenValidator
		validateChildrenValidator(ctx context.Context, input interface{}, translator Translator) interface{}
		// Returns requires
		getRequires() requires
		// Returns name
		getName() string
		// Returns current validators + r.Validators
		appendRuleSet(r ruleSet) ruleSet
		// Returns passed argument name from struct if exist
		get(name string) interface{}
		// Sets passed argument value instead of existing in name parameter if exists
		set(name string, value interface{})
		// Sets custom validators which are defined in generator
		setGeneratorCustomValidators(validators *Validators) ruleSet
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
	return o
}

func (o *ruleSetS) Len(length int) ruleSet {
	functionName := "len"
	o.validators[functionName] = lenRule(length)
	o.addOption(functionName, "length", fmt.Sprint(length))
	return o
}

func (o *ruleSetS) AlwaysCheckRules() ruleSet {
	o.required()
	return o
}

func (o *ruleSetS) Required() ruleSet {
	functionName := "required"
	o.validators[functionName] = requiredRule
	return o.AlwaysCheckRules()
}

func (o *ruleSetS) Optional() ruleSet {
	o.optional()
	return o
}

func (o *ruleSetS) NonZero() ruleSet {
	functionName := "non_zero"
	o.validators[functionName] = nonZeroRule
	return o.AlwaysCheckRules()
}

func (o *ruleSetS) NonNil() ruleSet {
	functionName := "non_nil"
	o.validators[functionName] = nonNilRule
	return o.AlwaysCheckRules()
}

func (o *ruleSetS) NonEmpty() ruleSet {
	functionName := "non_empty"
	o.validators[functionName] = nonEmptyRule
	return o.AlwaysCheckRules()
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

func (o *ruleSetS) RegisteredCustom(validatorKeys ...string) ruleSet {
	for _, key := range validatorKeys {
		vs := *o.customValidators
		if _, ok := o.validators[key]; ok {
			panic(fmt.Sprintf("%s is duplicate and has to be unique", key))
		}
		if function, ok := vs[key]; ok {
			o.validators[key] = function
		} else {
			panic(fmt.Sprintf("%s custom validator doesn't exist, it is really defined in generator?", key))
		}
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
	if o.childrenValidator == nil {
		v := &validatorS{rule: rule, rules: nil}
		o.childrenValidator = v
	} else {
		o.childrenValidator.getRule().appendRuleSet(rule)
	}
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
	choicesString := []string{}
	for i := 0; i < len(choices); i++ {
		choicesString = append(choicesString, fmt.Sprint(choices[i]))
	}
	o.addOption(functionName, "choices", strings.ReplaceAll(fmt.Sprint(choicesString), " ", ", "))
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

func (o *ruleSetS) WhenNotExistOne(choices ...string) ruleSet {
	functionName := "when_not_exist_one"
	o.requires[functionName] = whenNotExistOneRequireRule(choices...)
	o.validators[functionName] = requiredRule
	o.addOption(functionName, "choices", strings.ReplaceAll(fmt.Sprint(choices), " ", ", "))
	return o
}

func (o *ruleSetS) WhenNotExistAll(choices ...string) ruleSet {
	functionName := "when_not_exist_all"
	o.requires[functionName] = whenNotExistAllRequireRule(choices...)
	o.validators[functionName] = requiredRule
	o.addOption(functionName, "choices", strings.ReplaceAll(fmt.Sprint(choices), " ", ", "))
	return o
}

func (o *ruleSetS) String() ruleSet {
	functionName := "string"
	o.validators[functionName] = stringRule
	return o
}

func (o *ruleSetS) GetValidator() Validator {
	return o.deepValidator
}

func (o *ruleSetS) GetChildrenValidator() Validator {
	return o.childrenValidator
}

func (o *ruleSetS) appendSpecificMessages(key string, value string) {
	var sm = o.getSpecificMessages()
	if sm == nil {
		sm = Messages{}
	}
	sm[key] = value

	o.SpecificMessages(sm)
}

func (o *ruleSetS) SpecificMessages(specificMessages Messages) ruleSet {
	o.specificMessages = specificMessages
	return o
}

func (o *ruleSetS) getSpecificMessages() Messages {
	return o.specificMessages
}

func (o *ruleSetS) setChildrenValidator(input Validator) {
	if o.childrenValidator != nil {
		o.childrenValidator.getRule().appendRuleSet(input.getRule())
	} else {
		o.childrenValidator = input
	}
}

func (o *ruleSetS) getChildrenValidator() Validator {
	return o.childrenValidator
}

func (o *ruleSetS) hasChildrenValidator() bool {
	return o.childrenValidator != nil
}

func (o *ruleSetS) validateChildrenValidator(ctx context.Context, input interface{}, translator Translator) interface{} {
	return o.childrenValidator.Validate(ctx, input, translator)
}

func (o *ruleSetS) validate(ctx context.Context, input interface{}) []string {
	fails := []string{}
	for key, vFunction := range o.validators {
		if !vFunction(ctx, input) {
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

func (o *ruleSetS) setDeepValidator(input Validator) {
	o.deepValidator = input
}

func (o *ruleSetS) getDeepValidator() Validator {
	return o.deepValidator
}

func (o *ruleSetS) hasDeepValidator() bool {
	return o.deepValidator != nil
}

func (o *ruleSetS) validateDeepValidator(ctx context.Context, input interface{}, translator Translator) interface{} {
	return o.deepValidator.Validate(ctx, input, translator)
}

func (o *ruleSetS) getRequires() requires {
	return o.requires
}

func (o *ruleSetS) getName() string {
	return o.name
}

func (o *ruleSetS) appendRuleSet(r ruleSet) ruleSet {
	rValidators := r.get("validators").(Validators)
	for key, value := range rValidators {
		o.validators[key] = value
	}
	rOptions := r.get("options").(options)
	for key, value := range rOptions {
		o.options[key] = value
	}
	rDeepValidator, ok := r.get("deepValidator").(Validator)
	if ok && rDeepValidator != nil && o.deepValidator == nil {
		o.deepValidator = rDeepValidator
	}
	rChildrenValidator, ok := r.get("childrenValidator").(Validator)
	if ok && rChildrenValidator != nil && o.childrenValidator == nil {
		o.childrenValidator = rChildrenValidator
	} else if ok && o.childrenValidator != nil && rChildrenValidator != nil {
		childrenValidatorRuleSet := rChildrenValidator.getRule()
		o.childrenValidator.getRule().appendRuleSet(childrenValidatorRuleSet)
	}
	rSpecificMessages := r.get("specificMessages").(Messages)
	for key, value := range rSpecificMessages {
		o.specificMessages[key] = value
	}
	if o.isOptional && !r.get("isOptional").(bool) {
		o.isOptional = false
	}
	name := r.get("name").(string)
	if name != "" && o.name == "" {
		o.name = name
	}
	rRequires := r.get("requires").(requires)
	for key, value := range rRequires {
		o.requires[key] = value
	}

	return o
}

func (o *ruleSetS) get(name string) interface{} {
	switch name {
	case "childrenValidator":
		return o.childrenValidator
	case "deepValidator":
		return o.deepValidator
	case "isOptional":
		return o.isOptional
	case "name":
		return o.name
	case "options":
		return o.options
	case "requires":
		return o.requires
	case "specificMessages":
		return o.specificMessages
	case "validators":
		return o.validators
	default:
		panic(fmt.Sprintf("there is no item as %s", name))
	}
}

func (o *ruleSetS) set(name string, value interface{}) {
	switch name {
	case "childrenValidator":
		o.childrenValidator = value.(Validator)
	case "deepValidator":
		o.deepValidator = value.(Validator)
	case "isOptional":
		o.isOptional = value.(bool)
	case "name":
		o.name = value.(string)
	case "options":
		o.options = value.(options)
	case "requires":
		o.requires = value.(requires)
	case "specificMessages":
		o.specificMessages = value.(Messages)
	case "validators":
		o.validators = value.(Validators)
	default:
		panic(fmt.Sprintf("there is no item as %s", name))
	}
}

func (o *ruleSetS) setGeneratorCustomValidators(validators *Validators) ruleSet {
	o.customValidators = validators
	return o
}
