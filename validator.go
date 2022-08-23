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
func getErrorMessage(fieldName string, failKey string, value interface{}, options option, messages Messages, specificMessages SpecificMessages) string {
	snakeCaseFieldName := gStrings.SnakeCase(fieldName)
	if outMessages, ok := specificMessages[fieldName]; ok {
		if out, ok := outMessages[failKey]; ok {
			for key, value := range options {
				out = strings.ReplaceAll(out, "$"+key, value)
			}
			return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(out, "$fieldS", snakeCaseFieldName), "$field", fieldName), "$value", fmt.Sprint(value))
		}
	}

	if out, ok := messages[failKey]; ok {
		for key, value := range options {
			out = strings.ReplaceAll(out, "$"+key, value)
		}
		return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(out, "$fieldS", snakeCaseFieldName), "$field", fieldName), "$value", fmt.Sprint(value))
	} else {
		if defaultErrorMessage, ok := filters.DefaultValidatorErrorMessages[failKey]; ok {
			for key, value := range options {
				defaultErrorMessage = strings.ReplaceAll(defaultErrorMessage, "$"+key, value)
			}
			return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(defaultErrorMessage, "$fieldS", snakeCaseFieldName), "$field", fieldName), "$value", fmt.Sprint(value))
		} else {
			return fmt.Sprintf("error happened but no error message exists on '%s' rule", failKey)
		}
	}
}

func (o *validatorS) Validate(input interface{}) map[string]interface{} {
	output := map[string]interface{}{}
	inputValue := reflect.ValueOf(input)

	validate := func(rule ruleSet, onKeyInput interface{}, fieldName string) []string {
		halfOutput := []string{}
		fails := rule.validate(onKeyInput)
		if len(fails) != 0 {
			for _, failKey := range fails {
				halfOutput = append(halfOutput, getErrorMessage(fieldName, failKey, onKeyInput, rule.getOption(failKey), o.messages, o.specificMessages))
			}
		}

		if len(halfOutput) == 0 {
			return nil
		}

		return halfOutput
	}

	switch inputValue.Kind() {
	case reflect.Struct:
		for fieldName, rule := range o.rules {
			valueOnKeyInput := inputValue.FieldByName(fieldName)
			if valueOnKeyInput.IsValid() {
				value := valueOnKeyInput.Interface()
				if !rule.isRequired() && filters.IsEmptyNilZero(value) {
					continue
				}
				errors := validate(rule, value, fieldName)
				if errors != nil {
					output[fieldName] = errors
				}

				if rule.hasDeepValidator() && output[fieldName] == nil && (filters.Map(value) || filters.Struct(value)) {
					data := rule.validateDeepValidator(value)

					if len(data) != 0 {
						output[fieldName] = data
					}
				}

				if rule.hasChildrenRule() && output[fieldName] == nil && filters.Slice(value) {
					for i := 0; i < valueOnKeyInput.Len(); i++ {
						element := valueOnKeyInput.Index(i)
						childrenRule := rule.getChildrenRule()
						errors = validate(childrenRule, element.Interface(), fieldName)
						if errors != nil {
							if _, ok := output[fieldName]; !ok {
								output[fieldName] = map[string][]string{}
							}
							output[fieldName].(map[string][]string)[strconv.Itoa(i)] = errors
						}
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
				value := valueOnKeyInput.Interface()
				if !rule.isRequired() && filters.IsEmptyNilZero(value) {
					continue
				}
				errors := validate(rule, value, fieldName)
				if errors != nil {
					output[fieldName] = errors
				}

				if rule.hasDeepValidator() && output[fieldName] == nil && (filters.Map(value) || filters.Struct(value)) {
					data := rule.validateDeepValidator(value)

					if len(data) != 0 {
						output[fieldName] = data
					}
				}

				if rule.hasChildrenRule() && output[fieldName] == nil && filters.Slice(value) {
					for i := 0; i < valueOnKeyInput.Len(); i++ {
						element := valueOnKeyInput.Index(i)
						childrenRule := rule.getChildrenRule()
						errors = validate(childrenRule, element.Interface(), fieldName)
						if errors != nil {
							if _, ok := output[fieldName]; !ok {
								output[fieldName] = map[string][]string{}
							}
							output[fieldName].(map[string][]string)[strconv.Itoa(i)] = errors
						}
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
