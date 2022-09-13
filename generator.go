package galidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	gstrings "github.com/golodash/godash/strings"
)

type (
	// A struct to implement generator interface
	generatorS struct {
		defaultErrors map[string]string
	}

	// An interface to generate a validator or ruleSet
	generator interface {
		// Generates a validator interface which can be used to validate struct or map by some rules
		//
		// Please use CapitalCase for rules' keys (Important for getting data out of struct types)
		Validator(rules Rules, messages ...Messages) validator
		// Generates a validator interface which can be used to validate slice by one rule
		SliceValidator(r ruleSet, errorMessages ...Messages) validator
		// Generates a validator interface which can be used to validate struct or map by some rules
		//
		// Please use CapitalCase for rules' keys (Important for getting data out of struct types)
		ValidatorFromStruct(input interface{}) validator
		// Generates a ruleSet to validate passed information
		//
		// Variable name will be used in output of error keys
		RuleSet(name ...string) ruleSet
		// An alias for RuleSet function
		R(name ...string) ruleSet
	}
)

func (o *generatorS) Validator(rules Rules, errorMessages ...Messages) validator {
	var messages Messages = nil
	if len(errorMessages) != 0 {
		messages = errorMessages[0]
	}

	return &validatorS{rule: nil, rules: rules, messages: messages, specificMessages: SpecificMessages{}, defaultErrorMessages: o.defaultErrors}
}

func (o *generatorS) SliceValidator(r ruleSet, errorMessages ...Messages) validator {
	var messages Messages = nil
	if len(errorMessages) != 0 {
		messages = errorMessages[0]
	}

	if r == nil {
		panic("nil can not be passed")
	}
	if !r.hasChildrenValidator() {
		panic("passed ruleSet has no children")
	}

	return &validatorS{rule: r, rules: nil, messages: messages, specificMessages: SpecificMessages{}, defaultErrorMessages: o.defaultErrors}
}

func (o *generatorS) ValidatorFromStruct(input interface{}) validator {
	inputType := reflect.TypeOf(input)
	rules := Rules{}
	for reflect.ValueOf(input).Kind() == reflect.Ptr {
		input = reflect.ValueOf(input).Elem().Interface()
	}
	if inputType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("struct type should be passed and not %s", inputType.Kind()))
	}

	for i := 0; i < inputType.NumField(); i++ {
		element := inputType.Field(i)
		tags := []string{element.Tag.Get("g"), element.Tag.Get("galidator")}
		r := o.RuleSet(element.Tag.Get("json"))
		for _, fullTag := range tags {
			filters := strings.Split(fullTag, ",")
			for j := 0; j < len(filters); j++ {
				tag := strings.Split(filters[j], "=")
				funcName := gstrings.PascalCase(tag[0])
				parameters := []string{}
				if len(tag) == 2 {
					parameters = strings.Split(tag[1], "&")
				}

				if element.Type.Kind() == reflect.Struct {
					r.Complex(o.ValidatorFromStruct(element))
				}

				if element.Type.Kind() == reflect.Slice {
					addTypeCheck(r, element.Type.Kind())
					//! Attention Needed
					// r.Children(o.RuleSet())
				}

				switch funcName {
				case "Int":
					r.Int()
				case "Float":
					r.Float()
				case "Min":
					if len(parameters) == 1 {
						if p1, err := strconv.ParseFloat(parameters[0], 64); err == nil {
							r.Min(p1)
						}
					}
				case "Max":
					if len(parameters) == 1 {
						if p1, err := strconv.ParseFloat(parameters[0], 64); err == nil {
							r.Max(p1)
						}
					}
				case "LenRange":
					if len(parameters) == 2 {
						if p1, err := strconv.ParseInt(parameters[0], 10, 64); err == nil {
							if p2, err := strconv.ParseInt(parameters[1], 10, 64); err == nil {
								r.LenRange(int(p1), int(p2))
							}
						}
					}
				case "Len":
					if len(parameters) == 1 {
						if p1, err := strconv.ParseInt(parameters[0], 10, 64); err == nil {
							r.Len(int(p1))
						}
					}
				case "Required":
					r.Required()
				case "Optional":
					r.Optional()
				case "NonZero":
					r.NonZero()
				case "NonNil":
					r.NonNil()
				case "NonEmpty":
					r.NonEmpty()
				case "Email":
					r.Email()
				case "Regex":
					if len(parameters) > 1 {
						r.Regex(parameters[0])
					}
				case "Phone":
					r.Phone()
				case "Map":
					r.Map()
				case "Slice":
					r.Slice()
				case "Struct":
					r.Struct()
				case "Password":
					r.Password()
				case "OR", "Or":
					r.OR()
				case "XOR", "Xor":
					r.XOR()
				case "Choices":
					r.XOR()
				case "WhenExistOne":
					r.WhenExistOne()
				case "WhenExistAll":
					r.WhenExistAll()
				case "String":
					r.String()
				}
			}
		}
		rules[element.Name] = r
	}

	return o.Validator(rules)
}

func (o *generatorS) RuleSet(name ...string) ruleSet {
	var output = ""
	if len(name) != 0 {
		output = name[0]
	}
	return &ruleSetS{name: output, validators: Validators{}, requires: requires{}, options: options{}, isOptional: true, deepValidator: nil, childrenValidator: nil}
}

func (o *generatorS) R(name ...string) ruleSet {
	return o.RuleSet(name...)
}

// Returns a Validator Generator
func New(defaultErrors ...Messages) generator {
	defaultErrorsOutput := map[string]string{}
	for k, v := range defaultValidatorErrorMessages {
		defaultErrorsOutput[k] = v
	}
	if len(defaultErrors) > 0 {
		for k, v := range defaultErrors[0] {
			defaultErrorsOutput[k] = v
		}
	}

	return &generatorS{defaultErrors: defaultErrorsOutput}
}
