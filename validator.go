package galidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	filters "github.com/golodash/galidator/internal"
	gStrings "github.com/golodash/godash/strings"
)

type (
	// To create a map of rules
	Rules map[string]ruleSet

	// To specify errors for rules
	Messages map[string]string

	// To specify specific errors in specific rules
	SpecificMessages map[string]Messages

	// A struct to implement validator interface
	validatorS struct {
		// Stores keys of fields with their rules
		rules Rules
		// Stores custom error messages sent by user
		messages Messages
		// Stores custom field error messages sent by user
		specificMessages SpecificMessages
	}

	// Validator object
	validator interface {
		// Validates passed data and returns a map of possible validation errors happened on every field with failed validation.
		//
		// If no errors found, length of the output will be 0
		Validate(input interface{}) map[string]interface{}
		// Adds more specific error messages for specific rules in specific fields
		AddSpecificMessages(fieldMessages SpecificMessages) validator
	}
)

// Formats and returns error message associated with passed failKey
func getErrorMessage(fieldName string, failKey string, options option, messages Messages, specificMessages SpecificMessages) string {
	snakeCaseFieldName := gStrings.SnakeCase(fieldName)
	if outMessages, ok := specificMessages[fieldName]; ok {
		if out, ok := outMessages[failKey]; ok {
			for key, value := range options {
				out = strings.ReplaceAll(out, "$"+key, value)
			}
			return strings.ReplaceAll(strings.ReplaceAll(out, "$field", fieldName), "$fieldS", snakeCaseFieldName)
		}
	}

	if out, ok := messages[failKey]; ok {
		for key, value := range options {
			out = strings.ReplaceAll(out, "$"+key, value)
		}
		return strings.ReplaceAll(strings.ReplaceAll(out, "$field", fieldName), "$fieldS", snakeCaseFieldName)
	} else {
		if defaultErrorMessage, ok := filters.DefaultValidatorErrorMessages[failKey]; ok {
			for key, value := range options {
				defaultErrorMessage = strings.ReplaceAll(defaultErrorMessage, "$"+key, value)
			}
			return strings.ReplaceAll(strings.ReplaceAll(defaultErrorMessage, "$field", fieldName), "$fieldS", snakeCaseFieldName)
		} else {
			return fmt.Sprintf("error happened but no error message exists on '%s' rule", failKey)
		}
	}
}

func (o *validatorS) Validate(input interface{}) map[string]interface{} {
	output := map[string]interface{}{}
	inputValue := reflect.ValueOf(input)

	validate := func(rule ruleSet, onKeyInput interface{}, fieldName string) {
		valueOnKeyInput := reflect.ValueOf(onKeyInput)
		fails := rule.validate(valueOnKeyInput.Interface())
		if len(fails) != 0 {
			if output[fieldName] == nil {
				output[fieldName] = []string{}
			}
			for _, failKey := range fails {
				output[fieldName] = append(output[fieldName].([]string), getErrorMessage(fieldName, failKey, rule.getOption(failKey), o.messages, o.specificMessages))
			}
		}
	}

	switch inputValue.Kind() {
	case reflect.Struct:
		for fieldName, rule := range o.rules {
			valueOnKeyInput := inputValue.FieldByName(fieldName)
			if valueOnKeyInput.IsValid() {
				if !rule.isRequired() && filters.IsEmptyNilZero(valueOnKeyInput.Interface()) {
					continue
				}
				validate(rule, valueOnKeyInput.Interface(), fieldName)

				if rule.hasDeepValidator() && output[fieldName] == nil {
					data := rule.validateDeepValidator(valueOnKeyInput.Interface())

					if len(data) != 0 {
						output[fieldName] = data
					}
				}
			} else {
				panic(fmt.Sprintf("value on %s is not valid", fieldName))
			}
		}
	case reflect.Map:
		for fieldName, rule := range o.rules {
			valueOnKeyInput := inputValue.MapIndex(reflect.ValueOf(fieldName))
			if !valueOnKeyInput.IsValid() {
				valueOnKeyInput = inputValue.MapIndex(reflect.ValueOf(gStrings.SnakeCase(fieldName)))
			}
			if valueOnKeyInput.IsValid() {
				if !rule.isRequired() && filters.IsEmptyNilZero(valueOnKeyInput.Interface()) {
					continue
				}
				validate(rule, valueOnKeyInput.Interface(), fieldName)

				if rule.hasDeepValidator() && output[fieldName] == nil {
					data := rule.validateDeepValidator(valueOnKeyInput.Interface())

					if len(data) != 0 {
						output[fieldName] = data
					}
				}
			} else {
				panic(fmt.Sprintf("value on %s is not valid", fieldName))
			}
		}
	case reflect.Slice:
		for i := 0; i < inputValue.Len(); i++ {
			element := inputValue.Index(i)
			if element.IsValid() {
				data := o.Validate(element.Interface())
				if len(data) != 0 {
					output[strconv.Itoa(i)] = data
				}
			}
		}
	}

	if len(output) == 0 {
		return nil
	}

	return output
}

func (o *validatorS) AddSpecificMessages(fieldMessages SpecificMessages) validator {
	for fieldKey, errorMessages := range fieldMessages {
		if _, ok := o.rules[fieldKey]; !ok {
			continue
		}
		if _, ok := o.specificMessages[fieldKey]; !ok {
			o.specificMessages[fieldKey] = Messages{}
		}
		for key, errorMessage := range errorMessages {
			updatedKey := gStrings.SnakeCase(key)
			if updatedKey == key {
				continue
			}

			o.specificMessages[fieldKey][updatedKey] = errorMessage
			delete(errorMessages, key)
		}
	}

	return o
}
