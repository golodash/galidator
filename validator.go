package galidator

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	playgroundValidator "github.com/go-playground/validator/v10"
	gStrings "github.com/golodash/godash/strings"
)

type (
	// To create a map of rules
	Rules map[string]ruleSet

	// To specify errors for rules
	Messages map[string]string

	Translator func(string) string

	// A struct to implement Validator interface
	validatorS struct {
		// Slice validator has this item full
		rule ruleSet
		// Stores keys of fields with their rules
		rules Rules
		// Stores custom error messages sent by user
		messages *Messages
	}

	// Used just in decryptErrors function
	sliceValidationError []error

	// Validator interface
	Validator interface {
		// Validates passed data and returns a map of possible validation errors happened on every field with failed validation.
		//
		// If no errors found, output will be nil
		Validate(ctx context.Context, input interface{}, translator ...Translator) interface{}
		// Decrypts errors returned from gin's Bind process and returns proper error messages
		//
		// If returnUnmarshalErrorContext is true (default is true), if an error happened when
		// unmarshal process was happening, `UnmarshalError` message will return like: {"message":"unmarshal error"}
		//
		// If returnUnmarshalErrorContext is false, actual error message which tells you what went wrong
		// will return
		DecryptErrors(err error, returnUnmarshalErrorContext ...bool) interface{}
		// Sets passed default values if value field is nil
		SetDefaultOnNil(input interface{}, defaultValue interface{})
		// Sets passed default values if value field is zero
		SetDefault(input interface{}, defaultValue interface{})
		// Returns ruleSet of an input
		GetStructRule(input string) ruleSet
		// Returns the ruleSet of current validator
		GetRule() ruleSet
		// Returns Rules
		getRules() Rules
		// Returns rule
		getRule() ruleSet
		// Replaces passed messages with existing one
		setMessages(messages *Messages)
		// Returns messages
		getMessages() *Messages
	}
)

const (
	// Error message that will return when UnmarshalTypeError happens
	UnmarshalError = "unmarshal error"
)

func (err sliceValidationError) Error() string {
	return "error"
}

func getFormattedErrorMessage(message string, fieldName string, value interface{}, options option, translator ...Translator) string {
	var t Translator = nil
	if len(translator) != 0 {
		t = translator[0]
	}
	for key, value := range options {
		message = strings.ReplaceAll(message, "$"+key, value)
	}
	newValue := fmt.Sprint(value)
	if t != nil {
		newValue = t(newValue)
	}
	return strings.ReplaceAll(strings.ReplaceAll(message, "$field", fieldName), "$value", newValue)
}

// Formats and returns error message associated with passed ruleKey
func getRawErrorMessage(ruleKey string, messages Messages, specificMessage Messages, defaultErrorMessages map[string]string) string {
	// Search for error message in specific messages for the rule
	if out, ok := specificMessage[gStrings.SnakeCase(ruleKey)]; len(specificMessage) != 0 && ok {
		return out
	}

	// Search for error message in general messages for the validator
	if out, ok := messages[ruleKey]; ok {
		return out
	} else {
		// If no error message found, search if there is a default error message for rule error
		if defaultErrorMessage, ok := defaultErrorMessages[ruleKey]; ok {
			return defaultErrorMessage
		} else {
			// If no error message found, return that error message doesn't exist
			return fmt.Sprintf("error happened but no error message exists on '%s' rule key", ruleKey)
		}
	}
}

func (o *validatorS) GetRule() ruleSet {
	return o.rule
}

func (o *validatorS) GetStructRule(input string) ruleSet {
	if rule, ok := o.rules[input]; ok {
		return rule
	}
	return nil
}

