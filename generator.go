package galidator

type generatorS struct{}

func (o *generatorS) Generate(items Items, errorMessages Messages) validator {
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
