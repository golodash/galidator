package galidator

import (
	"github.com/golodash/godash/strings"
)

type (
	generatorS struct{}
	generator  interface {
		Generate(Rules, Messages) validator
		Rule() rule
	}
)

func (o *generatorS) Generate(rules Rules, errorMessages Messages) validator {
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

var generatorInstance = &generatorS{}

// Returns a Generator
func New() generator {
	return generatorInstance
}
