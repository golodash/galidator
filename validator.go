package galidator

import (
	"fmt"
	"reflect"
	"strings"

	gStrings "github.com/golodash/godash/strings"
)

type Items map[string]item

type Messages map[string]string

type validatorS struct {
	items    Items
	messages Messages
}

func getErrorMessage(fieldName string, failKey string, options option, messages Messages) string {
	failKey = strings.ToLower(failKey)
	if out, ok := messages[failKey]; ok {
		for key, value := range options {
			out = strings.ReplaceAll(out, "{"+key+"}", value)
		}
		return strings.ReplaceAll(out, "{field}", fieldName)
	} else {
		if defaultErrorMessage, ok := defaultValidatorErrorMessages[failKey]; ok {
			for key, value := range options {
				defaultErrorMessage = strings.ReplaceAll(defaultErrorMessage, "{"+key+"}", value)
			}
			return strings.ReplaceAll(defaultErrorMessage, "{field}", fieldName)
		} else {
			return fmt.Sprintf("error happened but no error message exists on '%s' item role", failKey)
		}
	}
}

func (o *validatorS) Validate(input interface{}) map[string][]string {
	output := map[string][]string{}
	inputValue := reflect.ValueOf(input)

	validate := func(role item, onKeyInput interface{}, fieldName string) {
		valueOnKeyInput := reflect.ValueOf(onKeyInput)
		fails := role.validate(valueOnKeyInput.Interface())
		if len(fails) != 0 {
			for _, failKey := range fails {
				fieldName = gStrings.SnakeCase(fieldName)
				output[fieldName] = append(output[fieldName], getErrorMessage(fieldName, failKey, role.GetOption(failKey), o.messages))
			}
		}
	}

	switch inputValue.Kind() {
	case reflect.Struct:
		for fieldName, role := range o.items {
			valueOnKeyInput := inputValue.FieldByName(fieldName)
			validate(role, valueOnKeyInput.Interface(), fieldName)
		}
	case reflect.Map:
		for fieldName, role := range o.items {
			valueOnKeyInput := inputValue.MapIndex(reflect.ValueOf(fieldName))
			validate(role, valueOnKeyInput.Interface(), fieldName)
		}
	}

	return output
}

type validator interface {
	Validate(interface{}) map[string][]string
}
