package datahelpers

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"stellarsky.ai/platform/codegen/data-service-generator/db/generator/defs"
)

func filterData(startCounter uint32) ([]defs.Filter, string, []defs.ParameterRef) {
	filters := []defs.Filter{
		{Attribute: "age", Operator: ">=", ParamName: "age"},
		{Attribute: "status", Operator: "IN", ParamName: "status"},
		{
			Operator: "OR",
			Conditions: []defs.Filter{
				{Attribute: "salary", Operator: "BETWEEN", ParamName: "salary_range"},
				{Attribute: "position", Operator: "=", ParamName: "position"},
				{Operator: "AND", Conditions: []defs.Filter{
					{Attribute: "deparment", Operator: "=", ParamName: "deparement"},
					{Attribute: "experience", Operator: ">=", ParamName: "experience"},
					{Operator: "AND", Conditions: []defs.Filter{
						{Attribute: "name", Transformation: "CHAR_LENGTH", Operator: ">", ParamName: "name"},
						{Attribute: "created_at", Transformation: "DATE", Operator: "=", ParamName: "created_at"},
					}},
				}},
			},
		},
		{
			Operator: "NOT",
			Conditions: []defs.Filter{
				{Attribute: "terminated", Operator: "=", ParamName: "terminated"},
			},
		},
	}

	expected := fmt.Sprintf("(age >= $%d AND status = ANY($%d) AND ((salary BETWEEN $%d AND $%d) OR position = $%d OR (deparment = $%d AND experience >= $%d AND (CHAR_LENGTH(name) > $%d AND DATE(created_at) = $%d))) AND NOT(terminated = $%d))",
		startCounter, startCounter+1, startCounter+2, startCounter+3, startCounter+4, startCounter+5, startCounter+6, startCounter+7, startCounter+8, startCounter+9)

	paramsMap := []defs.ParameterRef{
		{Name: "age", Index: -1},
		{Name: "status", Index: -1},
		{Name: "salary_range", Index: 0},
		{Name: "salary_range", Index: 1},
		{Name: "position", Index: -1},
		{Name: "deparement", Index: -1},
		{Name: "experience", Index: -1},
		{Name: "name", Index: -1},
		{Name: "created_at", Index: -1},
		{Name: "terminated", Index: -1},
	}
	return filters, expected, paramsMap
}

func valueMapData(startCounter uint32) ([]defs.Update, string, []defs.ParameterRef) {

	values := []defs.Update{
		{Attribute: "age", ParamName: "age"},
		{Attribute: "status", ParamName: "status"},
		{Attribute: "salary", ParamName: "salary"},
		{Attribute: "position", ParamName: "position"},
		{Attribute: "department", ParamName: "department"},
		{Attribute: "experience", ParamName: "experience"},
		{Attribute: "name", ParamName: "name"},
		{Attribute: "created_at", ParamName: "created_at"},
		{Attribute: "terminated", ParamName: "terminated"},
	}

	expected := fmt.Sprintf("age = $%d, status = $%d, salary = $%d, position = $%d, department = $%d, experience = $%d, name = $%d, created_at = $%d, terminated = $%d",
		startCounter, startCounter+1, startCounter+2, startCounter+3, startCounter+4, startCounter+5, startCounter+6, startCounter+7, startCounter+8)

	paramsMap := []defs.ParameterRef{
		{Name: "age", Index: -1},
		{Name: "status", Index: -1},
		{Name: "salary", Index: -1},
		{Name: "position", Index: -1},
		{Name: "department", Index: -1},
		{Name: "experience", Index: -1},
		{Name: "name", Index: -1},
		{Name: "created_at", Index: -1},
		{Name: "terminated", Index: -1},
	}
	return values, expected, paramsMap
}

