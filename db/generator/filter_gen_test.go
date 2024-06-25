package generator

import (
	"fmt"
	"testing"
)

// func TestParseFilter(t *testing.T) {
// 	tests := []struct {
// 		filter   Filter
// 		expected string
// 	}{
// 		{Filter{Attribute: "age", Operator: ">=", Value: 30}, "age >= 30"},
// 		{Filter{Attribute: "status", Operator: "IN", Value: []interface{}{"active", "pending"}}, "status IN ('active', 'pending')"},
// 		{Filter{
// 			Operator: "OR",
// 			Conditions: []Filter{
// 				{Attribute: "salary", Operator: ">", Value: 50000},
// 				{Attribute: "position", Operator: "=", Value: "manager"},
// 			},
// 		}, "(salary > 50000 OR position = 'manager')"},
// 		{Filter{
// 			Operator: "NOT",
// 			Conditions: []Filter{
// 				{Attribute: "terminated", Operator: "=", Value: true},
// 			},
// 		}, "NOT(terminated = true)"},
// 		{Filter{Attribute: "name", Transformation: "CHAR_LENGTH", Operator: ">", Value: 5}, "CHAR_LENGTH(name) > 5"},
// 	}

// 	for _, test := range tests {
// 		result, _ := ParseFilter(test.filter, false, uint32(0))
// 		if result != test.expected {
// 			t.Errorf("Expected '%s', but got '%s'", test.expected, result)
// 		}
// 	}
// }

func TestPrepareFilters(t *testing.T) {
	filters := Filters{
		Filters: []Filter{
			{Attribute: "age", Operator: ">=", ParamName: "age"},
			{Attribute: "status", Operator: "IN", ParamName: "status"},
			{
				Operator: "OR",
				Conditions: []Filter{
					{Attribute: "salary", Operator: "BETWEEN", ParamName: "salary_range"},
					{Attribute: "position", Operator: "=", ParamName: "position"},
					{Operator: "AND", Conditions: []Filter{
						{Attribute: "deparment", Operator: "=", ParamName: "deparement"},
						{Attribute: "experience", Operator: ">=", ParamName: "experience"}}},
				},
			},
			{
				Operator: "NOT",
				Conditions: []Filter{
					{Attribute: "terminated", Operator: "=", ParamName: "terminated"},
				},
			},
			{Attribute: "name", Transformation: "CHAR_LENGTH", Operator: ">", ParamName: "name"},
			{Attribute: "created_at", Transformation: "DATE", Operator: "=", ParamName: "created_at"},
		},
	}

	expected := "(age >= $1 AND status = ANY($2) AND ((salary BETWEEN $3 AND $4) OR position = $5 OR (deparment = $6 AND experience >= $7)) AND NOT(terminated = $8) AND CHAR_LENGTH(name) > $9 AND DATE(created_at) = $10)"
	expectedParamsMap := []ParameterRef{
		{Index: 0, Name: ""},
		{Index: -1, Name: "age"},
		{Index: -1, Name: "status"},
		{Index: 0, Name: "salary_range"},
		{Index: 1, Name: "salary_range"},
		{Index: -1, Name: "position"},
		{Index: -1, Name: "deparement"},
		{Index: -1, Name: "experience"},
		{Index: -1, Name: "terminated"},
		{Index: -1, Name: "name"},
		{Index: -1, Name: "created_at"},
	}
	result, paramsMap := PrepareFilters(filters.Filters)

	fmt.Printf("ParamsMap: %v\n", paramsMap)
	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}

	if len(paramsMap) != len(expectedParamsMap) {
		t.Errorf("Expected %v, but got %v", expectedParamsMap, paramsMap)
	}
}
