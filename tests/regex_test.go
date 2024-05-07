package tests

import (
	"context"
	"testing"

	"github.com/golodash/galidator/v2"
)

func TestRegex(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "pass",
			validator: g.Validator(g.R().Regex("^123$")),
			in:        "123",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "fail",
			validator: g.Validator(g.R().Regex("&123$").SpecificMessages(galidator.Messages{"regex": "regex failed"})),
			in:        "1233",
			panic:     false,
			expected:  []string{"regex failed"},
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