func TestPrepareFilters(t *testing.T) {
	filters := []defs.Filter{
		{Attribute: "age", Operator: ">=", ParamName: "age"},
		{Attribute: "status", Operator: "IN", ParamName: "status"},
		{
			Operator: "OR",
			Conditions: []defs.Filter{
				{Attribute: "salary", Operator: "BETWEEN", ParamName: "salary_range"},
				{Attribute: "position", Operator: "=", ParamName: "position"},
				{Operator: "AND", Conditions: []defs.Filter{
					{Attribute: "department", Operator: "=", ParamName: "department"},
					{Attribute: "experience", Operator: ">=", ParamName: "experience"}}},
			},
		},
		{
			Operator: "NOT",
			Conditions: []defs.Filter{
				{Attribute: "terminated", Operator: "=", ParamName: "terminated"},
			},
		},
		{Attribute: "name", Transformation: "CHAR_LENGTH", Operator: ">", ParamName: "name"},
		{Attribute: "created_at", Transformation: "DATE", Operator: "=", ParamName: "created_at"},
	}

	expected := "(age >= $1 AND status = ANY($2) AND ((salary BETWEEN $3 AND $4) OR position = $5 OR (department = $6 AND experience >= $7)) AND NOT(terminated = $8) AND CHAR_LENGTH(name) > $9 AND DATE(created_at) = $10)"
	expectedParamsMap := []defs.ParameterRef{
		{Index: -1, Name: "age"},
		{Index: -1, Name: "status"},
		{Index: 0, Name: "salary_range"},
		{Index: 1, Name: "salary_range"},
		{Index: -1, Name: "position"},
		{Index: -1, Name: "department"},
		{Index: -1, Name: "experience"},
		{Index: -1, Name: "terminated"},
		{Index: -1, Name: "name"},
		{Index: -1, Name: "created_at"},
	}
	result, paramsMap := PrepareFilters(filters)

	fmt.Printf("ParamsMap: %v\n", paramsMap)
	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}

	if len(paramsMap) != len(expectedParamsMap) {
		t.Errorf("Expected %v, but got %v", expectedParamsMap, paramsMap)
	}

	for i, param := range paramsMap {
		assert.Equal(t, expectedParamsMap[i], param, "Mismatch at index %d. Expected %v, but got %v", i, expectedParamsMap[i], paramsMap[i])
	}

	// Test case: complex filter
	complexFilters, complexExpected, complexExpectedParams := filterData(uint32(1))
	complexQuery, complexParams := PrepareFilters(complexFilters)
	if complexQuery != complexExpected {
		t.Errorf("Expected '%s', but got '%s'", complexExpected, complexQuery)
	}
	if len(complexParams) != len(complexExpectedParams) && !reflect.DeepEqual(complexParams, complexExpectedParams) {
		t.Errorf("Expected %v, but got %v", complexExpectedParams, complexParams)
	}
}

func TestPrepareUpdateStmt(t *testing.T) {
	// Test case: empty AccessConfig
	updateConfig := &defs.AccessConfig{}
	setClause, whereClause, paramsMap := PrepareUpdateStmt(updateConfig)
	if setClause != "" || whereClause != "" || len(paramsMap) != 0 {
		t.Errorf("Set clause should be empty, where clause should be empty, and params map should be empty")
	}

	// Test case: non-empty AccessConfig
	updateConfig = &defs.AccessConfig{
		Filter: []defs.Filter{
			{
				Attribute: "attribute1",
				Operator:  "=",
				ParamName: "param1",
			},
			{
				Attribute: "attribute2",
				Operator:  ">",
				ParamName: "param2",
			},
		},
		Set: []defs.Update{
			{
				Attribute: "attribute3",
				ParamName: "param3",
			},
			{
				Attribute: "attribute4",
				ParamName: "param4",
			},
		},
		Autoincrement:    []string{"attribute5", "attribute6"},
		CaptureTimestamp: []string{"attribute7", "attribute8"},
	}
	setClause, whereClause, paramsMap = PrepareUpdateStmt(updateConfig)
	expectedSetClause := "attribute_3 = $1, attribute_4 = $2, attribute_5 = attribute_5 + 1, attribute_6 = attribute_6 + 1, attribute_7 = NOW(), attribute_8 = NOW()"
	expectedWhereClause := "(attribute_1 = $3 AND attribute_2 > $4)"
	expectedParamsMap := []defs.ParameterRef{
		{
			Name:  "param3",
			Index: -1,
		},
		{
			Name:  "param4",
			Index: -1,
		},
		{
			Name:  "param1",
			Index: -1,
		},
		{
			Name:  "param2",
			Index: -1,
		},
	}
	if setClause != expectedSetClause {
		t.Errorf("Set clause should be %s, got: %s", expectedSetClause, setClause)
	}
	if whereClause != expectedWhereClause {
		t.Errorf("Where clause should be %s, got: %s", expectedWhereClause, whereClause)
	}
	if !reflect.DeepEqual(paramsMap, expectedParamsMap) {
		t.Errorf("Expected %v, but got %v", expectedParamsMap, paramsMap)
	}
}

