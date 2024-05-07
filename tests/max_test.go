package tests

import (
	"context"
	"testing"

	"github.com/golodash/galidator/v2"
)

func TestMax(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "fail-int",
			validator: g.Validator(g.R().Max(5).SpecificMessages(galidator.Messages{"max": "can be at most $max characters long"})),
			in:        6,
			panic:     false,
			expected:  []string{"can be at most 5 characters long"},
		},
		{
			name:      "fail-string",
			validator: g.Validator(g.R().Max(5).SpecificMessages(galidator.Messages{"max": "can be at most $max characters long"})),
			in:        "abcdef",
			panic:     false,
			expected:  []string{"can be at most 5 characters long"},
		},
		{
			name:      "fail-slice",
			validator: g.Validator(g.R().Max(5).SpecificMessages(galidator.Messages{"max": "can be at most $max characters long"})),
			in:        []int{1, 2, 3, 4, 5, 6},
			panic:     false,
			expected:  []string{"can be at most 5 characters long"},
		},
		{
			name:      "fail-map",
			validator: g.Validator(g.R().Max(5).SpecificMessages(galidator.Messages{"max": "can be at most $max characters long"})),
			in:        map[string]int{"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6},
			panic:     false,
			expected:  []string{"can be at most 5 characters long"},
		},
		{
			name:      "pass-int",
			validator: g.Validator(g.R().Max(5).SpecificMessages(galidator.Messages{"max": "can be at most $max characters long"})),
			in:        5,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-string",
			validator: g.Validator(g.R().Max(5).SpecificMessages(galidator.Messages{"max": "can be at most $max characters long"})),
			in:        "111",
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-slice",
			validator: g.Validator(g.R().Max(5).SpecificMessages(galidator.Messages{"max": "can be at most $max characters long"})),
			in:        []int{1, 2, 3},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "pass-map",
			validator: g.Validator(g.R().Max(5).SpecificMessages(galidator.Messages{"max": "can be at most $max characters long"})),
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
