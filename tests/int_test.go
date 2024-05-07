package tests

import (
	"context"
	"testing"

	"github.com/golodash/galidator"
)

func TestInt(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "pass",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        1,
			panic:     false,
			expected:  nil,
		},
		{
			name:      "fail-1",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        "1",
			panic:     false,
			expected:  []string{"1 is not int"},
		},
		{
			name:      "fail-2",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        map[string]string{"1": "1"},
			panic:     false,
			expected:  []string{"map[1:1] is not int"},
		},
		{
			name:      "int",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        int(1),
			panic:     false,
			expected:  nil,
		},
		{
			name:      "int8",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        int8(1),
			panic:     false,
			expected:  nil,
		},
		{
			name:      "int16",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        int16(1),
			panic:     false,
			expected:  nil,
		},
		{
			name:      "int32",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        int32(1),
			panic:     false,
			expected:  nil,
		},
		{
			name:      "int64",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        int64(1),
			panic:     false,
			expected:  nil,
		},
		{
			name:      "uint",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        uint(1),
			panic:     false,
			expected:  nil,
		},
		{
			name:      "uint8",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        uint8(1),
			panic:     false,
			expected:  nil,
		},
		{
			name:      "uint16",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        uint16(1),
			panic:     false,
			expected:  nil,
		},
		{
			name:      "uint32",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        uint32(1),
			panic:     false,
			expected:  nil,
		},
		{
			name:      "uint64",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        uint64(1),
			panic:     false,
			expected:  nil,
		},
		{
			name:      "float32",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        float32(1),
			panic:     false,
			expected:  []string{"1 is not int"},
		},
		{
			name:      "float64",
			validator: g.Validator(g.R().Int().SpecificMessages(galidator.Messages{"int": "$value is not int"})),
			in:        float64(1),
			panic:     false,
			expected:  []string{"1 is not int"},
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
