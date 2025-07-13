package tests

import (
	"context"
	"testing"
)

func TestRequired(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "pass-int",
			validator: g.Validator(g.R().Required()),
			in:        1,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-float",
			validator: g.Validator(g.R().Required()),
			in:        1.1,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-string",
			validator: g.Validator(g.R().Required()),
			in:        "1",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-slice",
			validator: g.Validator(g.R().Required()),
			in:        []int{1},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-map",
			validator: g.Validator(g.R().Required()),
			in:        map[int]int{1: 1},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-int",
			validator: g.Validator(g.R().Required()),
			in:        0,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-float",
			validator: g.Validator(g.R().Required()),
			in:        0.0,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "fail-string",
			validator: g.Validator(g.R().Required()),
			in:        "",
			panic:     false,
			expected:  []string{"required"},
		},
		{
			name:      "fail-slice",
			validator: g.Validator(g.R().Required()),
			in:        []int{},
			panic:     false,
			expected:  []string{"required"},
		},
		{
			name:      "fail-map",
			validator: g.Validator(g.R().Required()),
			in:        map[string]string{},
			panic:     false,
			expected:  []string{"required"},
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
