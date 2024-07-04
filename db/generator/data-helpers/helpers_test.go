package datahelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"stellarsky.ai/platform/codegen/data-service-generator/db/generator/defs"
)

func TestReadParamValues(t *testing.T) {
	// Test case 1: No parameters
	params := []defs.ParameterRef{}
	request := map[string]interface{}{}
	expected := []interface{}{}
	result, err := ReadParamValues(params, request)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	// Test case 2: Single parameter with index
	params = []defs.ParameterRef{{Name: "param1", Index: -1}}
	request = map[string]interface{}{"param1": 123}
	expected = []interface{}{123}
	result, err = ReadParamValues(params, request)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	// Test case 3: Single parameter with name
	params = []defs.ParameterRef{{Name: "param1", Index: -1}}
	request = map[string]interface{}{"param1": "abc"}
	expected = []interface{}{"abc"}
	result, err = ReadParamValues(params, request)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	// Test case 4: Multiple parameters
	params = []defs.ParameterRef{
		{Name: "param1", Index: 0},
		{Name: "param1", Index: 1},
		{Name: "param2", Index: -1},
	}
	request = map[string]interface{}{
		"param1": []interface{}{123, 456},
		"param2": "abc",
	}
	expected = []interface{}{123, 456, "abc"}
	result, err = ReadParamValues(params, request)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	// Test case 5: Missing parameter
	params = []defs.ParameterRef{{Name: "param1", Index: 0}}
	request = map[string]interface{}{}
	_, err = ReadParamValues(params, request)
	assert.Error(t, err)
}

func TestPostgresToGoType(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		expectedType string
	}{
		{"Test1", "bigint", "int64"},
		{"Test2", "integer", "int"},
		{"Test2A", "smallint", "int16"},
		{"Test3", "boolean", "bool"},
		{"Test4", "real", "float32"},
		{"Test4B", "numeric(5,2)", "float64"},
		{"Test5", "double precision", "float64"},
		{"Test6", "char", "string"},
		{"Test7", "varchar", "string"},
		{"Test7A", "varchar(255)", "string"},
		{"Test8", "text", "string"},
		{"Test9", "uuid", "string"},
		{"Test10", "bytea", "[]byte"},
		{"Test11", "date", "time.Time"},
		{"Test12", "time", "time.Time"},
		{"Test13", "timestamp", "time.Time"},
		{"Test14", "timestamptz", "time.Time"},
		{"Test15", "json", "interface{}"},
		{"Test16", "jsonb", "interface{}"},
		{"Test17", "[]bigint", "[]int64"},
		{"Test18", "[]integer", "[]int"},
		{"Test19", "[]boolean", "[]bool"},
		{"Test20", "[]real", "[]float32"},
		{"Test21", "[]double precision", "[]float64"},
		{"Test22", "[]char", "[]string"},
		{"Test23", "[]varchar", "[]string"},
		{"Test24", "[]text", "[]string"},
		{"Test25", "[]uuid", "[]string"},
		{"Test26", "[]bytea", "[][]byte"},
		{"Test27", "[]date", "[]time.Time"},
		{"Test28", "[]time", "[]time.Time"},
		{"Test29", "[]timestamp", "[]time.Time"},
		{"Test30", "[]timestamptz", "[]time.Time"},
		{"Test31", "[]json", "[]interface{}"},
		{"Test32", "[]jsonb", "[]interface{}"},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualType := PostgresToGoType(tc.input)
			if actualType != tc.expectedType {
				t.Errorf("PostgresToGoType(%s) = %s; want %s", tc.input, actualType, tc.expectedType)
			}
		})
	}
}
