package tests

import (
	"context"
	"testing"
)

var translates = map[string]string{
	"required": "this is required which is translated",
}

func testTranslator(input string) string {
	if out, ok := translates[input]; ok {
		return out
	}
	return input
}

func TestTranslator(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "fail-int",
			validator: g.Validator(g.R().Required()),
			in:        0,
			panic:     false,
			expected:  []string{"this is required which is translated"},
		},
		{
			name:      "fail-float",
			validator: g.Validator(g.R().Required()),
			in:        0.0,
			panic:     false,
			expected:  []string{"this is required which is translated"},
		},
		{
			name:      "fail-string",
			validator: g.Validator(g.R().Required()),
			in:        "",
			panic:     false,
			expected:  []string{"this is required which is translated"},
		},
		{
			name:      "fail-slice",
			validator: g.Validator(g.R().Required()),
			in:        []int{},
			panic:     false,
			expected:  []string{"this is required which is translated"},
		},
		{
			name:      "fail-map",
			validator: g.Validator(g.R().Required()),
			in:        map[string]string{},
			panic:     false,
			expected:  []string{"this is required which is translated"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			defer deferTestCases(t, s.panic, s.expected)

			output := s.validator.Validate(context.TODO(), s.in, testTranslator)
			check(t, s.expected, output)
		})
	}
}
