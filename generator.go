package galidator

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
		//
		// Variable name will be used in output of error keys
		RuleSet(name ...string) ruleSet
		// An alias for RuleSet function
		R(name ...string) ruleSet
	}
)

func (o *generatorS) Validator(rules Rules, errorMessages ...Messages) validator {
	var messages Messages = nil
	if len(errorMessages) != 0 {
		messages = errorMessages[0]
	}

	for key, message := range messages {
		messages[key] = message
		delete(messages, key)
	}

	return &validatorS{rules: rules, messages: messages, specificMessages: SpecificMessages{}}
}

func (o *generatorS) RuleSet(name ...string) ruleSet {
	var output = ""
	if len(name) != 0 {
		output = name[0]
	}
	return &ruleSetS{name: output, validators: Validators{}, requires: requires{}, options: options{}, isOptional: true, deepValidator: nil, childrenRule: nil}
}

func (o *generatorS) R(name ...string) ruleSet {
	return o.RuleSet(name...)
}

// A unique instance of generatorS to stop creating unnecessarily multiple instances of a generator
var generatorInstance = &generatorS{}

// Returns a Validator Generator
func New() generator {
	return generatorInstance
}
