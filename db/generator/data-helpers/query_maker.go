package datahelpers

import (
	"fmt"
	"strings"

	"stellarsky.ai/platform/codegen/data-service-generator/base"
	"stellarsky.ai/platform/codegen/data-service-generator/db/generator/defs"
)

func applyTransformation(attribute string, transformation string) string {
	if transformation == "" {
		return attribute
	}
	return fmt.Sprintf("%s(%s)", transformation, attribute)
}

func makePreparedCounter(counter *uint32) string {
	result := fmt.Sprintf("$%d", *counter)
	*counter += 1
	return result
}

func buildPrepareStmt(filter defs.Filter, counter *uint32, paramsMap *[]defs.ParameterRef) string {
	attribute := applyTransformation(filter.Attribute, filter.Transformation)
	switch strings.ToUpper(filter.Operator) {
	case OperatorEquals, OperatorNotEquals, OperatorLessThan, OperatorLessThanEquals, OperatorGreaterThan, OperatorGreaterThanEquals:
		result := fmt.Sprintf("%s %s %v", attribute, filter.Operator, makePreparedCounter(counter))
		*paramsMap = append(*paramsMap, defs.ParameterRef{
			Name:  filter.ParamName,
			Index: -1,
		})
		return result
	case OperatorBetween:
		lowEnd := makePreparedCounter(counter)
		highEnd := makePreparedCounter(counter)
		*paramsMap = append(*paramsMap, defs.ParameterRef{Name: filter.ParamName, Index: 0}, defs.ParameterRef{Name: filter.ParamName, Index: 1})
		result := fmt.Sprintf("(%s BETWEEN %v AND %v)", attribute, lowEnd, highEnd)
		return result
	case OperatorIn:
		result := fmt.Sprintf("%s = ANY(%s)", attribute, makePreparedCounter(counter))
		*paramsMap = append(*paramsMap, defs.ParameterRef{
			Name:  filter.ParamName,
			Index: -1,
		})
		return result
	case LogicalAnd, LogicalOr, LogicalNot:
		subConditions := make([]string, len(filter.Conditions))
		for i, sub := range filter.Conditions {
			subConditions[i] = buildPrepareStmt(sub, counter, paramsMap)
		}
		if filter.Operator == LogicalNot {
			return fmt.Sprintf("NOT(%s)", subConditions[0])
		}
		return fmt.Sprintf("(%s)", strings.Join(subConditions, fmt.Sprintf(" %s ", filter.Operator)))
	default:
		return ""
	}
}

func argsClause(updates []defs.Update, counter *uint32, paramsMap *[]defs.ParameterRef) string {
	var values []string
	for i := range updates {
		values = append(values, fmt.Sprintf("$%d", *counter))
		*paramsMap = append(*paramsMap, defs.ParameterRef{
			Name:  updates[i].ParamName,
			Index: -1,
		})
		*counter += 1
	}
	return strings.Join(values, ", ")
}

func setClause(updates []defs.Update, counter *uint32, paramsMap *[]defs.ParameterRef) string {
	var clauses []string
	for _, update := range updates {
		clauses = append(clauses, fmt.Sprintf("%s = $%d", update.Attribute, *counter))
		*paramsMap = append(*paramsMap, defs.ParameterRef{
			Name:  update.ParamName,
			Index: -1,
		})
		*counter += 1
	}
	return strings.Join(clauses, ", ")
}

func autoincrementClause(attributes []string) string {
	var clauses []string
	for _, attribute := range attributes {
		clauses = append(clauses, fmt.Sprintf("%s = %s + 1", attribute, attribute))
	}
	return strings.Join(clauses, ", ")
}

func captureTimestampClause(attributes []string) string {
	var clauses []string
	for _, attribute := range attributes {
		clauses = append(clauses, fmt.Sprintf("%s = NOW()", attribute))
	}
	return strings.Join(clauses, ", ")
}

func prepareFilters(filters []defs.Filter, counter *uint32, paramsMap *[]defs.ParameterRef) string {
	conditions := make([]string, len(filters))
	for i, filter := range filters {
		conditions[i] = buildPrepareStmt(filter, counter, paramsMap)
	}
	result := strings.Join(conditions, " AND ")
	base.LOG.Debug("Prepared filters for", "input filters", filters, "result", result, "paramsMap", paramsMap)
	if len(strings.Trim(result, " ")) == 0 {
		return ""
	}
	return fmt.Sprintf("(%s)", result)
}

// Returns a prepared condition in Postgresql format, example "column = $1 AND column2 = $2"
// and returns id to param name mapping in array, example ["", "column" "column2"], keeping first string empty to start with $1
func PrepareFilters(filters []defs.Filter) (string, []defs.ParameterRef) {
	counter := uint32(1)
	paramsMap := make([]defs.ParameterRef, 0)
	paramsMap = append(paramsMap, defs.ParameterRef{Index: 0})
	result := prepareFilters(filters, &counter, &paramsMap)
	return result, paramsMap
}

