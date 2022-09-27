package galidator

import (
	"reflect"
	"strings"
)

type (
	// A struct to implement generator interface
	generatorS struct {
		customValidators Validators
	}

	// An interface to generate a validator or ruleSet
	generator interface {
		// Call this function before calling Validator function so that assigning custom validators can be possible
		CustomValidators(validators Validators) generator
		// Generates a validator interface which can be used to validate struct or map by some rules
		//
		// Please use CapitalCase for rules' keys (Important for getting data out of struct types)
		Validator(rule interface{}, messages ...Messages) Validator
		// Generates a validator interface which can be used to validate struct or map or slice by passing an instance of them
		validator(input interface{}) Validator
		// Generates a ruleSet to validate passed information
		//
		// Variable name will be used in output of error keys
		RuleSet(name ...string) ruleSet
		// An alias for RuleSet function
		R(name ...string) ruleSet
	}
)

func (o *generatorS) CustomValidators(validators Validators) generator {
	o.customValidators = validators
	return o
}

func (o *generatorS) Validator(rule interface{}, errorMessages ...Messages) Validator {
	var messages Messages = nil
	if len(errorMessages) != 0 {
		messages = errorMessages[0]
	}

	var output Validator = nil
	switch v := rule.(type) {
	case ruleSet:
		output = &validatorS{rule: v, rules: nil, messages: &messages}
	default:
		if reflect.TypeOf(v).Kind() == reflect.Struct || reflect.TypeOf(v).Kind() == reflect.Slice {
			output = o.validator(v)
		} else {
			panic("'rule' has to be a ruleSet or a struct instance")
		}
	}
	deepPassMessages(output, messages)

	return output
}

func (o *generatorS) validator(input interface{}) Validator {
	for reflect.ValueOf(input).Kind() == reflect.Ptr {
		input = reflect.ValueOf(input).Elem().Interface()
	}

	inputValue := reflect.ValueOf(input)
	inputType := reflect.TypeOf(input)
	r := o.RuleSet()
	if inputType.Kind() == reflect.Struct {
		rules := Rules{}
		for i := 0; i < inputType.NumField(); i++ {
			elementT := inputType.Field(i)
			element := inputValue.Field(i)
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

					normalFuncName, funcName := applyRules(r, tag, o, true)

					addSpecificMessage(r, funcName, elementT.Tag.Get(normalFuncName))
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

// Returns a Validator Generator
func New() generator {
	return &generatorS{}
}
