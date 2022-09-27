package tests

import (
	"testing"
)

func TestEmail(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "pass-1",
			validator: g.Validator(g.R().Email()),
			in:        "m@g.c",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-2",
			validator: g.Validator(g.R().Email()),
			in:        "m@g",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "fail-1",
			validator: g.Validator(g.R().Email()),
			in:        "m@",
			panic:     false,
			expected:  []string{"not a valid email address"},
		},
		{
			name:      "fail-2",
			validator: g.Validator(g.R().Email()),
			in:        "m",
			panic:     false,
			expected:  []string{"not a valid email address"},
		},
		{
			name:      "fail-3",
			validator: g.Validator(g.R().Email()),
			in:        "m@g.",
			panic:     false,
			expected:  []string{"not a valid email address"},
		},
		{
			name:      "fail-4",
			validator: g.Validator(g.R().Email()),
			in:        "m.",
			panic:     false,
			expected:  []string{"not a valid email address"},
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