// Given AccessConfig, for update we need to generate the prepared statement
// We need to generate the set clause part and where clause part
func PrepareUpdateStmt(updateConfig *defs.AccessConfig) (string, string, []defs.ParameterRef) {
	counter := uint32(1)
	paramsMap := make([]defs.ParameterRef, 0)
	paramsMap = append(paramsMap, defs.ParameterRef{Index: 0})

	setClause := setClause(updateConfig.Set, &counter, &paramsMap)
	autoincClause := autoincrementClause(updateConfig.Autoincrement)
	captureClause := captureTimestampClause(updateConfig.CaptureTimestamp)
	whereClause := prepareFilters(updateConfig.Filter, &counter, &paramsMap)

	base.LOG.Debug("Prepared update stmt for", "updateConfig", updateConfig, "setClause", setClause,
		"autoincClause", autoincClause, "captureClause", captureClause, "whereClause", whereClause, "paramsMap", paramsMap)
	allClauses := make([]string, 0, 3)

	if setClause != "" {
		allClauses = append(allClauses, setClause)
	}
	if autoincClause != "" {
		allClauses = append(allClauses, autoincClause)
	}
	if captureClause != "" {
		allClauses = append(allClauses, captureClause)
	}
	if len(allClauses) > 0 {
		return strings.Join(allClauses, ", "), whereClause, paramsMap
	}
	return "", whereClause, paramsMap
}

func PrepareAddStmt(addConfig *defs.AccessConfig) (string, []defs.ParameterRef) {
	counter := uint32(1)
	paramsMap := make([]defs.ParameterRef, 0)
	paramsMap = append(paramsMap, defs.ParameterRef{Index: 0})
	insertClause := argsClause(addConfig.Values, &counter, &paramsMap)
	return insertClause, paramsMap
}

func PrepareAddOrReplaceStmt(addConfig *defs.AccessConfig) (string, string, []defs.ParameterRef) {
	counter := uint32(1)
	paramsMap := make([]defs.ParameterRef, 0)
	paramsMap = append(paramsMap, defs.ParameterRef{Index: 0})
	insertClause := argsClause(addConfig.Values, &counter, &paramsMap)
	setClause := setClause(addConfig.Values, &counter, &paramsMap)
	return insertClause, setClause, paramsMap
}

func PrepareDeleteStmt(deleteConfig *defs.AccessConfig) (string, []defs.ParameterRef) {
	counter := uint32(1)
	paramsMap := make([]defs.ParameterRef, 0)
	paramsMap = append(paramsMap, defs.ParameterRef{Index: 0})
	whereClause := prepareFilters(deleteConfig.Filter, &counter, &paramsMap)
	return whereClause, paramsMap
}

func MakeFindQuery(table string, accessConfig *defs.AccessConfig) (string, []defs.ParameterRef) {
	filterClause, paramsMap := PrepareFilters(accessConfig.Filter)
	whereClause := fmt.Sprintf("(1 = 1) AND %s", filterClause)
	return fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(accessConfig.Attributes, ", "), table, whereClause), paramsMap
}

func MakeUpdateQuery(table string, updateConfig *defs.AccessConfig) (string, []defs.ParameterRef) {
	setClause, filterClause, paramsMap := PrepareUpdateStmt(updateConfig)
	whereClause := fmt.Sprintf("(1 = 1) AND %s", filterClause)
	return fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, setClause, whereClause), paramsMap
}

func MakeAddQuery(table string, addConfig *defs.AccessConfig) (string, []defs.ParameterRef) {
	insertClause, paramsMap := PrepareAddStmt(addConfig)
	attributes := make([]string, 0, len(addConfig.Attributes))
	for _, attr := range addConfig.Values {
		attributes = append(attributes, attr.Attribute)
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", table, strings.Join(attributes, ", "), insertClause), paramsMap
}

func MakeAddOrReplaceQuery(table string, addConfig *defs.AccessConfig) (string, []defs.ParameterRef) {
	insertClause, setClause, paramsMap := PrepareAddOrReplaceStmt(addConfig)
	attributes := make([]string, 0, len(addConfig.Attributes))
	for _, attr := range addConfig.Values {
		attributes = append(attributes, attr.Attribute)
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO UPDATE SET %s RETURNING id, (xmax = 0)",
		table, strings.Join(attributes, ", "), insertClause, setClause), paramsMap
}

func MakeDeleteQuery(table string, deleteConfig *defs.AccessConfig) (string, []defs.ParameterRef) {
	filterClause, paramsMap := PrepareDeleteStmt(deleteConfig)
	whereClause := fmt.Sprintf("(1 = 1) AND %s", filterClause)
	return fmt.Sprintf("DELETE FROM %s WHERE %s", table, whereClause), paramsMap
}
