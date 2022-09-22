package galidator

import (
	"reflect"
	"strconv"
	"strings"

	gStrings "github.com/golodash/godash/strings"
)

type (
	// A struct to implement generator interface
	generatorS struct{}

	// An interface to generate a validator or ruleSet
	generator interface {
		// Generates a validator interface which can be used to validate struct or map by some rules
		//
		// Please use CapitalCase for rules' keys (Important for getting data out of struct types)
		Validator(rule interface{}, messages ...Messages) validator
		// Generates a validator interface which can be used to validate struct or map or slice by passing an instance of them
		validator(input interface{}) validator
		// Generates a ruleSet to validate passed information
		//
		// Variable name will be used in output of error keys
		RuleSet(name ...string) ruleSet
		// An alias for RuleSet function
		R(name ...string) ruleSet
	}
)

func (o *generatorS) Validator(rule interface{}, errorMessages ...Messages) validator {
	var messages Messages = nil
	if len(errorMessages) != 0 {
		messages = errorMessages[0]
	}

	switch v := rule.(type) {
	case ruleSet:
		return &validatorS{rule: v, rules: nil, messages: messages, specificMessages: SpecificMessages{}}
	default:
		if reflect.TypeOf(v).Kind() == reflect.Struct {
			return o.validator(v)
		}
		panic("'rule' has to be a ruleSet or a struct instance")
	}
}

func (o *generatorS) validator(input interface{}) validator {
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
					tag := strings.Split(filters[j], "=")
					funcName := gStrings.PascalCase(tag[0])
					parameters := []string{}
					if len(tag) == 2 {
						parameters = strings.Split(tag[1], "&")
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
					case "LenRange", "Lenrange":
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
					case "NonZero", "Nonzero":
						r.NonZero()
					case "NonNil", "Nonnil":
						r.NonNil()
					case "NonEmpty", "Nonempty":
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
					case "Or", "OR":
						r.OR()
					case "Xor", "XOr", "XOR":
						r.XOR()
					case "Choices":
						r.XOR()
					case "WhenExistOne", "Whenexistone":
						r.WhenExistOne()
					case "WhenExistAll", "Whenexistall":
						r.WhenExistAll()
					case "String":
						r.String()
					}
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
	return &ruleSetS{name: output, validators: Validators{}, requires: requires{}, options: options{}, isOptional: true, deepValidator: nil, childrenValidator: nil}
}

func (o *generatorS) R(name ...string) ruleSet {
	return o.RuleSet(name...)
}

// Returns a Validator Generator
func New() generator {
	return &generatorS{}
}
