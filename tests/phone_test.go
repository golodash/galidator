package tests

import (
	"context"
	"testing"

	"github.com/golodash/galidator"
)

func TestPhone(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "pass-1",
			validator: g.Validator(g.R().Phone()),
			in:        "+989123456789",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-2",
			validator: g.Validator(g.R().Phone()),
			in:        "+989127769381",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-3",
			validator: g.Validator(g.R().Phone()),
			in:        "+989127769383",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "fail-1",
			validator: g.Validator(g.R().Phone().SpecificMessages(galidator.Messages{"phone": "phone failed"})),
			in:        "09101234567",
			panic:     false,
			expected:  []string{"phone failed"},
		},
		{
			name:      "fail-2",
			validator: g.Validator(g.R().Phone().SpecificMessages(galidator.Messages{"phone": "phone failed"})),
			in:        "+989123456",
			panic:     false,
			expected:  []string{"phone failed"},
		},
		{
			name:      "fail-3",
			validator: g.Validator(g.R().Phone().SpecificMessages(galidator.Messages{"phone": "phone failed"})),
			in:        "091012345",
			panic:     false,
			expected:  []string{"phone failed"},
		},
		{
			name:      "fail-3",
			validator: g.Validator(g.R().Phone().SpecificMessages(galidator.Messages{"phone": "phone failed"})),
			in:        "+989 10 015 4789",
			panic:     false,
			expected:  []string{"phone failed"},
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
