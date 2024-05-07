package tests

import (
	"context"
	"testing"
)

func TestNonZero(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "pass-int",
			validator: g.Validator(g.R().NonZero()),
			in:        1,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-float",
			validator: g.Validator(g.R().NonZero()),
			in:        1.1,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-string",
			validator: g.Validator(g.R().NonZero()),
			in:        "1",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-slice",
			validator: g.Validator(g.R().NonZero()),
			in:        []int{},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-map",
			validator: g.Validator(g.R().NonZero()),
			in:        map[int]int{},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "fail-int",
			validator: g.Validator(g.R().NonZero()),
			in:        0,
			panic:     false,
			expected:  []string{"can not be 0"},
		},
		{
			name:      "fail-float",
			validator: g.Validator(g.R().NonZero()),
			in:        0.0,
			panic:     false,
			expected:  []string{"can not be 0"},
		},
		{
			name:      "fail-string",
			validator: g.Validator(g.R().NonZero()),
			in:        "",
			panic:     false,
			expected:  []string{"can not be 0"},
		},
		{
			name:      "fail-nil",
			validator: g.Validator(g.R().NonZero()),
			in:        nil,
			panic:     false,
			expected:  nil,
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
