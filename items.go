package galidator

import (
	"fmt"
	"strings"
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
	functionName := strings.ToLower("int")
	o.validators[functionName] = Int
	return o
}

// Adds float validator
func (o *itemS) Float() item {
	functionName := strings.ToLower("float")
	o.validators[functionName] = Float
	return o
}

func (o *itemS) Min(min float64) item {
	functionName := strings.ToLower("min")
	o.validators[functionName] = Min(min)
	precision := determinePrecision(min)
	o.AddOption("min", "min", fmt.Sprintf("%."+precision+"f", min))
	return o
}

func (o *itemS) Max(max float64) item {
	functionName := strings.ToLower("max")
	o.validators[functionName] = Max(max)
	precision := determinePrecision(max)
	o.AddOption("max", "max", fmt.Sprintf("%."+precision+"f", max))
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
	Min(float64) item
	Max(float64) item

	GetOption(key string) option
	AddOption(key string, subKey string, value string)
}
