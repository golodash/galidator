package tests

import (
	"context"
	"testing"
)

type choices struct {
	Name []string `g:"c.choices=1&2&3" c.choices:"choices failed"`
}

type choices2 struct {
	Name []string `g:"child.choices=1&2&3" child.choices:"choices failed"`
}

func TestChoices(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "success",
			validator: g.Validator(choices{}),
			in:        choices{Name: []string{"1", "2", "3", "1"}},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "fail",
			validator: g.Validator(choices{}),
			in:        choices{Name: []string{"a", "2", "3", "1"}},
			panic:     false,
			expected:  map[string]map[string][]string{"Name": {"0": []string{"choices failed"}}},
		},
		{
			name:      "success_2",
			validator: g.Validator(choices2{}),
			in:        choices2{Name: []string{"a", "2", "3", "1"}},
			panic:     false,
			expected:  map[string]map[string][]string{"Name": {"0": []string{"choices failed"}}},
		},
		{
			name:      "fail_2",
			validator: g.Validator(choices2{}),
			in:        choices2{Name: []string{"1", "2", "cd", "sw"}},
			panic:     false,
			expected:  map[string]map[string][]string{"Name": {"2": []string{"choices failed"}, "3": []string{"choices failed"}}},
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
