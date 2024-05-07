package tests

import (
	"context"
	"testing"

	"github.com/golodash/galidator"
)

func TestOR(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "fail-1",
			validator: g.Validator(g.R().OR(g.R().String().Min(2), g.R().Int().Max(3)), galidator.Messages{"or": "or failed"}),
			in:        "2",
			panic:     false,
			expected:  []string{"or failed"},
		},
		{
			name:      "fail-2",
			validator: g.Validator(g.R().OR(g.R().String().Min(2), g.R().Int().Max(3)), galidator.Messages{"or": "or failed"}),
			in:        4,
			panic:     false,
			expected:  []string{"or failed"},
		},
		{
			name:      "fail-3",
			validator: g.Validator(g.R().OR(g.R().String().Min(2), g.R().Int().Max(3)), galidator.Messages{"or": "or failed"}),
			in:        []string{"s"},
			panic:     false,
			expected:  []string{"or failed"},
		},
		{
			name:      "pass-1",
			validator: g.Validator(g.R().OR(g.R().String().Min(2), g.R().Int().Max(3)), galidator.Messages{"or": "or failed"}),
			in:        "111",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-2",
			validator: g.Validator(g.R().OR(g.R().String().Min(2), g.R().Int().Max(3)), galidator.Messages{"or": "or failed"}),
			in:        2,
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
