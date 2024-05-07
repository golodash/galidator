package tests

import (
	"context"
	"testing"

	"github.com/golodash/galidator"
)

func TestXOR(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "fail-both_fail",
			validator: g.Validator(g.R().XOR(g.R().Min(2), g.R().Max(1).String()), galidator.Messages{"xor": "xor failed"}),
			in:        1,
			panic:     false,
			expected:  []string{"xor failed"},
		},
		{
			name:      "fail-both_succeed",
			validator: g.Validator(g.R().XOR(g.R().Min(2), g.R().Min(1).String()), galidator.Messages{"xor": "xor failed"}),
			in:        "22",
			panic:     false,
			expected:  []string{"xor failed"},
		},
		{
			name:      "pass-1",
			validator: g.Validator(g.R().XOR(g.R().Min(2), g.R().Min(1).String()), galidator.Messages{"xor": "xor failed"}),
			in:        35,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-2",
			validator: g.Validator(g.R().XOR(g.R().Min(2), g.R().Min(1).String()), galidator.Messages{"xor": "xor failed"}),
			in:        "3",
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
