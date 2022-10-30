package galidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	gStrings "github.com/golodash/godash/strings"
)

type (
	// To create a map of rules
	Rules map[string]ruleSet

	// To specify errors for rules
	Messages map[string]string

	// A struct to implement Validator interface
	validatorS struct {
		// Slice validator has this item full
		rule ruleSet
		// Stores keys of fields with their rules
		rules Rules
		// Stores custom error messages sent by user
		messages *Messages
	}

	// Validator interface
	Validator interface {
		// Validates passed data and returns a map of possible validation errors happened on every field with failed validation.
		//
		// If no errors found, output will be nil
		Validate(input interface{}) interface{}
		// Returns Rules
		getRules() Rules
		// Returns rule
		getRule() ruleSet
		// Replaces passed messages with existing one
		setMessages(messages Messages)
	}
)

// Formats and returns error message associated with passed ruleKey
func getErrorMessage(fieldName string, ruleKey string, value interface{}, options option, messages Messages, specificMessage Messages, defaultErrorMessages map[string]string) string {
	// Search for error message in specific messages for the rule
	if out, ok := specificMessage[gStrings.SnakeCase(ruleKey)]; len(specificMessage) != 0 && ok {
		for key, value := range options {
			out = strings.ReplaceAll(out, "$"+key, value)
		}
		return strings.ReplaceAll(strings.ReplaceAll(out, "$field", fieldName), "$value", fmt.Sprint(value))
	}

	// Search for error message in general messages for the validator
	if out, ok := messages[ruleKey]; ok {
		for key, value := range options {
			out = strings.ReplaceAll(out, "$"+key, value)
		}
		return strings.ReplaceAll(strings.ReplaceAll(out, "$field", fieldName), "$value", fmt.Sprint(value))
	} else {
		// If no error message found, search if there is a default error message for rule error
		if defaultErrorMessage, ok := defaultErrorMessages[ruleKey]; ok {
			for key, value := range options {
				defaultErrorMessage = strings.ReplaceAll(defaultErrorMessage, "$"+key, value)
			}
			return strings.ReplaceAll(strings.ReplaceAll(defaultErrorMessage, "$field", fieldName), "$value", fmt.Sprint(value))
		} else {
			// If no error message found, return that error message doesn't exist
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
				var m Messages = nil
				var sm Messages = ruleSet.getSpecificMessages()
				if o.messages != nil {
					m = *o.messages
				}
				halfOutput = append(halfOutput, getErrorMessage(fieldName, failKey, onKeyInput, ruleSet.getOption(failKey), m, sm, defaultValidatorErrorMessages))
			}
		}

		return halfOutput
	}

	if o.rules != nil {
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

					if reflect.ValueOf(data).IsValid() && reflect.ValueOf(data).Len() != 0 {
						output[fieldName] = data
					}
				}

				if ruleSet.hasChildrenValidator() && output[fieldName] == nil && sliceRule(value) {
					valueOnKeyInput = reflect.ValueOf(valueOnKeyInput.Interface())
					for i := 0; i < valueOnKeyInput.Len(); i++ {
						element := valueOnKeyInput.Index(i)
						errors := ruleSet.validateChildrenValidator(element.Interface())
						if reflect.ValueOf(errors).IsValid() && reflect.ValueOf(errors).Len() != 0 {
							if _, ok := output[fieldName]; !ok {
								output[fieldName] = map[string]interface{}{}
							}
							output[fieldName].(map[string]interface{})[strconv.Itoa(i)] = errors
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

					if reflect.ValueOf(data).IsValid() && reflect.ValueOf(data).Len() != 0 {
						output[fieldName] = data
					}
				}

				if ruleSet.hasChildrenValidator() && output[fieldName] == nil && sliceRule(value) {
					valueOnKeyInput = reflect.ValueOf(valueOnKeyInput.Interface())
					for i := 0; i < valueOnKeyInput.Len(); i++ {
						element := valueOnKeyInput.Index(i)
						errors := ruleSet.validateChildrenValidator(element.Interface())
						if reflect.ValueOf(errors).IsValid() && reflect.ValueOf(errors).Len() != 0 {
							if _, ok := output[fieldName]; !ok {
								output[fieldName] = map[string]interface{}{}
							}
							output[fieldName].(map[string]interface{})[strconv.Itoa(i)] = errors
						}
					}
				}
			}
		default:
			return []string{"invalid input"}
		}
	} else if o.rule != nil {
		if !o.rule.isRequired() && isEmptyNilZero(input) {
			return nil
		}

		errors := validate(o.rule, input, o.rule.getName())
		if len(errors) != 0 {
			return errors
		}

		switch inputValue.Kind() {
		case reflect.Slice:
			if o.rule.hasChildrenValidator() {
				for i := 0; i < inputValue.Len(); i++ {
					element := inputValue.Index(i)
					errors := o.rule.validateChildrenValidator(element.Interface())

					if reflect.ValueOf(errors).IsValid() && reflect.ValueOf(errors).Len() != 0 {
						if _, ok := output[strconv.Itoa(i)]; !ok {
							output[strconv.Itoa(i)] = errors
						}
					}
				}
			}
		default:
			if o.rule.hasDeepValidator() {
				errors := o.rule.validateDeepValidator(input)
				if reflect.ValueOf(errors).IsValid() && reflect.ValueOf(errors).Len() != 0 {
					return errors
				}
			}
		}
	} else {
		return []string{"invalid validator"}
	}

	if len(output) == 0 {
		return nil
	}

	return output
}

func (o *validatorS) getRules() Rules {
	return o.rules
}

func (o *validatorS) getRule() ruleSet {
	return o.rule
}

func (o *validatorS) setMessages(messages Messages) {
	o.messages = &messages
}
