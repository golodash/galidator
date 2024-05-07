package tests

import (
	"context"
	"testing"
)

func TestNonEmpty(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "pass-int",
			validator: g.Validator(g.R().NonEmpty()),
			in:        0,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-float",
			validator: g.Validator(g.R().NonEmpty()),
			in:        0,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-string",
			validator: g.Validator(g.R().NonEmpty()),
			in:        "",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-slice",
			validator: g.Validator(g.R().NonEmpty()),
			in:        []int{1},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-map",
			validator: g.Validator(g.R().NonEmpty()),
			in:        map[int]int{1: 1},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-nil",
			validator: g.Validator(g.R().NonEmpty()),
			in:        nil,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "fail-slice",
			validator: g.Validator(g.R().NonEmpty()),
			in:        []int{},
			panic:     false,
			expected:  []string{"can not be empty"},
		},
		{
			name:      "fail-map",
			validator: g.Validator(g.R().NonEmpty()),
			in:        map[int]int{},
			panic:     false,
			expected:  []string{"can not be empty"},
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
