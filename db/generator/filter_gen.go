package generator

import (
	"fmt"
	"strings"

	"stellarsky.ai/platform/codegen/data-service-generator/base"
)

type Filter struct {
	Attribute      string   `yaml:"attribute,omitempty"`
	Transformation string   `yaml:"transformation,omitempty"`
	Operator       string   `yaml:"operator"`
	ParamName      string   `yaml:"param_name,omitempty"`
	Conditions     []Filter `yaml:"conditions,omitempty"`
}

type Filters struct {
	Filters []Filter `yaml:"filter"`
}

type ParameterRef struct {
	Name  string `yaml:"name"`
	Index int32  `yaml:"index"`
}

func applyTransformation(attribute string, transformation string) string {
	if transformation == "" {
		return attribute
	}
	return fmt.Sprintf("%s(%s)", transformation, attribute)
}

// func ParseFilter(filter Filter, prepare bool, counter uint32) (string, uint32) {
// 	attribute := applyTransformation(filter.Attribute, filter.Transformation)
// 	switch strings.ToUpper(filter.Operator) {
// 	case OperatorEquals, OperatorNotEquals, OperatorLessThan, OperatorLessThanEquals, OperatorGreaterThan, OperatorGreaterThanEquals:
// 		result := fmt.Sprintf("%s %s %v", attribute, filter.Operator, generateValueByType(filter.Value, prepare, counter))
// 		counter += 1
// 		return result, counter
// 	case OperatorBetween:
// 		values := filter.Value.([]interface{})
// 		result := fmt.Sprintf("%s BETWEEN %v AND %v", attribute, generateValueByType(values[0], prepare, counter), generateValueByType(values[1], prepare, counter+1))
// 		counter += 2
// 		return result, counter
// 	case OperatorIn:
// 		if prepare {
// 			var valueFormat string
// 			valueFormat, counter = makePreparedCounter(counter)
// 			return fmt.Sprintf("%s IN (%s)", attribute, valueFormat), counter
// 		} else {
// 			values := filter.Value.([]interface{})
// 			valueStrs := make([]string, len(values))
// 			for i, val := range values {
// 				valueStrs[i] = generateValueByType(val, prepare, counter)
// 				counter += 1
// 			}
// 			return fmt.Sprintf("%s IN (%s)", attribute, strings.Join(valueStrs, ", ")), counter
// 		}
// 	case LogicalAnd, LogicalOr, LogicalNot:
// 		subConditions := make([]string, len(filter.Conditions))
// 		for i, subFilter := range filter.Conditions {
// 			subConditions[i], counter = ParseFilter(subFilter, prepare, counter)
// 		}
// 		if filter.Operator == LogicalNot {
// 			return fmt.Sprintf("NOT(%s)", subConditions[0]), counter
// 		}
// 		return fmt.Sprintf("(%s)", strings.Join(subConditions, fmt.Sprintf(" %s ", filter.Operator))), counter
// 	default:
// 		return "", counter
// 	}
// }

// func retriveValuesList(filter Filter, values []interface{}) []interface{} {
// 	switch strings.ToUpper(filter.Operator) {
// 	case OperatorEquals, OperatorNotEquals, OperatorLessThan, OperatorLessThanEquals, OperatorGreaterThan, OperatorGreaterThanEquals:
// 		values = append(values, filter.Value)
// 		return values
// 	case OperatorBetween:
// 		values = append(values, filter.Value.([]interface{})...)
// 		return values
// 	case OperatorIn:
// 		values = append(values, filter.Value.([]interface{})...)
// 		return values
// 	case LogicalAnd, LogicalOr, LogicalNot:
// 		for _, subFilter := range filter.Conditions {
// 			values = retriveValuesList(subFilter, values)
// 		}
// 		return values
// 	default:
// 		return values

// 	}
// }

