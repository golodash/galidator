package tests

import (
	"context"
	"testing"
)

func TestCustomValidatorsFromGenerators(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "fail-1",
			validator: g.Validator(g.R().RegisteredCustom("custom_validator")),
			in:        "some text",
			panic:     false,
			expected:  []string{"custom error"},
		},
		{
			name:      "fail-2",
			validator: g.Validator(g.R().RegisteredCustom("second_custom_validator")),
			in:        "some text",
			panic:     false,
			expected:  []string{"second custom error"},
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
