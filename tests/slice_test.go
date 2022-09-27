package tests

import (
	"testing"
)

func TestSlice(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "fail-int",
			validator: g.Validator(g.R().Slice()),
			in:        1,
			panic:     false,
			expected:  []string{"not a slice"},
		},
		{
			name:      "fail-float",
			validator: g.Validator(g.R().Slice()),
			in:        1.3,
			panic:     false,
			expected:  []string{"not a slice"},
		},
		{
			name:      "fail-map",
			validator: g.Validator(g.R().Slice()),
			in:        map[string]int{"1": 1},
			panic:     false,
			expected:  []string{"not a slice"},
		},
		{
			name:      "pass",
			validator: g.Validator(g.R().Slice()),
			in:        []int{1},
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
