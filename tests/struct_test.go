package tests

import (
	"testing"
)

type structTest struct{}

func TestStruct(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "fail-int",
			validator: g.Validator(g.R().Struct()),
			in:        1,
			panic:     false,
			expected:  []string{"not a struct"},
		},
		{
			name:      "fail-float",
			validator: g.Validator(g.R().Struct()),
			in:        1.3,
			panic:     false,
			expected:  []string{"not a struct"},
		},
		{
			name:      "fail-slice",
			validator: g.Validator(g.R().Struct()),
			in:        []int{1},
			panic:     false,
			expected:  []string{"not a struct"},
		},
		{
			name:      "pass",
			validator: g.Validator(g.R().Struct()),
			in:        structTest{},
			panic:     false,
			expected:  nil,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			defer deferTestCases(t, s.panic, s.expected)

			output := s.validator.Validate(s.in)
			check(t, s.expected, output)
		})
	}
}
