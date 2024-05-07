package tests

import (
	"context"
	"testing"

	"github.com/golodash/galidator"
)

func TestFloat(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "pass",
			validator: g.Validator(g.R().Float().SpecificMessages(galidator.Messages{"float": "$value is not float"})),
			in:        1.1,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "fail",
			validator: g.Validator(g.R().Float().SpecificMessages(galidator.Messages{"float": "$value is not float"})),
			in:        "1",
			panic:     false,
			expected:  []string{"1 is not float"},
		},
		{
			name:      "float",
			validator: g.Validator(g.R().Float().SpecificMessages(galidator.Messages{"float": "$value is not float"})),
			in:        int(1),
			panic:     false,
			expected:  []string{"1 is not float"},
		},
		{
			name:      "int8",
			validator: g.Validator(g.R().Float().SpecificMessages(galidator.Messages{"float": "$value is not float"})),
			in:        int8(1),
			panic:     false,
			expected:  []string{"1 is not float"},
		},
		{
			name:      "int16",
			validator: g.Validator(g.R().Float().SpecificMessages(galidator.Messages{"float": "$value is not float"})),
			in:        int16(1),
			panic:     false,
			expected:  []string{"1 is not float"},
		},
		{
			name:      "int32",
			validator: g.Validator(g.R().Float().SpecificMessages(galidator.Messages{"float": "$value is not float"})),
			in:        int32(1),
			panic:     false,
			expected:  []string{"1 is not float"},
		},
		{
			name:      "int64",
			validator: g.Validator(g.R().Float().SpecificMessages(galidator.Messages{"float": "$value is not float"})),
			in:        int64(1),
			panic:     false,
			expected:  []string{"1 is not float"},
		},
		{
			name:      "uint",
			validator: g.Validator(g.R().Float().SpecificMessages(galidator.Messages{"float": "$value is not float"})),
			in:        uint(1),
			panic:     false,
			expected:  []string{"1 is not float"},
		},
		{
			name:      "uint8",
			validator: g.Validator(g.R().Float().SpecificMessages(galidator.Messages{"float": "$value is not float"})),
			in:        uint8(1),
			panic:     false,
			expected:  []string{"1 is not float"},
		},
		{
			name:      "uint16",
			validator: g.Validator(g.R().Float().SpecificMessages(galidator.Messages{"float": "$value is not float"})),
			in:        uint16(1),
			panic:     false,
			expected:  []string{"1 is not float"},
		},
		{
			name:      "uint32",
			validator: g.Validator(g.R().Float().SpecificMessages(galidator.Messages{"float": "$value is not float"})),
			in:        uint32(1),
			panic:     false,
			expected:  []string{"1 is not float"},
		},
		{
			name:      "uint64",
			validator: g.Validator(g.R().Float().SpecificMessages(galidator.Messages{"float": "$value is not float"})),
			in:        uint64(1),
			panic:     false,
			expected:  []string{"1 is not float"},
		},
		{
			name:      "float32",
			validator: g.Validator(g.R().Float().SpecificMessages(galidator.Messages{"float": "$value is not float"})),
			in:        float32(1),
			panic:     false,
			expected:  nil,
		},
		{
			name:      "float64",
			validator: g.Validator(g.R().Float().SpecificMessages(galidator.Messages{"float": "$value is not float"})),
			in:        float64(1),
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
