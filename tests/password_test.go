package tests

import (
	"context"
	"testing"

	"github.com/golodash/galidator/v2"
)

func TestPassword(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "fail-1",
			validator: g.Validator(g.R().Password().SpecificMessages(galidator.Messages{"password": "password failed"})),
			in:        "1234567891011121314151617181920",
			panic:     false,
			expected:  []string{"password failed"},
		},
		{
			name:      "fail-2",
			validator: g.Validator(g.R().Password().SpecificMessages(galidator.Messages{"password": "password failed"})),
			in:        "123456789M",
			panic:     false,
			expected:  []string{"password failed"},
		},
		{
			name:      "fail-3",
			validator: g.Validator(g.R().Password().SpecificMessages(galidator.Messages{"password": "password failed"})),
			in:        "123456789Mh",
			panic:     false,
			expected:  []string{"password failed"},
		},
		{
			name:      "fail-4",
			validator: g.Validator(g.R().Password().SpecificMessages(galidator.Messages{"password": "password failed"})),
			in:        "123Ms!",
			panic:     false,
			expected:  []string{"password failed"},
		},
		{
			name:      "pass",
			validator: g.Validator(g.R().Password().SpecificMessages(galidator.Messages{"password": "password failed"})),
			in:        "123456789Mh!",
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
