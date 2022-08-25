package galidator

import (
	"github.com/golodash/godash/strings"
)

type (
	// A struct to implement generator interface
	generatorS struct{}

	// An interface to generate a validator or ruleSet
	generator interface {
		// Generates a validator interface which can be used to validate some data by some rules
		//
		// Please use CapitalCase for rules' keys (Important for getting data out of struct types)
		Validator(rules Rules, messages ...Messages) validator
		// Generates a ruleSet to validate passed information
		RuleSet() ruleSet
		// An alias for RuleSet function
		R() ruleSet
	}
)

func (o *generatorS) Validator(rules Rules, errorMessages ...Messages) validator {
	var messages Messages = nil
	if len(errorMessages) != 0 {
		messages = errorMessages[0]
	}

	for key, message := range messages {
		updatedKey := strings.SnakeCase(key)
		if updatedKey == key {
			continue
		}

		messages[updatedKey] = message
		delete(messages, key)
	}

	return &validatorS{rules: rules, messages: messages, specificMessages: SpecificMessages{}}
}

func (o *generatorS) RuleSet() ruleSet {
	return &ruleSetS{validators: Validators{}, options: options{}, isOptional: true, deepValidator: nil, childrenRule: nil}
}

func (o *generatorS) R() ruleSet {
	return o.RuleSet()
}

// A unique instance of generatorS to stop creating unnecessarily multiple instances of a generator
var generatorInstance = &generatorS{}

// Returns a Validator Generator
func New() generator {
	return generatorInstance
}
