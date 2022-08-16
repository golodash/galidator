package galidator

import (
	"fmt"

	"github.com/golodash/galidator/independents"
)

type (
	option     map[string]string
	options    map[string]option
	validators map[string]func(interface{}) bool
	ruleS      struct {
		validators validators
		options    options
	}
	rule interface {
		validate(interface{}) []string

		Int() rule
		Float() rule
		Min(min float64) rule
		Max(max float64) rule
		Len(from int, to int) rule
		Required() rule

		getOption(key string) option
		addOption(key string, subKey string, value string)
	}
)

// Adds int validator
func (o *ruleS) Int() rule {
	o.validators["int"] = independents.Int
	return o
}

// Adds float validator
func (o *ruleS) Float() rule {
	o.validators["float"] = independents.Float
	return o
}

func (o *ruleS) Min(min float64) rule {
	functionName := "min"
	o.validators[functionName] = independents.Min(min)
	precision := determinePrecision(min)
	o.addOption(functionName, "min", fmt.Sprintf("%."+precision+"f", min))
	return o
}

func (o *ruleS) Max(max float64) rule {
	functionName := "max"
	o.validators[functionName] = independents.Max(max)
	precision := determinePrecision(max)
	o.addOption(functionName, "max", fmt.Sprintf("%."+precision+"f", max))
	return o
}

func (o *ruleS) Len(from, to int) rule {
	functionName := "len"
	o.validators[functionName] = independents.Len(from, to)
	o.addOption(functionName, "from", fmt.Sprintf("%d", from))
	o.addOption(functionName, "to", fmt.Sprintf("%d", to))
	return o
}

func (o *ruleS) Required() rule {
	functionName := "required"
	o.validators[functionName] = independents.Required
	return o
}

// Validates all validators that did set for the field
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
