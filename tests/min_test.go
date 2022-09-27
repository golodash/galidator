package tests

import (
	"testing"

	"github.com/golodash/galidator"
)

func TestMin(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "pass",
			validator: g.Validator(g.R().Min(3).SpecificMessages(galidator.Messages{"min": "has to be at least $min characters long"})),
			in:        "111",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "fail-string",
			validator: g.Validator(g.R().Min(3).SpecificMessages(galidator.Messages{"min": "has to be at least $min characters long"})),
			in:        "ab",
			panic:     false,
			expected:  []string{"has to be at least 3 characters long"},
		},
		{
			name:      "fail-slice",
			validator: g.Validator(g.R().Min(3).SpecificMessages(galidator.Messages{"min": "has to be at least $min characters long"})),
			in:        []int{1, 2},
			panic:     false,
			expected:  []string{"has to be at least 3 characters long"},
		},
		{
			name:      "fail-map",
			validator: g.Validator(g.R().Min(3).SpecificMessages(galidator.Messages{"min": "has to be at least $min characters long"})),
			in:        map[string]int{"1": 1, "2": 2},
			panic:     false,
			expected:  []string{"has to be at least 3 characters long"},
		},
		{
			name:      "pass-string",
			validator: g.Validator(g.R().Min(3).SpecificMessages(galidator.Messages{"min": "has to be at least $min characters long"})),
			in:        "111",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-slice",
			validator: g.Validator(g.R().Min(3).SpecificMessages(galidator.Messages{"min": "has to be at least $min characters long"})),
			in:        []int{1, 2, 3},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-map",
			validator: g.Validator(g.R().Min(3).SpecificMessages(galidator.Messages{"min": "has to be at least $min characters long"})),
			in:        map[string]int{"1": 1, "2": 2, "3": 3},
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