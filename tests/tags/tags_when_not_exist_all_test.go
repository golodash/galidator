package tests

import (
	"context"
	"testing"
)

type WhenNotExistAllRequest struct {
	IdNumber       string `json:"id_number" g:"when_not_exist_all=PassportNumber" when_not_exist_all:"error"`
	PassportNumber string `json:"passport_number" g:"when_not_exist_all=IdNumber" when_not_exist_all:"error"`
}

func TestTagsWhenNotExistAll(t *testing.T) {
	scenarios := []scenario{
		{
			name:      "success_1",
			validator: g.Validator(WhenNotExistAllRequest{}),
			in:        WhenNotExistAllRequest{IdNumber: "1", PassportNumber: "1"},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "success_2",
			validator: g.Validator(WhenNotExistAllRequest{}),
			in:        WhenNotExistAllRequest{IdNumber: "1", PassportNumber: ""},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "success_3",
			validator: g.Validator(WhenNotExistAllRequest{}),
			in:        WhenNotExistAllRequest{IdNumber: "", PassportNumber: "1"},
			panic:     false,
			expected:  nil,
		},
		{
			name:      "fail",
			validator: g.Validator(WhenNotExistAllRequest{}),
			in:        WhenNotExistAllRequest{IdNumber: "", PassportNumber: ""},
			panic:     false,
			expected:  map[string][]string{"id_number": {"error"}, "passport_number": {"error"}},
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
