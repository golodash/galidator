package galidator

import (
	"github.com/golodash/godash/strings"
)

type (
	// A struct to implement generator interface
	generatorS struct{}

	// An interface to generate a validator or rule
	generator interface {
		// Generates a validator interface which can be used to validate some data by some filters
		Validator(Rules, Messages) validator
		// Generates a rule to validate passed information
		Rule() rule
	}
)

func (o *generatorS) Validator(rules Rules, errorMessages Messages) validator {
	for key, errorMessage := range errorMessages {
		updatedKey := strings.SnakeCase(key)
		if updatedKey == key {
			continue
		}

		errorMessages[updatedKey] = errorMessage
		delete(errorMessages, key)
	}

	return &validatorS{rules: rules, messages: errorMessages}
}

func (o *generatorS) Rule() rule {
	return &ruleS{validators: validators{}, options: options{}}
}

// A unique instance of generatorS to stop creating unnecessarily multiple instances of a generator
var generatorInstance = &generatorS{}

// Returns a Generator
func New() generator {
	return generatorInstance
}
