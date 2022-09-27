package tests

import (
	"testing"

	"github.com/golodash/galidator"
)

func TestType(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "pass",
			validator: g.Validator(g.R().Type(1.0).SpecificMessages(galidator.Messages{"type": "not float64"})),
			in:        1.1,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "fail",
			validator: g.Validator(g.R().Type(1.0).SpecificMessages(galidator.Messages{"type": "not float64"})),
			in:        "1",
			panic:     false,
			expected:  []string{"not float64"},
		},
		{
			name:      "float",
			validator: g.Validator(g.R().Type(1.0).SpecificMessages(galidator.Messages{"type": "not float64"})),
			in:        int(1),
			panic:     false,
			expected:  []string{"not float64"},
		},
		{
			name:      "int8",
			validator: g.Validator(g.R().Type(1.0).SpecificMessages(galidator.Messages{"type": "not float64"})),
			in:        int8(1),
			panic:     false,
			expected:  []string{"not float64"},
		},
		{
			name:      "int16",
			validator: g.Validator(g.R().Type(1.0).SpecificMessages(galidator.Messages{"type": "not float64"})),
			in:        int16(1),
			panic:     false,
			expected:  []string{"not float64"},
		},
		{
			name:      "int32",
			validator: g.Validator(g.R().Type(1.0).SpecificMessages(galidator.Messages{"type": "not float64"})),
			in:        int32(1),
			panic:     false,
			expected:  []string{"not float64"},
		},
		{
			name:      "int64",
			validator: g.Validator(g.R().Type(1.0).SpecificMessages(galidator.Messages{"type": "not float64"})),
			in:        int64(1),
			panic:     false,
			expected:  []string{"not float64"},
		},
		{
			name:      "uint",
			validator: g.Validator(g.R().Type(1.0).SpecificMessages(galidator.Messages{"type": "not float64"})),
			in:        uint(1),
			panic:     false,
			expected:  []string{"not float64"},
		},
		{
			name:      "uint8",
			validator: g.Validator(g.R().Type(1.0).SpecificMessages(galidator.Messages{"type": "not float64"})),
			in:        uint8(1),
			panic:     false,
			expected:  []string{"not float64"},
		},
		{
			name:      "uint16",
			validator: g.Validator(g.R().Type(1.0).SpecificMessages(galidator.Messages{"type": "not float64"})),
			in:        uint16(1),
			panic:     false,
			expected:  []string{"not float64"},
		},
		{
			name:      "uint32",
			validator: g.Validator(g.R().Type(1.0).SpecificMessages(galidator.Messages{"type": "not float64"})),
			in:        uint32(1),
			panic:     false,
			expected:  []string{"not float64"},
		},
		{
			name:      "uint64",
			validator: g.Validator(g.R().Type(1.0).SpecificMessages(galidator.Messages{"type": "not float64"})),
			in:        uint64(1),
			panic:     false,
			expected:  []string{"not float64"},
		},
		{
			name:      "float32",
			validator: g.Validator(g.R().Type(1.0).SpecificMessages(galidator.Messages{"type": "not float64"})),
			in:        float32(1),
			panic:     false,
			expected:  []string{"not float64"},
		},
		{
			name:      "float64",
			validator: g.Validator(g.R().Type(1.0).SpecificMessages(galidator.Messages{"type": "not float64"})),
			in:        float64(1),
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
