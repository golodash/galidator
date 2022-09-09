package galidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
		Validate(input interface{}) interface{}
		// Adds more specific error messages for specific rules in specific fields
		AddSpecificMessages(fieldMessages SpecificMessages) validator
	}
)

// Formats and returns error message associated with passed ruleKey
func getErrorMessage(fieldName string, ruleKey string, value interface{}, options option, messages Messages, specificMessages SpecificMessages) string {
	if outMessages, ok := specificMessages[fieldName]; ok {
		if out, ok := outMessages[ruleKey]; ok {
			for key, value := range options {
				out = strings.ReplaceAll(out, "$"+key, value)
			}
			return strings.ReplaceAll(strings.ReplaceAll(out, "$field", fieldName), "$value", fmt.Sprint(value))
		}
	}

	if out, ok := messages[ruleKey]; ok {
		for key, value := range options {
			out = strings.ReplaceAll(out, "$"+key, value)
		}
		return strings.ReplaceAll(strings.ReplaceAll(out, "$field", fieldName), "$value", fmt.Sprint(value))
	} else {
		if defaultErrorMessage, ok := defaultValidatorErrorMessages[ruleKey]; ok {
			for key, value := range options {
				defaultErrorMessage = strings.ReplaceAll(defaultErrorMessage, "$"+key, value)
			}
			return strings.ReplaceAll(strings.ReplaceAll(defaultErrorMessage, "$field", fieldName), "$value", fmt.Sprint(value))
		} else {
			return fmt.Sprintf("error happened but no error message exists on '%s' rule key", ruleKey)
		}
	}
}

func (o *validatorS) Validate(input interface{}) interface{} {
	for reflect.ValueOf(input).Kind() == reflect.Ptr {
		input = reflect.ValueOf(input).Elem().Interface()
	}

	output := map[string]interface{}{}
	inputValue := reflect.ValueOf(input)

	validate := func(ruleSet ruleSet, onKeyInput interface{}, fieldName string) []string {
		halfOutput := []string{}
		fails := ruleSet.validate(onKeyInput)
		if len(fails) != 0 {
			for _, failKey := range fails {
				halfOutput = append(halfOutput, getErrorMessage(fieldName, failKey, onKeyInput, ruleSet.getOption(failKey), o.messages, o.specificMessages))
			}
		}

		return halfOutput
	}

	switch inputValue.Kind() {
	case reflect.Struct:
		for fieldName, ruleSet := range o.rules {
			valueOnKeyInput := inputValue.FieldByName(fieldName)
			if ruleSet.getName() != "" {
				fieldName = ruleSet.getName()
			}

			if !valueOnKeyInput.IsValid() {
				valueOnKeyInput = inputValue.MapIndex(reflect.ValueOf(fieldName))
			}
			if !valueOnKeyInput.IsValid() {
				panic(fmt.Sprintf("value on %s is not valid", fieldName))
			}

			value := valueOnKeyInput.Interface()
			// Just continue if no requires are set and field is empty, nil or zero
			requires, isRequired := determineRequires(input, value, ruleSet.getRequires())
			if (!ruleSet.isRequired() && !isRequired) && isEmptyNilZero(value) {
				continue
			}

			errors := validate(ruleSet, value, fieldName)
			dels := []int{}
			for i, key := range errors {
				if _, ok := requires[key]; ok {
					dels = append(dels, i)
				}
			}
			j := 0
			for _, key := range dels {
				errors = append(errors[:key-j], errors[key+1-j:]...)
				j++
			}
			if len(errors) != 0 {
				output[fieldName] = errors
			}

			if ruleSet.hasDeepValidator() && output[fieldName] == nil && (mapRule(value) || structRule(value) || sliceRule(value)) {
				data := ruleSet.validateDeepValidator(value)

				if reflect.ValueOf(data).Len() != 0 {
					output[fieldName] = data
				}
			}

			if ruleSet.hasChildrenRule() && output[fieldName] == nil && sliceRule(value) {
				for i := 0; i < valueOnKeyInput.Len(); i++ {
					element := valueOnKeyInput.Index(i)
					childrenRule := ruleSet.getChildrenRule()
					errors = validate(childrenRule, element.Interface(), fieldName)
					if len(errors) != 0 {
						if _, ok := output[fieldName]; !ok {
							output[fieldName] = map[string][]string{}
						}
						output[fieldName].(map[string][]string)[strconv.Itoa(i)] = errors
					}
				}
			}
		}
	case reflect.Map:
		for fieldName, ruleSet := range o.rules {
			valueOnKeyInput := inputValue.MapIndex(reflect.ValueOf(fieldName))
			if ruleSet.getName() != "" {
				fieldName = ruleSet.getName()
			}

			if !valueOnKeyInput.IsValid() {
				valueOnKeyInput = inputValue.MapIndex(reflect.ValueOf(fieldName))
			}
			if !valueOnKeyInput.IsValid() {
				panic(fmt.Sprintf("value on %s is not valid", fieldName))
			}

			value := valueOnKeyInput.Interface()
			if !ruleSet.isRequired() && isEmptyNilZero(value) {
				continue
			}
			errors := validate(ruleSet, value, fieldName)
			if len(errors) != 0 {
				output[fieldName] = errors
			}

			if ruleSet.hasDeepValidator() && output[fieldName] == nil && (mapRule(value) || structRule(value) || sliceRule(value)) {
				data := ruleSet.validateDeepValidator(value)

				if reflect.ValueOf(data).Len() != 0 {
					output[fieldName] = data
				}
			}

			if ruleSet.hasChildrenRule() && output[fieldName] == nil && sliceRule(value) {
				for i := 0; i < valueOnKeyInput.Len(); i++ {
					element := valueOnKeyInput.Index(i)
					childrenRule := ruleSet.getChildrenRule()
					errors = validate(childrenRule, element.Interface(), fieldName)
					if len(errors) != 0 {
						if _, ok := output[fieldName]; !ok {
							output[fieldName] = map[string][]string{}
						}
						output[fieldName].(map[string][]string)[strconv.Itoa(i)] = errors
					}
				}
			}
		}
	case reflect.Slice:
		for i := 0; i < inputValue.Len(); i++ {
			element := inputValue.Index(i)
			if element.IsValid() {
				data := o.Validate(element.Interface())
				if reflect.ValueOf(data).IsValid() && reflect.ValueOf(data).Len() != 0 {
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
			o.specificMessages[fieldKey][key] = errorMessage
			delete(errorMessages, key)
		}
	}

	return o
}
