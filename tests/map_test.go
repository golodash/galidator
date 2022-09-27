package tests

import (
	"testing"
)

func TestMap(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "fail-int",
			validator: g.Validator(g.R().Map()),
			in:        1,
			panic:     false,
			expected:  []string{"not a map"},
		},
		{
			name:      "fail-float",
			validator: g.Validator(g.R().Map()),
			in:        1.3,
			panic:     false,
			expected:  []string{"not a map"},
		},
		{
			name:      "fail-slice",
			validator: g.Validator(g.R().Map()),
			in:        []int{1},
			panic:     false,
			expected:  []string{"not a map"},
		},
		{
			name:      "pass",
			validator: g.Validator(g.R().Map()),
			in:        map[string]string{"1": "1"},
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
