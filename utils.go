package galidator

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	gStrings "github.com/golodash/godash/strings"
)

// Determines the precision of a float number for print
func determinePrecision(number float64) string {
	for i := 0; ; i++ {
		ten := math.Pow10(i)
		if math.Floor(ten*number) == ten*number {
			return fmt.Sprint(i)
		}
	}
}

// Returns true if input is nil
func isNil(input interface{}) bool {
	return !reflect.ValueOf(input).IsValid() || input == nil || (reflect.TypeOf(input).Kind() == reflect.Ptr && reflect.ValueOf(input).IsNil())
}

// Returns true if input is map or slice and has 0 elements
//
// Returns false if input is not map or slice
func hasZeroItems(input interface{}) bool {
	inputValue := reflect.ValueOf(input)
	switch inputValue.Kind() {
	case reflect.Slice, reflect.Map:
		return inputValue.Len() == 0
	default:
		return false
	}
}

// A helper function to use inside code
func isEmptyNilZero(input interface{}) bool {
	return !requiredRule(context.TODO(), input)
}

// Returns values of passed fields from passed struct or map
func getValues(all interface{}, fields ...string) []interface{} {
	fieldsValues := []interface{}{}
	allValue := reflect.ValueOf(all)

	if allValue.Kind() == reflect.Map {
		for _, key := range fields {
			element := allValue.MapIndex(reflect.ValueOf(key))
			if !element.IsValid() {
				element = allValue.MapIndex(reflect.ValueOf(key))
			}

			if !element.IsValid() {
				panic(fmt.Sprintf("value on %s field is not valid", key))
			}

			fieldsValues = append(fieldsValues, element.Interface())
		}
	} else if allValue.Kind() == reflect.Struct {
		for _, key := range fields {
			element := allValue.FieldByName(key)
			if !element.IsValid() {
				panic(fmt.Sprintf("value on %s field is not valid", key))
			}

			fieldsValues = append(fieldsValues, element.Interface())
		}
	}

	return fieldsValues
}

// Returns a list of keys for requires which determine not required and a bool which determines if we need to validate or not
func determineRequires(all interface{}, input interface{}, requires requires) (map[string]interface{}, bool) {
	output := map[string]interface{}{}
	if len(requires) == 0 {
		return output, false
	}
	for key, req := range requires {
		if !req(all)(input) {
			output[key] = 1
		}
	}

	return output, len(output) == len(requires)
}

// Passes messages to other validators
func deepPassMessages(v Validator, messages *Messages) {
	v.setMessages(messages)
	r := v.getRule()
	if r != nil {
		if v1 := r.getChildrenValidator(); v1 != nil {
			deepPassMessages(v1, messages)
		}
		if v2 := r.getDeepValidator(); v2 != nil {
			deepPassMessages(v2, messages)
		}
	}
	rs := v.getRules()
	if rs == nil {
		return
	}
	for _, r := range rs {
		if r != nil {
			if v1 := r.getChildrenValidator(); v1 != nil {
				deepPassMessages(v1, messages)
			}
			if v2 := r.getDeepValidator(); v2 != nil {
				deepPassMessages(v2, messages)
			}
		}
	}
}

// Adds one specific message to passed ruleSet if message is not a empty string
func addSpecificMessage(r ruleSet, funcName, message string) {
	splits := strings.SplitN(funcName, ".", 2)
	if firstElementSnake := gStrings.SnakeCase(splits[0]); len(splits) > 1 && (firstElementSnake == "c" || firstElementSnake == "child") {
		addSpecificMessage(r.getChildrenValidator().getRule(), splits[1], message)
		return
	}
	funcName = gStrings.SnakeCase(funcName)
	if message != "" {
		r.appendSpecificMessages(funcName, message)
	}
}

// Adds rules which are inside passed slice of strings called tag
func applyRules(r ruleSet, tag []string, o *generatorS, orXor bool) (normalFuncName string) {
	normalFuncName = strings.TrimSpace(tag[0])
	funcName := normalFuncName
	splits := strings.SplitN(normalFuncName, ".", 2)
	if firstElementSnake := gStrings.SnakeCase(splits[0]); len(splits) > 1 && (firstElementSnake == "c" || firstElementSnake == "child") {
		if r.getChildrenValidator() == nil {
			r.Children(o.R())
		}
		if len(tag) > 1 {
			applyRules(r.getChildrenValidator().getRule(), []string{splits[1], tag[1]}, o, orXor)
		} else {
			applyRules(r.getChildrenValidator().getRule(), []string{splits[1]}, o, orXor)
		}
		return normalFuncName
	} else if len(splits) == 1 {
		funcName = gStrings.PascalCase(normalFuncName)
	} else {
		panic(fmt.Sprintf("can't understand %s rule", normalFuncName))
	}

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
	case "Or", "OR":
		if orXor {
			rules := []ruleSet{}
			parametersSeparated := strings.Split(tag[1], "|")
			for _, parameters := range parametersSeparated {
				rule := o.R()
				parameter := strings.Split(parameters, "+")
				for _, p := range parameter {
					applyRules(rule, strings.SplitN(p, "=", 2), o, false)
				}
				rules = append(rules, rule)
			}

			r.OR(rules...)
		} else {
			panic("OR or XOR inside another OR or XOR is not possible")
		}
	case "Xor", "XOr", "XOR":
		if orXor {
			rules := []ruleSet{}
			parametersSeparated := strings.Split(tag[1], "|")
			for _, parameters := range parametersSeparated {
				rule := o.R()
				parameter := strings.Split(parameters, "+")
				for _, p := range parameter {
					applyRules(rule, strings.SplitN(p, "=", 2), o, false)
				}
				rules = append(rules, rule)
			}

			r.XOR(rules...)
		} else {
			panic("OR or XOR inside another OR or XOR is not possible")
		}
	case "Choices":
		params := []interface{}{}
		for _, item := range parameters {
			params = append(params, item)
		}
		r.Choices(params...)
	case "WhenExistOne":
		r.WhenExistOne(parameters...)
	case "WhenExistAll":
		r.WhenExistAll(parameters...)
	case "WhenNotExistOne":
		r.WhenNotExistOne(parameters...)
	case "WhenNotExistAll":
		r.WhenNotExistAll(parameters...)
	case "String":
		r.String()
	case "Children", "Custom", "Complex", "Type":
		panic(fmt.Sprintf("take a look at documentations, %s rule does not work in tags like this", funcName))
	default:
		if normalFuncName != "" {
			if function, ok := o.customValidators[normalFuncName]; ok {
				r.Custom(Validators{
					normalFuncName: function,
				})
			} else {
				panic(fmt.Sprintf("%s custom validator did not find, call CustomValidators function before calling Validator function", normalFuncName))
			}
		}
	}

	return
}
