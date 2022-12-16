package galidator

import (
	"reflect"
	"strings"

	"github.com/golodash/godash/generals"
)

type (
	// A struct to implement generator interface
	generatorS struct {
		// Custom validators
		customValidators Validators
		// Custom error messages
		messages Messages
	}

	// An interface to generate a validator or ruleSet
	generator interface {
		// Overrides current validators(if there is one) with passed validators
		//
		// Call this method before calling `generator.Validator` method to have effect
		CustomValidators(validators Validators) generator
		// Overrides current messages(if there is one) with passed messages
		//
		// Call this method before calling `generator.Validator` method to have effect
		CustomMessages(messages Messages) generator
		// Generates a validator interface which can be used to validate struct or map by some rules.
		//
		// `input` can be a ruleSet or a struct instance.
		//
		// Please use CapitalCase for rules' keys (Important for getting data out of struct types)
		Validator(input interface{}, messages ...Messages) Validator
		// Generates a validator interface which can be used to validate struct or map or slice by passing an instance of them
		validator(input interface{}) Validator
		// Generates a ruleSet to validate passed information
		//
		// Variable name will be used in output of error keys
		RuleSet(name ...string) ruleSet
		// An alias for RuleSet function
		R(name ...string) ruleSet
		// Generates a complex validator to validate maps and structs
		ComplexValidator(rules Rules, messages ...Messages) Validator
	}
)

func (o *generatorS) CustomValidators(validators Validators) generator {
	o.customValidators = validators
	return o
}

func (o *generatorS) CustomMessages(messages Messages) generator {
	o.messages = messages
	return o
}

func (o *generatorS) Validator(rule interface{}, errorMessages ...Messages) Validator {
	var messages Messages = o.messages
	if len(errorMessages) != 0 {
		for key, value := range errorMessages[0] {
			messages[key] = value
		}
	}
	switch rule.(type) {
	case ruleSet:
		break
	default:
		for reflect.TypeOf(rule).Kind() == reflect.Ptr {
			rule = reflect.ValueOf(rule).Elem().Interface()
		}
	}

	var output Validator = nil
	switch v := rule.(type) {
	case ruleSet:
		output = &validatorS{rule: v, rules: nil, messages: nil}
	default:
		if reflect.TypeOf(v).Kind() == reflect.Struct || reflect.TypeOf(v).Kind() == reflect.Slice {
			output = o.validator(v)
		} else {
			panic("'rule' has to be a ruleSet or a struct instance")
		}
	}
	deepPassMessages(output, generals.Duplicate(messages).(Messages))

	return output
}

func (o *generatorS) validator(input interface{}) Validator {
	inputValue := reflect.ValueOf(input)
	inputType := reflect.TypeOf(input)
	r := o.RuleSet()
	if inputType.Kind() == reflect.Struct {
		rules := Rules{}
		for i := 0; i < inputType.NumField(); i++ {
			elementT := inputType.Field(i)
			element := inputValue.Field(i)
			for element.Type().Kind() == reflect.Ptr {
				elementType := element.Type().Elem()
				element = reflect.New(elementType).Elem()
				elementT.Type = element.Type()
			}
			tags := []string{elementT.Tag.Get("g"), elementT.Tag.Get("galidator")}
			r = o.RuleSet(elementT.Tag.Get("json"))
			addTypeCheck(r, elementT.Type.Kind())

			if elementT.Type.Kind() == reflect.Struct || elementT.Type.Kind() == reflect.Map {
				validator := o.validator(element.Interface())
				r.setDeepValidator(validator)
			} else if elementT.Type.Kind() == reflect.Slice {
				child := elementT.Type.Elem()
				if child.Kind() != reflect.Slice && child.Kind() != reflect.Struct && child.Kind() != reflect.Map {
					r.Children(o.R().Type(child))
				} else {
					validator := o.validator(reflect.Zero(elementT.Type.Elem()).Interface())
					r.setChildrenValidator(validator)
				}
			}

			for _, fullTag := range tags {
				filters := strings.Split(fullTag, ",")
				for j := 0; j < len(filters); j++ {
					tag := strings.SplitN(filters[j], "=", 2)

					normalFuncName := applyRules(r, tag, o, true)

					addSpecificMessage(r, normalFuncName, elementT.Tag.Get(normalFuncName))
				}
			}

			// Support for binding tag that is used for Bind actions
			if bindingTags := elementT.Tag.Get("binding"); bindingTags != "" {
				splits := strings.Split(bindingTags, ",")
				for j := 0; j < len(splits); j++ {
					addSpecificMessage(r, splits[j], elementT.Tag.Get(splits[j]))
				}
			}
			rules[elementT.Name] = r
		}

		return &validatorS{rules: rules}
	} else if inputType.Kind() == reflect.Slice {
		child := inputType.Elem()
		if child.Kind() != reflect.Slice && child.Kind() != reflect.Struct && child.Kind() != reflect.Map {
			r.Children(o.R().Type(child))
		} else {
			validator := o.validator(reflect.Zero(child).Interface())
			r.setChildrenValidator(validator)
		}

		return &validatorS{rule: r}
	} else if inputType.Kind() == reflect.Map {
		panic("do not use map in Validator creation based on struct elements tags")
	} else {
		r.Type(inputType)

		return &validatorS{rule: r}
	}
}

func (o *generatorS) RuleSet(name ...string) ruleSet {
	var output = ""
	if len(name) != 0 {
		output = name[0]
	}
	return &ruleSetS{name: output, validators: Validators{}, requires: requires{}, options: options{}, isOptional: true}
}

func (o *generatorS) R(name ...string) ruleSet {
	return o.RuleSet(name...)
}

func (o *generatorS) ComplexValidator(rules Rules, errorMessages ...Messages) Validator {
	var messages Messages = o.messages
	if len(errorMessages) != 0 {
		for key, value := range errorMessages[0] {
			messages[key] = value
		}
	}
	output := &validatorS{rule: nil, rules: rules, messages: nil}

	deepPassMessages(output, generals.Duplicate(messages).(Messages))
	return output
}

// Returns a new Generator
func NewGenerator() generator {
	return &generatorS{
		messages:         Messages{},
		customValidators: Validators{},
	}
}

// An alias for `NewGenerator` funcion
func New() generator {
	return NewGenerator()
}

// An alias for `NewGenerator` funcion
func G() generator {
	return NewGenerator()
}
