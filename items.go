package galidator

import (
	"fmt"

	"github.com/golodash/galidator/independents"
)

type (
	option     map[string]string
	options    map[string]option
	validators map[string]func(interface{}) bool
	itemS      struct {
		validators validators
		options    options
	}
	item interface {
		validate(interface{}) []string

		Int() item
		Float() item
		Min(min float64) item
		Max(max float64) item
		Len(from int, to int) item
		Required() item

		getOption(key string) option
		addOption(key string, subKey string, value string)
	}
)

// Adds int validator
func (o *itemS) Int() item {
	o.validators["int"] = independents.Int
	return o
}

// Adds float validator
func (o *itemS) Float() item {
	o.validators["float"] = independents.Float
	return o
}

func (o *itemS) Min(min float64) item {
	functionName := "min"
	o.validators[functionName] = independents.Min(min)
	precision := determinePrecision(min)
	o.addOption(functionName, "min", fmt.Sprintf("%."+precision+"f", min))
	return o
}

func (o *itemS) Max(max float64) item {
	functionName := "max"
	o.validators[functionName] = independents.Max(max)
	precision := determinePrecision(max)
	o.addOption(functionName, "max", fmt.Sprintf("%."+precision+"f", max))
	return o
}

func (o *itemS) Len(from, to int) item {
	functionName := "len"
	o.validators[functionName] = independents.Len(from, to)
	o.addOption(functionName, "from", fmt.Sprintf("%d", from))
	o.addOption(functionName, "to", fmt.Sprintf("%d", to))
	return o
}

func (o *itemS) Required() item {
	functionName := "required"
	o.validators[functionName] = independents.Required
	return o
}

// Validates all validators that did set for the field
func (o *itemS) validate(input interface{}) []string {
	fails := []string{}
	for key, vFunction := range o.validators {
		if !vFunction(input) {
			fails = append(fails, key)
		}
	}

	return fails
}

func (o *itemS) getOption(key string) option {
	if option, ok := o.options[key]; ok {
		return option
	}
	return option{}
}

func (o *itemS) addOption(key string, subKey string, value string) {
	if option, ok := o.options[key]; ok {
		option[subKey] = value
		return
	}
	o.options[key] = option{subKey: value}
}