func (o *validatorS) Validate(ctx context.Context, input interface{}, translator ...Translator) interface{} {
	for reflect.ValueOf(input).Kind() == reflect.Ptr {
		input = reflect.ValueOf(input).Elem().Interface()
	}
	var t Translator = nil
	if len(translator) != 0 {
		t = translator[0]
	}

	output := map[string]interface{}{}
	inputValue := reflect.ValueOf(input)

	validate := func(ruleSet ruleSet, onKeyInput interface{}, fieldName string) []string {
		for reflect.ValueOf(onKeyInput).IsValid() && reflect.TypeOf(onKeyInput).Kind() == reflect.Ptr {
			onKeyInputValue := reflect.ValueOf(onKeyInput).Elem()
			if onKeyInputValue.IsValid() {
				onKeyInput = onKeyInputValue.Interface()
			} else {
				onKeyInput = nil
			}
		}

		halfOutput := []string{}
		fails := ruleSet.validate(ctx, onKeyInput)
		if len(fails) != 0 {
			for _, failKey := range fails {
				var m Messages = nil
				var sm Messages = ruleSet.getSpecificMessages()
				if o.messages != nil {
					m = *o.messages
				}
				message := getRawErrorMessage(failKey, m, sm, defaultValidatorErrorMessages)
				if t != nil {
					message = t(message)
				}
				message = getFormattedErrorMessage(message, fieldName, onKeyInput, ruleSet.getOption(failKey), t)
				halfOutput = append(halfOutput, message)
			}
		}

		return halfOutput
	}

	if o.rules != nil {
		switch inputValue.Kind() {
		case reflect.Struct:
			for fieldName, ruleSet := range o.rules {
				valueOnKeyInput := inputValue.FieldByName(fieldName)
				typeOnKeyInput, found := inputValue.Type().FieldByName(fieldName)
				if ruleSet.getName() != "" {
					fieldName = ruleSet.getName()
				}

				if !valueOnKeyInput.IsValid() {
					valueOnKeyInput = inputValue.FieldByName(fieldName)
				}
				if !valueOnKeyInput.IsValid() {
					panic(fmt.Sprintf("value on %s is not valid", fieldName))
				}

				// If not exported
				if !found || typeOnKeyInput.PkgPath != "" {
					break
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

				if ruleSet.hasDeepValidator() && output[fieldName] == nil && (mapRule(ctx, value) || structRule(ctx, value) || sliceRule(ctx, value)) {
					data := ruleSet.validateDeepValidator(ctx, value, t)

					if reflect.ValueOf(data).IsValid() && reflect.ValueOf(data).Len() != 0 {
						output[fieldName] = data
					}
				}

				if ruleSet.hasChildrenValidator() && output[fieldName] == nil && sliceRule(ctx, value) {
					valueOnKeyInput = reflect.ValueOf(valueOnKeyInput.Interface())
					for i := 0; i < valueOnKeyInput.Len(); i++ {
						element := valueOnKeyInput.Index(i)
						errors := ruleSet.validateChildrenValidator(ctx, element.Interface(), t)
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

				if ruleSet.hasDeepValidator() && output[fieldName] == nil && (mapRule(ctx, value) || structRule(ctx, value) || sliceRule(ctx, value)) {
					data := ruleSet.validateDeepValidator(ctx, value, t)

					if reflect.ValueOf(data).IsValid() && reflect.ValueOf(data).Len() != 0 {
						output[fieldName] = data
					}
				}

				if ruleSet.hasChildrenValidator() && output[fieldName] == nil && sliceRule(ctx, value) {
					valueOnKeyInput = reflect.ValueOf(valueOnKeyInput.Interface())
					for i := 0; i < valueOnKeyInput.Len(); i++ {
						element := valueOnKeyInput.Index(i)
						errors := ruleSet.validateChildrenValidator(ctx, element.Interface(), t)
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
					errors := o.rule.validateChildrenValidator(ctx, element.Interface(), t)

					if reflect.ValueOf(errors).IsValid() && reflect.ValueOf(errors).Len() != 0 {
						if _, ok := output[strconv.Itoa(i)]; !ok {
							output[strconv.Itoa(i)] = errors
						}
					}
				}
			}
		default:
			if o.rule.hasDeepValidator() {
				errors := o.rule.validateDeepValidator(ctx, input, t)
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

func (o *validatorS) getMessages() *Messages {
	return o.messages
}

func (o *validatorS) DecryptErrors(err error, returnUnmarshalErrorContext ...bool) interface{} {
	if err == nil {
		return nil
	}

	unmarshalError := false
	if len(returnUnmarshalErrorContext) > 0 {
		unmarshalError = returnUnmarshalErrorContext[0]
	}

	return decryptErrors(err, o, unmarshalError)
}

func (o *validatorS) SetDefaultOnNil(input interface{}, defaultValue interface{}) {
	if reflect.TypeOf(input).Kind() != reflect.Ptr {
		panic("Please send data as a pointer like: &input")
	}

	o.setDefaultOn(input, defaultValue, true, false)
}

func (o *validatorS) SetDefault(input interface{}, defaultValue interface{}) {
	if reflect.TypeOf(input).Kind() != reflect.Ptr {
		panic("Please send data as a pointer like: &input")
	}

	o.setDefaultOn(input, defaultValue, false, true)
}

func (o *validatorS) setDefaultOn(input interface{}, defaultInput interface{}, onNil, onZero bool) {
	inputValue := reflect.ValueOf(input)
	defaultValue := reflect.ValueOf(defaultInput)

	for inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
	}

	for defaultValue.Kind() == reflect.Ptr {
		defaultValue = defaultValue.Elem()
	}

	if o.rules != nil {
		switch inputValue.Type().Kind() {
		case reflect.Struct:
			for fieldName, ruleSet := range o.rules {
				valueOnKeyInput := inputValue.FieldByName(fieldName)
				valueOnDefaultKeyInput := defaultValue.FieldByName(fieldName)
				if ruleSet.getName() != "" {
					fieldName = ruleSet.getName()
				}

				if !valueOnKeyInput.IsValid() {
					valueOnKeyInput = inputValue.FieldByName(fieldName)
				}
				if !valueOnDefaultKeyInput.IsValid() {
					valueOnDefaultKeyInput = defaultValue.FieldByName(fieldName)
				}
				if !valueOnKeyInput.IsValid() || !valueOnDefaultKeyInput.IsValid() {
					panic(fmt.Sprintf("value on %s is not valid", fieldName))
				}

				value := valueOnKeyInput.Interface()
				if (onNil && isNil(value)) || (onZero && !nonZeroRule(nil, value)) {
					if valueOnKeyInput.Kind() != reflect.Ptr {
						panic(fmt.Sprintf("value on %s on the side you want to copy defaults on nil(first input) has to be pointer", fieldName))
					}
					if valueOnDefaultKeyInput.Kind() != reflect.Ptr {
						valueOnDefaultKeyInput = valueOnDefaultKeyInput.Addr()
					}
					valueOnKeyInput.Set(valueOnDefaultKeyInput)
				} else {
					if ruleSet.hasDeepValidator() && (mapRule(nil, value) || structRule(nil, value) || sliceRule(nil, value)) {
						deepValidator := ruleSet.getDeepValidator()

						if onNil {
							deepValidator.SetDefaultOnNil(valueOnKeyInput.Interface(), valueOnDefaultKeyInput.Interface())
						} else if onZero {
							deepValidator.SetDefault(valueOnKeyInput.Interface(), valueOnDefaultKeyInput.Interface())
						}
					}
				}
			}
		case reflect.Map:
			for fieldName, ruleSet := range o.rules {
				valueOnKeyInput := inputValue.MapIndex(reflect.ValueOf(fieldName))
				valueOnDefaultKeyInput := defaultValue.MapIndex(reflect.ValueOf(fieldName))
				if ruleSet.getName() != "" {
					fieldName = ruleSet.getName()
				}

				if !valueOnKeyInput.IsValid() {
					valueOnKeyInput = inputValue.MapIndex(reflect.ValueOf(fieldName))
				}
				if !valueOnDefaultKeyInput.IsValid() {
					valueOnDefaultKeyInput = defaultValue.MapIndex(reflect.ValueOf(fieldName))
				}
				if !valueOnKeyInput.IsValid() || !valueOnDefaultKeyInput.IsValid() {
					panic(fmt.Sprintf("value on %s is not valid", fieldName))
				}

				value := valueOnKeyInput.Interface()
				if (onNil && isNil(value)) || (onZero && !nonZeroRule(context.TODO(), value)) {
					valueOnKeyInput.Set(valueOnDefaultKeyInput)
				} else {
					if ruleSet.hasDeepValidator() && (mapRule(context.TODO(), value) || structRule(context.TODO(), value) || sliceRule(context.TODO(), value)) {
						deepValidator := ruleSet.getDeepValidator()

						if onNil {
							deepValidator.SetDefaultOnNil(valueOnKeyInput.Interface(), valueOnDefaultKeyInput.Interface())
						} else if onZero {
							deepValidator.SetDefault(valueOnKeyInput.Interface(), valueOnDefaultKeyInput.Interface())
						}
					}
				}
			}
		default:
			return
		}
	} else if o.rule != nil {
		value := inputValue.Interface()
		if o.rule.hasDeepValidator() && (mapRule(context.TODO(), value) || structRule(context.TODO(), value) || sliceRule(context.TODO(), value)) {
			deepValidator := o.rule.getDeepValidator()

			if onNil {
				deepValidator.SetDefaultOnNil(value, defaultInput)
			} else if onZero {
				deepValidator.SetDefault(value, defaultInput)
			}
		}
	}
}

func getFieldName(path string) string {
	return strings.SplitN(path, ".", 2)[0]
}

func decryptPath(path string, v Validator, errorField playgroundValidator.FieldError) interface{} {
	var (
		splits    = strings.SplitN(path, ".", 2)
		fieldName = splits[0]
		output    = map[string]interface{}{}
	)
	if len(splits) == 2 {
		path = splits[1]
	} else {
		path = ""
	}
	if path == "" {
		if rs := v.getRules(); len(rs) != 0 {
			if r, ok := rs[fieldName]; ok {
				if r.getName() != "" {
					fieldName = r.getName()
				}
				message := getRawErrorMessage(errorField.Tag(), *v.getMessages(), r.getSpecificMessages(), defaultValidatorErrorMessages)
				// Do not add translator, translator is for Validate process
				message = getFormattedErrorMessage(message, fieldName, errorField.Value(), r.getOption(errorField.Tag()))
				return message
			} else {
				panic("error structure does not match with validator structure")
			}
		} else {
			panic("error structure does not match with validator structure")
		}
	} else if rs := v.getRules(); len(rs) != 0 {
		re, _ := regexp.Compile(`(\w+)\[(\d+)\]`)
		slicePieces := re.FindStringSubmatch(fieldName)
		arrayItem := -1
		if len(slicePieces) > 0 {
			fieldName = slicePieces[1]
			arrayItem, _ = strconv.Atoi(slicePieces[2])
		}
		if r, ok := rs[fieldName]; ok {
			if deep := r.getDeepValidator(); deep != nil {
				if r.getName() != "" {
					fieldName = r.getName()
				}
				output[fieldName] = decryptPath(path, deep, errorField)
			} else if children := r.getChildrenValidator(); arrayItem != -1 && children != nil {
				if splits := strings.Split(path, "."); len(splits) == 1 && len(re.FindStringSubmatch(splits[0])) == 0 {
					output[fieldName] = map[int]interface{}{arrayItem: map[string]interface{}{splits[0]: decryptPath(path, children, errorField)}}
				} else {
					output[fieldName] = map[int]interface{}{arrayItem: decryptPath(path, children, errorField)}
				}
			} else {
				panic("error structure does not match with validator structure")
			}
		} else {
			panic("error structure does not match with validator structure")
		}
	} else {
		panic("error structure does not match with validator structure")
	}

	return output
}

func decryptErrors(err error, v Validator, unmarshalError bool) interface{} {
	output := map[string]interface{}{}
	if e, ok := err.(playgroundValidator.ValidationErrors); ok {
		for i := 0; i < len(e); i++ {
			errorField := e[i]
			splits := strings.SplitN(errorField.StructNamespace(), ".", 2)
			path := splits[1]
			out := decryptPath(path, v, errorField)
			if outMap, ok := out.(map[string]interface{}); ok && len(outMap) != 0 {
				for key, value := range outMap {
					output[key] = value
				}
			} else if outString, ok := out.(string); ok {
				name := getFieldName(path)
				if r := v.getRules()[name]; r != nil && r.getName() != "" {
					name = r.getName()
				}
				output[name] = outString
			}
		}
	} else if e, ok := err.(sliceValidationError); ok {
		for i := 0; i < len(e); i++ {
			if r := v.getRule(); r != nil {
				if deep := r.getDeepValidator(); deep != nil {
					if out := decryptErrors(e[i], deep, unmarshalError); out != nil {
						output[strconv.Itoa(i)] = out
					}
				} else {
					panic("error structure does not match with validator structure")
				}
			} else {
				panic("error structure does not match with validator structure")
			}
		}
	} else if e, ok := err.(*json.UnmarshalTypeError); ok {
		if unmarshalError {
			return e.Error()
		} else {
			return UnmarshalError
		}
	} else {
		panic("passed error is not supported, errors returned from methods like BindJson(in gin framework) is supported")
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

func (o *validatorS) setMessages(messages *Messages) {
	o.messages = messages
}
