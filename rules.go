package galidator

import (
	"fmt"

	"github.com/golodash/galidator/filters"
)

type (
	// A map with data recorded to be used in returning error messages
	option map[string]string

	// A map of different rules as key with their own option
	options map[string]option

	// A map full of validators which is assigned for a single key in a validator struct
	validators map[string]func(interface{}) bool

	// A struct to implement rule interface
	ruleS struct {
		// Used to validate user's data
		validators validators
		// Used in returning error messages
		options options
	}

	// An interface with some functions to satisfy validation purpose
	rule interface {
		// Validates all validators defined
		validate(interface{}) []string

		// Checks if input (can be)/is int
		Int() rule
		// Checks if input (can be)/is float
		Float() rule
		// Checks if input acts like: input >= min or len(input) >= min
		Min(min float64) rule
		// Checks if input acts like: input <= max or len(input) <= max
		Max(max float64) rule
		// Checks if input acts like: len(input) >= from && len(input) <= to
		//
		// If from == -1, no check on from will happen
		// If to == -1, no check on to will happen
		LenRange(from int, to int) rule
		// Checks if input acts like: len(input) == length
		Len(length int) rule
		// Checks if input is not zero(0, "", '') or nil or empty
		Required() rule
		// Checks if input is not zero(0, "", '')
		NonZero() rule
		// Checks if input is not nil
		NonNil() rule
		// Checks if input has items inside it
		NonEmpty() rule
		// Returns option of the passed rule key
		getOption(key string) option
		// Adds a new subKey with a value associated with it to option of passed rule key
		addOption(key string, subKey string, value string)
	}
)

func (o *ruleS) Int() rule {
	o.validators["int"] = filters.Int
	return o
}

func (o *ruleS) Float() rule {
	o.validators["float"] = filters.Float
	return o
}

func (o *ruleS) Min(min float64) rule {
	functionName := "min"
	o.validators[functionName] = filters.Min(min)
	precision := determinePrecision(min)
	o.addOption(functionName, "min", fmt.Sprintf("%."+precision+"f", min))
	return o
}

func (o *ruleS) Max(max float64) rule {
	functionName := "max"
	o.validators[functionName] = filters.Max(max)
	precision := determinePrecision(max)
	o.addOption(functionName, "max", fmt.Sprintf("%."+precision+"f", max))
	return o
}

func (o *ruleS) LenRange(from, to int) rule {
	functionName := "len_range"
	o.validators[functionName] = filters.LenRange(from, to)
	o.addOption(functionName, "from", fmt.Sprintf("%d", from))
	o.addOption(functionName, "to", fmt.Sprintf("%d", to))
	return o
}

func (o *ruleS) Len(length int) rule {
	functionName := "len"
	o.validators[functionName] = filters.Len(length)
	o.addOption(functionName, "length", fmt.Sprint(length))
	return o
}

func (o *ruleS) Required() rule {
	functionName := "required"
	o.validators[functionName] = filters.Required
	return o
}

func (o *ruleS) NonZero() rule {
	functionName := "non_zero"
	o.validators[functionName] = filters.NonZero
	return o
}

func (o *ruleS) NonNil() rule {
	functionName := "non_nil"
	o.validators[functionName] = filters.NonNil
	return o
}

func (o *ruleS) NonEmpty() rule {
	functionName := "non_empty"
	o.validators[functionName] = filters.NonEmpty
	return o
}

func (o *ruleS) validate(input interface{}) []string {
	fails := []string{}
	for key, vFunction := range o.validators {
		if !vFunction(input) {
			fails = append(fails, key)
		}
	}

	return fails
}

func (o *ruleS) getOption(key string) option {
	if option, ok := o.options[key]; ok {
		return option
	}
	return option{}
}

func (o *ruleS) addOption(key string, subKey string, value string) {
	if option, ok := o.options[key]; ok {
		option[subKey] = value
		return
	}
	o.options[key] = option{subKey: value}
}
