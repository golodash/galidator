package tests

import (
	"testing"

	"github.com/golodash/galidator"
)

func TestSuperComplex(t *testing.T) {
	v := g.Validator(g.R().Slice().Children(
		g.R().Map().Complex(galidator.Rules{
			"id":    g.R().Int().Min(1).AlwaysCheckRules(),
			"rules": g.R().Slice().Children(g.R().Choices("1", "2", "3")),
			"users": g.R().Slice().Children(g.R().Complex(galidator.Rules{
				"name": g.R().String(),
			})),
		}),
	))

	scenarios := []scenario{
		{
			name:      "pass-1",
			validator: v,
			in: []interface{}{map[string]interface{}{
				"id":    1,
				"rules": []string{"2", "1"},
				"users": nil,
			}},
			panic:    false,
			expected: nil,
		},
		{
			name:      "pass-2",
			validator: v,
			in: []interface{}{map[string]interface{}{
				"id":    1,
				"rules": []string{"2", "1"},
				"users": []interface{}{
					map[string]string{
						"name": "1",
					},
					map[string]string{
						"name": "2",
					},
					map[string]string{
						"name": "3",
					},
				},
			}},
			panic:    false,
			expected: nil,
		},
		{
			name:      "fail-1",
			validator: v,
			in: []interface{}{map[string]interface{}{
				"id":    1,
				"rules": []string{"2", "1"},
				"users": []interface{}{
					map[string]string{
						"name": "1",
					},
					map[string]string{
						"name": "2",
					},
					map[string]int{
						"name": 3,
					},
				},
			}},
			panic:    false,
			expected: map[string]interface{}{"0": map[string]interface{}{"users": map[string]interface{}{"2": map[string]interface{}{"name": []string{"not a string"}}}}},
		},
		{
			name:      "fail-2",
			validator: v,
			in: []interface{}{map[string]interface{}{
				"id":    0,
				"rules": []string{"2", "1"},
				"users": []interface{}{
					map[string]string{
						"name": "1",
					},
					map[string]string{
						"name": "2",
					},
					map[string]string{
						"name": "3",
					},
				},
			}},
			panic:    false,
			expected: map[string]interface{}{"0": map[string]interface{}{"id": []string{"id's length must be higher equal to 1"}}},
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
