package tests

import (
	"testing"

	"github.com/golodash/galidator"
	"github.com/golodash/godash/slices"
)

type fromStructUsers struct {
	Name string `json:"name"`
}

type fromStructTest struct {
	ID    int               `json:"id" g:"min=1"`
	Rules []string          `json:"rules" g:"custom_choices" custom_choices:"not included in allowed choices"`
	Users []fromStructUsers `json:"users"`
}

var choices = []string{
	"1", "2", "3",
}

func custom_choices(input interface{}) bool {
	for _, item := range input.([]string) {
		if slices.FindIndex(choices, item) == -1 {
			return false
		}
	}
	return true
}

func TestFromStruct(t *testing.T) {
	v := g.CustomValidators(galidator.Validators{
		"custom_choices": custom_choices,
	}).Validator(fromStructTest{})

	scenarios := []scenario{
		{
			name:      "pass-1",
			validator: v,
			in: fromStructTest{
				ID:    1,
				Rules: []string{"2", "1"},
				Users: nil,
			},
			panic:    false,
			expected: nil,
		},
		{
			name:      "pass-2",
			validator: v,
			in: fromStructTest{
				ID:    1,
				Rules: []string{"2", "1"},
				Users: []fromStructUsers{
					{
						Name: "1",
					},
					{
						Name: "2",
					},
					{
						Name: "3",
					},
				},
			},
			panic:    false,
			expected: nil,
		},
		{
			name:      "fail-1",
			validator: v,
			in: fromStructTest{
				ID:    1,
				Rules: []string{"2", "1", "5"},
				Users: []fromStructUsers{
					{
						Name: "1",
					},
					{
						Name: "2",
					},
					{
						Name: "3",
					},
				},
			},
			panic:    false,
			expected: map[string]interface{}{"rules": []string{"not included in allowed choices"}},
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
