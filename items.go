package galidator

import "strings"

type Items map[string]item

type itemS struct {
	validates map[string]func(interface{}) bool
}

// Adds int validator
func (o *itemS) Int() item {
	functionName := strings.ToLower(getFunctionName(Int))
	o.validates[functionName] = Int
	return o
}

// Adds float validator
func (o *itemS) Float() item {
	functionName := strings.ToLower("float")
	o.validates[functionName] = Float
	return o
}

// Validates all validators that did set for the field
func (o *itemS) validate(input interface{}) []string {
	fails := []string{}
	for key, vFunction := range o.validates {
		if !vFunction(input) {
			fails = append(fails, key)
		}
	}

	return fails
}

type item interface {
	validate(interface{}) []string
	Int() item
	Float() item
}
