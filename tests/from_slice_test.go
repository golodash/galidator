package tests

import (
	"context"
	"testing"

	"github.com/golodash/galidator/v2"
)

func TestFromSlice(t *testing.T) {
	v := g.CustomValidators(galidator.Validators{
		"custom_choices": custom_choices,
	}).Validator([]fromStructTest{})

	scenarios := []scenario{
		{
			name:      "pass-1",
			validator: v,
			in: []fromStructTest{
				{
					ID:    1,
					Rules: []string{"2", "1"},
					Users: nil,
				},
			},
			panic:    false,
			expected: nil,
		},
		{
			name:      "pass-2",
			validator: v,
			in: []fromStructTest{
				{
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
			},
			panic:    false,
			expected: nil,
		},
		{
			name:      "fail-1",
			validator: v,
			in: []fromStructTest{
				{
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
			},
			panic:    false,
			expected: map[string]interface{}{"0": map[string]interface{}{"rules": []string{"not included in allowed choices"}}},
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
