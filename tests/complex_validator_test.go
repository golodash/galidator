package tests

import (
	"context"
	"testing"

	"github.com/golodash/galidator/v2"
)

type ComplexValidator struct {
	Name        string
	Description string
}

func TestComplexValidator(t *testing.T) {
	v := g.ComplexValidator(galidator.Rules{
		"Name":        g.RuleSet("name").Required(),
		"Description": g.RuleSet("description").Required(),
	})

	scenarios := []scenario{
		{
			name:      "pass-1-map",
			validator: v,
			in: map[string]interface{}{
				"name":        "name",
				"description": "description",
			},
			panic:    false,
			expected: nil,
		},
		{
			name:      "fail-1-map",
			validator: v,
			in: map[string]interface{}{
				"name":        "",
				"description": "description",
			},
			panic:    false,
			expected: map[string]interface{}{"name": []string{"required"}},
		},
		{
			name:      "fail-2-map",
			validator: v,
			in: map[string]interface{}{
				"name":        "name",
				"description": "",
			},
			panic:    false,
			expected: map[string]interface{}{"description": []string{"required"}},
		},
		{
			name:      "fail-3-map",
			validator: v,
			in: map[string]interface{}{
				"name":        "",
				"description": "",
			},
			panic:    false,
			expected: map[string]interface{}{"name": []string{"required"}, "description": []string{"required"}},
		},
		{
			name:      "pass-1-struct",
			validator: v,
			in: ComplexValidator{
				Name:        "name",
				Description: "description",
			},
			panic:    false,
			expected: nil,
		},
		{
			name:      "fail-1-struct",
			validator: v,
			in: ComplexValidator{
				Name:        "",
				Description: "description",
			},
			panic:    false,
			expected: map[string]interface{}{"name": []string{"required"}},
		},
		{
			name:      "fail-2-struct",
			validator: v,
			in: ComplexValidator{
				Name:        "name",
				Description: "",
			},
			panic:    false,
			expected: map[string]interface{}{"description": []string{"required"}},
		},
		{
			name:      "fail-3-struct",
			validator: v,
			in: ComplexValidator{
				Name:        "",
				Description: "",
			},
			panic:    false,
			expected: map[string]interface{}{"name": []string{"required"}, "description": []string{"required"}},
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
