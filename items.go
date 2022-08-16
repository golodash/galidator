package galidator

import (
	"fmt"
)

type option map[string]string
type options map[string]option
type validators map[string]func(interface{}) bool
type itemS struct {
	validators validators
	options    options
}

// Adds int validator
func (o *itemS) Int() item {
	o.validators["int"] = Int
	return o
}

// Adds float validator
func (o *itemS) Float() item {
	o.validators["float"] = Float
	return o
}

func (o *itemS) Min(min float64) item {
	functionName := "min"
	o.validators[functionName] = Min(min)
	precision := determinePrecision(min)
	o.AddOption(functionName, "min", fmt.Sprintf("%."+precision+"f", min))
	return o
}

func (o *itemS) Max(max float64) item {
	functionName := "max"
	o.validators[functionName] = Max(max)
	precision := determinePrecision(max)
	o.AddOption(functionName, "max", fmt.Sprintf("%."+precision+"f", max))
	return o
}

func (o *itemS) Len(from, to int) item {
	functionName := "len"
	o.validators[functionName] = Len(from, to)
	o.AddOption(functionName, "from", fmt.Sprintf("%d", from))
	o.AddOption(functionName, "to", fmt.Sprintf("%d", to))
	return o
}

func (o *itemS) Required() item {
	functionName := "required"
	o.validators[functionName] = Required
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

func (o *itemS) GetOption(key string) option {
	if option, ok := o.options[key]; ok {
		return option
	}
	return option{}
}

func (o *itemS) AddOption(key string, subKey string, value string) {
	if option, ok := o.options[key]; ok {
		option[subKey] = value
		return
	}
	o.options[key] = option{subKey: value}
}

type item interface {
	validate(interface{}) []string

	Int() item
	Float() item
	Min(min float64) item
	Max(max float64) item
	Len(from int, to int) item
	Required() item

	GetOption(key string) option
	AddOption(key string, subKey string, value string)
}
