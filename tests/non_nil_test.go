package tests

import (
	"context"
	"testing"
)

func TestNonNil(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "pass-int",
			validator: g.Validator(g.R().NonNil()),
			in:        1,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-float",
			validator: g.Validator(g.R().NonNil()),
			in:        1.1,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-string",
			validator: g.Validator(g.R().NonNil()),
			in:        "1",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-slice",
			validator: g.Validator(g.R().NonNil()),
			in:        []int{},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-map",
			validator: g.Validator(g.R().NonNil()),
			in:        map[int]int{},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "fail-nil",
			validator: g.Validator(g.R().NonNil()),
			in:        nil,
			panic:     false,
			expected:  []string{"can not be nil"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			defer deferTestCases(t, s.panic, s.expected)

			output := s.validator.Validate(context.TODO(), s.in)
			check(t, s.expected, output)
		})
	}
}