func buildPrepareStmt(filter Filter, counter *uint32, paramsMap *[]ParameterRef) string {
	attribute := applyTransformation(filter.Attribute, filter.Transformation)
	switch strings.ToUpper(filter.Operator) {
	case OperatorEquals, OperatorNotEquals, OperatorLessThan, OperatorLessThanEquals, OperatorGreaterThan, OperatorGreaterThanEquals:
		result := fmt.Sprintf("%s %s %v", attribute, filter.Operator, makePreparedCounter(counter))
		*paramsMap = append(*paramsMap, ParameterRef{
			Name:  filter.ParamName,
			Index: -1,
		})
		return result
	case OperatorBetween:
		lowEnd := makePreparedCounter(counter)
		highEnd := makePreparedCounter(counter)
		*paramsMap = append(*paramsMap, ParameterRef{Name: filter.ParamName, Index: 0}, ParameterRef{Name: filter.ParamName, Index: 1})
		result := fmt.Sprintf("(%s BETWEEN %v AND %v)", attribute, lowEnd, highEnd)
		return result
	case OperatorIn:
		result := fmt.Sprintf("%s = ANY(%s)", attribute, makePreparedCounter(counter))
		*paramsMap = append(*paramsMap, ParameterRef{
			Name:  filter.ParamName,
			Index: -1,
		})
		return result
	case LogicalAnd, LogicalOr, LogicalNot:
		subConditions := make([]string, len(filter.Conditions))
		for i, subFilter := range filter.Conditions {
			subConditions[i] = buildPrepareStmt(subFilter, counter, paramsMap)
		}
		if filter.Operator == LogicalNot {
			return fmt.Sprintf("NOT(%s)", subConditions[0])
		}
		return fmt.Sprintf("(%s)", strings.Join(subConditions, fmt.Sprintf(" %s ", filter.Operator)))
	default:
		return ""
	}
}

func makePreparedCounter(counter *uint32) string {
	result := fmt.Sprintf("$%d", *counter)
	*counter += 1
	return result
}

func generateValueByType(value interface{}, prepare bool, counter uint32) string {
	if prepare {
		return fmt.Sprintf("$%d", counter)
	}

	var valueToFormat string
	switch value.(type) {
	case string:
		valueToFormat = fmt.Sprintf("'%v'", value)
	case int, int8, int16, int32, int64, uint, uint16, uint32, uint64, bool, float32, float64:
		valueToFormat = fmt.Sprintf("%v", value)
	default:
		valueToFormat = fmt.Sprintf("%v", value)
	}
	base.LOG.Debug("GenerateValueByType", "value", value, "valueToFormat", valueToFormat)
	return valueToFormat
}

// func BuildWhereCondition(filters Filters) string {
// 	return BuildWhereFilters(filters.Filters, false)
// }

// func BuildWhereFilters(filters []Filter, prepared bool) string {
// 	conditions := make([]string, len(filters))
// 	counter := uint32(1)
// 	for i, filter := range filters {
// 		conditions[i], counter = ParseFilter(filter, prepared, counter)
// 	}
// 	return strings.Join(conditions, " AND ")
// }

// Returns a prepared condition in Postgresql format, example "column = $1 AND column2 = $2"
// and returns id to param name mapping in array, example ["", "column" "column2"], keeping first string empty to start with $1
func PrepareFilters(filters []Filter) (string, []ParameterRef) {
	conditions := make([]string, len(filters))
	counter := uint32(1)
	paramsMap := make([]ParameterRef, 0)
	paramsMap = append(paramsMap, ParameterRef{Index: 0})
	for i, filter := range filters {
		conditions[i] = buildPrepareStmt(filter, &counter, &paramsMap)
	}
	return fmt.Sprintf("(%s)", strings.Join(conditions, " AND ")), paramsMap
}

// func ReadValuesList(filters []Filter) []interface{} {
// 	result := make([]interface{}, 0)
// 	for _, filter := range filters {
// 		values := make([]interface{}, 0)
// 		values = retriveValuesList(filter, values)
// 		result = append(result, values...)
// 	}
// 	return result
// }
