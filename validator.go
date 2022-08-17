package galidator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golodash/galidator/filters"
	gStrings "github.com/golodash/godash/strings"
)

type (
	// Stores keys of fields with their rules
	Rules map[string]rule

	// Stores custom error messages sent by user
	Messages map[string]string

	// A struct to implement validator interface
	validatorS struct {
		rules    Rules
		messages Messages
	}

	// Validator object
	validator interface {
		// Validates passed data and returns a map of possible validation errors happened on every field with failed validation.
		//
		// If no errors found, length of the output will be 0
		Validate(input interface{}) map[string][]string
	}
)

// Formats and returns error message associated with passed failKey
func getErrorMessage(fieldName string, failKey string, options option, messages Messages) string {
	if out, ok := messages[failKey]; ok {
		for key, value := range options {
			out = strings.ReplaceAll(out, "{"+key+"}", value)
		}
		return strings.ReplaceAll(out, "{field}", fieldName)
	} else {
		if defaultErrorMessage, ok := filters.DefaultValidatorErrorMessages[failKey]; ok {
			for key, value := range options {
				defaultErrorMessage = strings.ReplaceAll(defaultErrorMessage, "{"+key+"}", value)
			}
			return strings.ReplaceAll(defaultErrorMessage, "{field}", fieldName)
		} else {
			return fmt.Sprintf("error happened but no error message exists on '%s' rule", failKey)
		}
	}
}

func (o *validatorS) Validate(input interface{}) map[string][]string {
	output := map[string][]string{}
	inputValue := reflect.ValueOf(input)

	validate := func(role rule, onKeyInput interface{}, fieldName string) {
		valueOnKeyInput := reflect.ValueOf(onKeyInput)
		fails := role.validate(valueOnKeyInput.Interface())
		if len(fails) != 0 {
			for _, failKey := range fails {
				fieldName = gStrings.SnakeCase(fieldName)
				output[fieldName] = append(output[fieldName], getErrorMessage(fieldName, failKey, role.getOption(failKey), o.messages))
			}
		}
	}

	switch inputValue.Kind() {
	case reflect.Struct:
		for fieldName, role := range o.rules {
			valueOnKeyInput := inputValue.FieldByName(fieldName)
			validate(role, valueOnKeyInput.Interface(), fieldName)
		}
	case reflect.Map:
		for fieldName, role := range o.rules {
			valueOnKeyInput := inputValue.MapIndex(reflect.ValueOf(fieldName))
			validate(role, valueOnKeyInput.Interface(), fieldName)
		}
	}

	return output
}
