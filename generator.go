package galidator

import (
	"github.com/golodash/godash/strings"
)

type generatorS struct{}

func (o *generatorS) Generate(items Items, errorMessages Messages) validator {
	for key, errorMessage := range errorMessages {
		updatedKey := strings.SnakeCase(key)
		if updatedKey == key {
			continue
		}

		errorMessages[updatedKey] = errorMessage
		delete(errorMessages, key)
	}

	return &validatorS{items: items, messages: errorMessages}
}

func (o *generatorS) Item() item {
	return &itemS{validators: validators{}, options: options{}}
}

type generator interface {
	Generate(Items, Messages) validator
	Item() item
}

// Returns a Generator

var generatorValue = &generatorS{}

func New() generator {
	return generatorValue
}
