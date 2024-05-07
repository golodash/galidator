package tests

import (
	"context"
	"testing"

	"github.com/golodash/galidator/v2"
)

func TestWhenNotExistOne(t *testing.T) {
	v := g.Validator(g.R().Complex(galidator.Rules{
		"id":       g.R().Int(),
		"name":     g.R().String(),
		"username": g.R().WhenNotExistOne("id", "name").String().SpecificMessages(galidator.Messages{"when_not_exist_one": "we are required now"}),
	}))

	scenarios := []scenario{
		{
			name:      "fail-1",
			validator: v,
			in: map[string]interface{}{
				"id":       1,
				"name":     "",
				"username": "",
			},
			panic:    false,
			expected: map[string][]string{"username": {"we are required now"}},
		},
		{
			name:      "fail-2",
			validator: v,
			in: map[string]interface{}{
				"id":       0,
				"name":     "name",
				"username": "",
			},
			panic:    false,
			expected: map[string][]string{"username": {"we are required now"}},
		},
		{
			name:      "fail-3",
			validator: v,
			in: map[string]interface{}{
				"id":       0,
				"name":     "",
				"username": "",
			},
			panic:    false,
			expected: map[string][]string{"username": {"we are required now"}},
		},
		{
			name:      "pass-1",
			validator: v,
			in: map[string]interface{}{
				"id":       5,
				"name":     "name",
				"username": "",
			},
			panic:    false,
			expected: nil,
		},
		{
			name:      "pass-2",
			validator: v,
			in: map[string]interface{}{
				"id":       1,
				"name":     "",
				"username": "ss",
			},
			panic:    false,
			expected: nil,
		},
		{
			name:      "pass-3",
			validator: v,
			in: map[string]interface{}{
				"id":       0,
				"name":     "name",
				"username": "ss",
			},
			panic:    false,
			expected: nil,
		},
		{
			name:      "pass-4",
			validator: v,
			in: map[string]interface{}{
				"id":       1,
				"name":     "name",
				"username": "ss",
			},
			panic:    false,
			expected: nil,
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