func TestMakeUpdateQuery(t *testing.T) {
	table := "test_table"

	updateList, expectedSetClause, expectedUpdateParams := valueMapData(uint32(1))
	filter, expectedFilterClause, expectedFilterParams := filterData(uint32(len(expectedUpdateParams)) + 1)
	updateConfig := &defs.AccessConfig{
		Set:    updateList,
		Filter: filter,
	}

	expectedQuery := fmt.Sprintf("UPDATE test_table SET %v WHERE (1 = 1) AND %v", expectedSetClause, expectedFilterClause)
	expectedParams := []defs.ParameterRef{}
	expectedParams = append(expectedParams, expectedUpdateParams...)
	expectedParams = append(expectedParams, expectedFilterParams...)
	query, params := MakeUpdateQuery(table, updateConfig)
	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedParams, params)
}

func TestMakeAddQuery(t *testing.T) {
	table := "test_table"

	addConfig := &defs.AccessConfig{
		Values: []defs.Update{{Attribute: "attr1", ParamName: "p1"}, {Attribute: "attr2", ParamName: "p2"}},
	}
	expectedQuery := "INSERT INTO test_table (attr_1, attr_2) VALUES ($1, $2) RETURNING id"
	expectedParams := []defs.ParameterRef{{Name: "p1", Index: -1}, {Name: "p2", Index: -1}}

	query, params := MakeAddQuery(table, addConfig)
	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedParams, params)
}

func TestMakeAddOrReplaceQuery(t *testing.T) {
	table := "test_table"
	addConfig := &defs.AccessConfig{
		Values: []defs.Update{{Attribute: "attr1", ParamName: "p1"}, {Attribute: "attr2", ParamName: "p2"}},
	}
	expectedQuery := "INSERT INTO test_table (attr_1, attr_2) VALUES ($1, $2) ON CONFLICT DO UPDATE SET attr_1 = $3, attr_2 = $4 RETURNING id, (xmax = 0)"
	expectedParams := []defs.ParameterRef{{Name: "p1", Index: -1}, {Name: "p2", Index: -1}, {Name: "p1", Index: -1}, {Name: "p2", Index: -1}}

	query, params := MakeAddOrReplaceQuery(table, addConfig)
	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedParams, params)
}

func TestMakeDeleteQuery(t *testing.T) {
	table := "test_table"
	deleteConfig := &defs.AccessConfig{
		Filter: []defs.Filter{
			{
				Attribute: "attr1",
				Operator:  "=",
				ParamName: "p1",
			},
			{
				Attribute: "attr2",
				Operator:  ">",
				ParamName: "p2",
			},
		},
	}
	expectedQuery := "DELETE FROM test_table WHERE (1 = 1) AND (attr_1 = $1 AND attr_2 > $2)"
	expectedParams := []defs.ParameterRef{{Name: "p1", Index: -1}, {Name: "p2", Index: -1}}

	query, params := MakeDeleteQuery(table, deleteConfig)
	assert.Equal(t, expectedQuery, query)
	assert.Equal(t, expectedParams, params)
}
