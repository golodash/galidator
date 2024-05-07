package tests

import (
	"context"
	"testing"

	"github.com/golodash/galidator"
)

func TestChoices(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "fail-1",
			validator: g.Validator(g.R().Choices("1", 2, 3.0), galidator.Messages{"choices": "choices failed"}),
			in:        "2",
			panic:     false,
			expected:  []string{"choices failed"},
		},
		{
			name:      "fail-2",
			validator: g.Validator(g.R().Choices("1", 2, 3.0), galidator.Messages{"choices": "choices failed"}),
			in:        1,
			panic:     false,
			expected:  []string{"choices failed"},
		},
		{
			name:      "fail-3",
			validator: g.Validator(g.R().Choices("1", 2, 3.0), galidator.Messages{"choices": "choices failed"}),
			in:        3,
			panic:     false,
			expected:  []string{"choices failed"},
		},
		{
			name:      "fail-4",
			validator: g.Validator(g.R().Choices("1", 2, 3.0), galidator.Messages{"choices": "choices failed"}),
			in:        map[string]string{"s": "s"},
			panic:     false,
			expected:  []string{"choices failed"},
		},
		{
			name:      "pass-1",
			validator: g.Validator(g.R().Choices("1", 2, 3.0), galidator.Messages{"choices": "choices failed"}),
			in:        "1",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-2",
			validator: g.Validator(g.R().Choices("1", 2, 3.0), galidator.Messages{"choices": "choices failed"}),
			in:        2,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-3",
			validator: g.Validator(g.R().Choices("1", 2, 3.0), galidator.Messages{"choices": "choices failed"}),
			in:        3.0,
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
