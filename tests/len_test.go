package tests

import (
	"context"
	"testing"

	"github.com/golodash/galidator"
)

func TestLen(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "pass",
			validator: g.Validator(g.R().Len(3).SpecificMessages(galidator.Messages{"len": "must be exactly $length characters long"})),
			in:        "111",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "fail",
			validator: g.Validator(g.R().Len(3).SpecificMessages(galidator.Messages{"len": "must be exactly $length characters long"})),
			in:        12,
			panic:     false,
			expected:  []string{"must be exactly 3 characters long"},
		},
		{
			name:      "fail-string",
			validator: g.Validator(g.R().Len(3).SpecificMessages(galidator.Messages{"len": "must be exactly $length characters long"})),
			in:        "11",
			panic:     false,
			expected:  []string{"must be exactly 3 characters long"},
		},
		{
			name:      "fail-slice",
			validator: g.Validator(g.R().Len(3).SpecificMessages(galidator.Messages{"len": "must be exactly $length characters long"})),
			in:        []int{1, 2, 3, 4, 5, 6},
			panic:     false,
			expected:  []string{"must be exactly 3 characters long"},
		},
		{
			name:      "fail-map",
			validator: g.Validator(g.R().Len(3).SpecificMessages(galidator.Messages{"len": "must be exactly $length characters long"})),
			in:        map[string]int{"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6},
			panic:     false,
			expected:  []string{"must be exactly 3 characters long"},
		},
		{
			name:      "pass-string",
			validator: g.Validator(g.R().Len(3).SpecificMessages(galidator.Messages{"len": "must be exactly $length characters long"})),
			in:        "111",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-slice",
			validator: g.Validator(g.R().Len(3).SpecificMessages(galidator.Messages{"len": "must be exactly $length characters long"})),
			in:        []int{1, 2, 3},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-map",
			validator: g.Validator(g.R().Len(3).SpecificMessages(galidator.Messages{"len": "must be exactly $length characters long"})),
			in:        map[string]int{"1": 1, "2": 2, "3": 3},
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
