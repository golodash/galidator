package tests

import (
	"testing"

	"github.com/golodash/galidator"
	"github.com/golodash/godash/generals"
)

type scenario struct {
	name      string
	validator galidator.Validator
	in        interface{}
	panic     bool
	expected  interface{}
}

var g = galidator.New()

// Used in test cases to prevent code breaking
func deferTestCases(t *testing.T, crash bool, expected interface{}) {
	err := recover()

	if err != nil && !crash {
		t.Errorf("expected = %v, err = %s", expected, err)
	}
}

// Checks if `expected` and `output` are the same.
//
// If they are the same, returns true.
//
// If not, raises an error in `t` and returns false.
func check(t *testing.T, expected, output interface{}) bool {
	if !generals.Same(output, expected) {
		t.Errorf("expected = %v, output = %v", expected, output)
		return false
	}
	return true
}
