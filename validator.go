package galidator

import (
	"fmt"
	"reflect"
	"strings"
)

type Messages map[string]string

type validatorS struct {
	items    Items
	messages Messages
}

func getErrorMessage(fieldName string, failKey string, messages Messages) string {
	failKey = strings.ToLower(failKey)
	if out, ok := messages[failKey]; ok {
		return out
	} else {
		if defaultKey, ok := defaultValidatorErrorMessages[failKey]; ok {
			return fmt.Sprintf(defaultKey, fieldName)
		} else {
			return fmt.Sprintf("error happened but no error message exists on '%s' item role", failKey)
		}
	}
}

func (o *validatorS) Validate(input interface{}) map[string][]string {
	output := map[string][]string{}
	inputValue := reflect.ValueOf(input)

	switch inputValue.Kind() {
	case reflect.Struct:
		for fieldName, validator := range o.items {
			valueOnKeyInput := inputValue.FieldByName(fieldName)
			fails := validator.validate(valueOnKeyInput.Interface())
			if len(fails) != 0 {
				for _, failKey := range fails {
					output[fieldName] = append(output[fieldName], getErrorMessage(fieldName, failKey, o.messages))
				}
			}
		}
	case reflect.Map:
		for fieldName, validator := range o.items {
			valueOnKeyInput := inputValue.MapIndex(reflect.ValueOf(fieldName))
			fails := validator.validate(valueOnKeyInput.Interface().(string))
			if len(fails) != 0 {
				for _, failKey := range fails {
					output[fieldName] = append(output[fieldName], getErrorMessage(fieldName, failKey, o.messages))
				}
			}
		}
	}

	return output
}

type validator interface {
	Validate(interface{}) map[string][]string
}
