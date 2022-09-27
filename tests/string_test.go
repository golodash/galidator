package tests

import (
	"testing"
)

func TestString(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "fail-int",
			validator: g.Validator(g.R().String()),
			in:        1,
			panic:     false,
			expected:  []string{"not a string"},
		},
		{
			name:      "fail-float",
			validator: g.Validator(g.R().String()),
			in:        1.3,
			panic:     false,
			expected:  []string{"not a string"},
		},
		{
			name:      "fail-slice",
			validator: g.Validator(g.R().String()),
			in:        []int{1},
			panic:     false,
			expected:  []string{"not a string"},
		},
		{
			name:      "fail-map",
			validator: g.Validator(g.R().String()),
			in:        map[int]int{1: 3},
			panic:     false,
			expected:  []string{"not a string"},
		},
		{
			name:      "pass",
			validator: g.Validator(g.R().String()),
			in:        "1",
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
