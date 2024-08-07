package datahelpers

import (
	"fmt"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
	"stellarsky.ai/platform/codegen/data-service-generator/config"
	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang"
	"stellarsky.ai/platform/codegen/data-service-generator/db/generator/defs"
)

// Dialect represents a specific SQL dialect
type Dialect interface {
	GetName() string
	FormatIdentifier(name string) string
	GetPlaceholder(index int) string
	DatabaseType(typeId int64) (string, error)
}

// BaseDialect implements common functionality for all dialects
type BaseDialect struct {
	name string
}

func (d *BaseDialect) GetName() string {
	return d.name
}

func (d *BaseDialect) FormatIdentifier(name string) string {
	return fmt.Sprintf("`%s`", strcase.ToSnake(name))
}

func (d *BaseDialect) GetPlaceholder(index int) string {
	return "?"
}

func (d *BaseDialect) DatabaseType(typeId int64) (string, error) {
	return "", fmt.Errorf("DatabaseType not implemented for dialect %s", d.GetName())
}

// Query is the top-level interface for all query types
type Query interface {
	Build() (string, []interface{})
}

// BaseQuery provides common functionality for all query types
type BaseQuery struct {
	dialect Dialect
	parts   []string
	args    []interface{}
}

func (q *BaseQuery) AddPart(part string) {
	q.parts = append(q.parts, part)
}

func (q *BaseQuery) AddArg(arg interface{}) {
	q.args = append(q.args, arg)
}

func (q *BaseQuery) Build() (string, []interface{}) {
	return strings.Join(q.parts, " "), q.args
}

func (q *BaseQuery) Dialect() Dialect {
	return q.dialect
}

type PostgresDialect struct {
	BaseDialect
	placeholderIndex int
}

func NewPostgresDialect() *PostgresDialect {
	return &PostgresDialect{BaseDialect: BaseDialect{name: "postgres"}}
}

func (d *PostgresDialect) GetPlaceholder(index int) string {
	d.placeholderIndex++
	return fmt.Sprintf("$%d", d.placeholderIndex)
}

func (d *PostgresDialect) ResetPlaceholderIndex() {
	d.placeholderIndex = 0
}

func (d *PostgresDialect) DatabaseType(typeId int64) (string, error) {
	return GetPostgresType(typeId), nil
}

type PreparedStmtBuilder struct {
	dialect          Dialect
	modelName        string
	accessConfig     defs.AccessConfig
	placeholderIndex int
	params           []defs.ParameterRef
}

func NewPreparedStmtBuilder(modelName string, accessConfig defs.AccessConfig) *PreparedStmtBuilder {
	return &PreparedStmtBuilder{
		dialect:          NewPostgresDialect(),
		modelName:        modelName,
		accessConfig:     accessConfig,
		placeholderIndex: 0,
		params:           []defs.ParameterRef{},
	}
}

func (psb *PreparedStmtBuilder) resetPlaceholderIndex() {
	psb.placeholderIndex = 0
	psb.params = []defs.ParameterRef{}
}

func (psb *PreparedStmtBuilder) getNextPlaceholder() string {
	psb.placeholderIndex++
	return psb.dialect.GetPlaceholder(psb.placeholderIndex)
}

func (psb *PreparedStmtBuilder) addParam(name string, index int32) {
	psb.params = append(psb.params, defs.ParameterRef{Name: name, Index: index})
}

func (psb *PreparedStmtBuilder) addFuncParam(funcName string, args ...interface{}) {
	psb.params = append(psb.params, defs.ParameterRef{FuncName: funcName, FuncArgs: args})
}

func (psb *PreparedStmtBuilder) buildPreparedStmtFilterClause(filter defs.Filter) string {
	if len(filter.Conditions) > 0 {
		var subclauses []string
		for _, subfilter := range filter.Conditions {
			subclause := psb.buildPreparedStmtFilterClause(subfilter)
			subclauses = append(subclauses, subclause)
		}

		var combinedClause string
		switch strings.ToUpper(filter.Operator) {
		case LogicalAND:
			combinedClause = strings.Join(subclauses, fmt.Sprintf(" %s ", LogicalAND))
		case LogicalOR:
			combinedClause = strings.Join(subclauses, fmt.Sprintf(" %s ", LogicalOR))
		case LogicalNOT:
			combinedClause = fmt.Sprintf("%s(%s)", LogicalNOT, strings.Join(subclauses, fmt.Sprintf(" %s ", LogicalAND)))
		default:
			combinedClause = strings.Join(subclauses, fmt.Sprintf(" %s ", LogicalAND))
		}

		return "(" + combinedClause + ")"
	}

	attr := psb.dialect.FormatIdentifier(filter.Attribute)
	if filter.Transformation != "" {
		attr = fmt.Sprintf("%s(%s)", filter.Transformation, attr)
	}

	switch filter.Operator {
	case OperatorEqual, OperatorNotEqual, OperatorGreaterThan, OperatorLessThan, OperatorGreaterThanOrEqual, OperatorLessThanOrEqual:
		psb.addParam(filter.ParamName, -1)
		return fmt.Sprintf("%s %s %s", attr, filter.Operator, psb.getNextPlaceholder())
	case OperatorIN:
		psb.addParam(filter.ParamName, -1)
		return fmt.Sprintf("%s %s ANY(%s)", attr, OperatorEqual, psb.getNextPlaceholder())
	case OperatorNOTIN:
		psb.addParam(filter.ParamName, -1)
		return fmt.Sprintf("%s %s ALL(%s)", attr, OperatorNotEqual, psb.getNextPlaceholder())
	case OperatorLIKE, OperatorNOTLIKE:
		psb.addParam(filter.ParamName, -1)
		return fmt.Sprintf("%s %s %s", attr, filter.Operator, psb.getNextPlaceholder())
	case OperatorBETWEEN, OperatorNOTBETWEEN:
		psb.addParam(filter.ParamName, 0)
		psb.addParam(filter.ParamName, 1)
		return fmt.Sprintf("(%s %s %s AND %s)", attr, filter.Operator, psb.getNextPlaceholder(), psb.getNextPlaceholder())
	case OperatorIS, OperatorISNOT:
		psb.addParam(filter.ParamName, -1)
		return fmt.Sprintf("%s %s %s", attr, filter.Operator, psb.getNextPlaceholder())
	case OperatorEXISTS, OperatorNOTEXISTS:
		return fmt.Sprintf("%s (%s)", filter.Operator, filter.ParamName)
	default:
		return ""
	}
}

func (psb *PreparedStmtBuilder) buildPreparedStmtSetClause(updates []defs.Update) string {
	var clauses []string

	for _, update := range updates {
		clauses = append(clauses, fmt.Sprintf("%s %s %s",
			psb.dialect.FormatIdentifier(update.Attribute), OperatorEqual, psb.getNextPlaceholder()))
		psb.addParam(update.ParamName, -1)
	}

	return strings.Join(clauses, ", ")
}

func (psb *PreparedStmtBuilder) buildWhereClause() string {
	if len(psb.accessConfig.Filter) == 0 {
		return ""
	}

	whereClause := psb.buildPreparedStmtFilterClause(defs.Filter{
		Operator:   LogicalAND,
		Conditions: psb.accessConfig.Filter,
	})

	return fmt.Sprintf("%s %s", KeywordWHERE, whereClause)
}

func (psb *PreparedStmtBuilder) BuildFindPreparedStmt() (string, []defs.ParameterRef) {
	psb.resetPlaceholderIndex()
	query := strings.Builder{}

	query.WriteString(KeywordSELECT + " ")
	if len(psb.accessConfig.Attributes) > 0 {
		attributeTokens := make([]string, len(psb.accessConfig.Attributes))
		for i, attr := range psb.accessConfig.Attributes {
			attributeTokens[i] = psb.dialect.FormatIdentifier(attr)
		}
		query.WriteString(strings.Join(attributeTokens, ", "))
	} else {
		query.WriteString("*")
	}
	query.WriteString(fmt.Sprintf(" %s %s", KeywordFROM, psb.dialect.FormatIdentifier(psb.modelName)))

	whereClause := psb.buildWhereClause()
	if whereClause != "" {
		query.WriteString(" " + whereClause)
	}

	return query.String(), psb.params
}

func (psb *PreparedStmtBuilder) BuildDeletePreparedStmt() (string, []defs.ParameterRef) {
	psb.resetPlaceholderIndex()
	query := strings.Builder{}

	query.WriteString(fmt.Sprintf("%s %s %s", KeywordDELETE, KeywordFROM, psb.dialect.FormatIdentifier(psb.modelName)))

	whereClause := psb.buildWhereClause()
	if whereClause != "" {
		query.WriteString(" " + whereClause)
	}

	return query.String(), psb.params
}

func (psb *PreparedStmtBuilder) BuildUpdatePreparedStmt() (string, []defs.ParameterRef) {
	psb.resetPlaceholderIndex()
	query := strings.Builder{}

	query.WriteString(fmt.Sprintf("%s %s %s ", KeywordUPDATE, psb.dialect.FormatIdentifier(psb.modelName), KeywordSET))

	setClauses := []string{}

	// Regular SET clauses
	if len(psb.accessConfig.Set) > 0 {
		setClauses = append(setClauses, psb.buildPreparedStmtSetClause(psb.accessConfig.Set))
	}

	// Autoincrement columns
	for _, attr := range psb.accessConfig.Autoincrement {
		setClauses = append(setClauses, fmt.Sprintf("%s = %s + 1", psb.dialect.FormatIdentifier(attr), psb.dialect.FormatIdentifier(attr)))
	}

	// Capture timestamp columns
	for _, attr := range psb.accessConfig.CaptureTimestamp {
		setClauses = append(setClauses, fmt.Sprintf("%s = NOW()", psb.dialect.FormatIdentifier(attr)))
	}

	query.WriteString(strings.Join(setClauses, ", "))

	whereClause := psb.buildWhereClause()
	if whereClause != "" {
		query.WriteString(" " + whereClause)
	}

	return query.String(), psb.params
}

func (psb *PreparedStmtBuilder) BuildAddPreparedStmt() (string, []defs.ParameterRef) {
	psb.resetPlaceholderIndex()
	query := strings.Builder{}

	attributes := make([]string, 0, len(psb.accessConfig.Values))
	values := make([]string, 0, len(psb.accessConfig.Values))

	attributes = append(attributes, psb.dialect.FormatIdentifier("id"))
	values = append(values, psb.getNextPlaceholder())
	psb.addFuncParam("UUIDV7")

	for _, update := range psb.accessConfig.Values {
		attributes = append(attributes, psb.dialect.FormatIdentifier(update.Attribute))
		values = append(values, psb.getNextPlaceholder())
		psb.addParam(update.ParamName, -1)
	}

	query.WriteString(fmt.Sprintf("%s %s %s (%s) %s (%s) %s id",
		KeywordINSERT, KeywordINTO,
		psb.dialect.FormatIdentifier(psb.modelName),
		strings.Join(attributes, ", "),
		KeywordVALUES,
		strings.Join(values, ", "),
		KeywordRETURNING))

	return query.String(), psb.params
}

func (psb *PreparedStmtBuilder) BuildAddOrReplacePreparedStmt() (string, []defs.ParameterRef) {
	psb.resetPlaceholderIndex()
	query := strings.Builder{}

	attributes := make([]string, 0, len(psb.accessConfig.Values))
	values := make([]string, 0, len(psb.accessConfig.Values))
	updateClauses := make([]string, 0, len(psb.accessConfig.Values))

	attributes = append(attributes, psb.dialect.FormatIdentifier("id"))
	values = append(values, psb.getNextPlaceholder())
	psb.addFuncParam("UUIDV7")
	for _, update := range psb.accessConfig.Values {
		attr := psb.dialect.FormatIdentifier(update.Attribute)
		placeholder := psb.getNextPlaceholder()

		attributes = append(attributes, attr)
		values = append(values, placeholder)
		psb.addParam(update.ParamName, -1)
	}

	for _, update := range psb.accessConfig.Values {
		attr := psb.dialect.FormatIdentifier(update.Attribute)
		placeholder := psb.getNextPlaceholder()
		updateClauses = append(updateClauses, fmt.Sprintf("%s = %s", attr, placeholder))
		psb.addParam(update.ParamName, -1)
	}

	query.WriteString(fmt.Sprintf("%s %s %s (%s) %s (%s) %s %s %s %s %s %s %s id, (xmax = 0) AS inserted",
		KeywordINSERT, KeywordINTO, psb.dialect.FormatIdentifier(psb.modelName),
		strings.Join(attributes, ", "),
		KeywordVALUES,
		strings.Join(values, ", "),
		KeywordON, KeywordCONFLICT, KeywordDO, KeywordUPDATE, KeywordSET,
		strings.Join(updateClauses, ", "),
		KeywordRETURNING))

	return query.String(), psb.params
}

type SchemaBuilder struct {
	dialect      Dialect
	databaseName string
	dataConfig   defs.DataConfig
}

func NewSchemaBuilder(dialect Dialect, databaseName string, config defs.DataConfig) *SchemaBuilder {
	return &SchemaBuilder{
		dialect:      dialect,
		databaseName: databaseName,
		dataConfig:   config,
	}
}

func (sb *SchemaBuilder) BuildCreateDatabase() string {
	return fmt.Sprintf("CREATE DATABASE %s;", sb.dialect.FormatIdentifier(sb.dataConfig.FamilyName))
}

func (sb *SchemaBuilder) BuildCreateUser() string {
	return fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s';", sb.dataConfig.DatabaseConfig.UserName, sb.dataConfig.DatabaseConfig.Password)
}

func (sb *SchemaBuilder) BuildGrantPermissions() string {
	return fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s;",
		sb.dialect.FormatIdentifier(sb.dataConfig.FamilyName),
		sb.dialect.FormatIdentifier(sb.dataConfig.DatabaseConfig.UserName))
}

func (sb *SchemaBuilder) BuildCreateTable(model *defs.ModelConfig) string {
	columns := []string{
		fmt.Sprintf("%s`id` UUID PRIMARY KEY DEFAULT gen_random_uuid()", golang.Indent),
		fmt.Sprintf("%s`version` INTEGER NOT NULL DEFAULT 1", golang.Indent),
		fmt.Sprintf("%s`updated_at` TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP", golang.Indent),
	}

	// Add columns from ModelConfig
	for _, attr := range model.Attributes {
		attrName, attrType := sb.mapAttributeTypeToSQL(attr)
		columns = append(columns, fmt.Sprintf("%s%s %s", golang.Indent, sb.dialect.FormatIdentifier(attrName), attrType))
	}

	// Create table SQL
	createTableSQL := fmt.Sprintf("CREATE TABLE %s (\n%s\n);",
		sb.dialect.FormatIdentifier(model.Name),
		strings.Join(columns, ",\n"))

	// Add indexes based on filters
	seenIndexes := sb.generateIndexesFromFilters(model.GetAllFilters())
	sb.generateIndexes(model.Model.GetIndexes(), seenIndexes)
	indexSQL := sb.generateIndexSQL(model.Name, seenIndexes)
	return createTableSQL + "\n\n" + indexSQL + "\n"
}

func (sb *SchemaBuilder) mapAttributeTypeToSQL(attrId int64) (string, string) {
	// Map attribute types to SQL types
	// Implement this based on your specific type mappings
	attribute, ok := config.Attributes[attrId]
	if !ok {
		return "", ""
	}
	attrType, _ := sb.dialect.DatabaseType(attribute.TypeId)
	return attribute.Name, attrType

}

type indexItem struct {
	joinedAttrs string // e.g. "attr1, attr2, attr3"
	isUnique    bool   // e.g. "UNIQUE INDEX" or "INDEX"
}

func (sb *SchemaBuilder) newIndexItem(attrs []string, isUnique bool) indexItem {
	snakeAttrs := []string{}
	for _, attr := range attrs {
		snakeAttrs = append(snakeAttrs, sb.dialect.FormatIdentifier(strcase.ToSnake(attr)))
	}
	return indexItem{joinedAttrs: strings.Join(snakeAttrs, ", "), isUnique: isUnique}
}

func (sb *SchemaBuilder) generateIndexesFromFilters(allFilters []defs.Filter) map[string]indexItem {
	seenIndexes := map[string]indexItem{}

	// we go one level
	oneLevelFilters := allFilters[:]
	for _, filter := range oneLevelFilters {
		if filter.Attribute != "" {
			oneLevelFilters = append(oneLevelFilters, filter)
		}

		for _, condition := range filter.Conditions {
			if condition.Attribute != "" {
				oneLevelFilters = append(oneLevelFilters, condition)
			}
		}
	}

	for _, filter := range oneLevelFilters {
		if filter.Attribute != "" {
			item := sb.newIndexItem([]string{filter.Attribute}, false)
			if _, ok := seenIndexes[filter.Attribute]; !ok {
				seenIndexes[filter.Attribute] = item
			}
		}
	}

	return seenIndexes
}

func (sb *SchemaBuilder) generateIndexes(allIndexes []defs.Index, seenIndexes map[string]indexItem) {

	for _, index := range allIndexes {
		indexAttrs := make([]string, len(index.Attributes))
		for i, attr := range index.Attributes {
			attrName := config.Attributes[attr].Name
			indexAttrs[i] = strcase.ToSnake(attrName)
		}
		indexAttr := strings.Join(indexAttrs, ", ")

		if item, ok := seenIndexes[indexAttr]; !ok || (!item.isUnique && index.IsUnique) {
			seenIndexes[index.IndexName] = sb.newIndexItem(indexAttrs, index.IsUnique)
		}
	}

}

func (sb *SchemaBuilder) generateIndexSQL(modelName string, indexItemMap map[string]indexItem) string {
	var indexSQLs []string
	for _, item := range indexItemMap {
		indexType := "INDEX"
		if item.isUnique {
			indexType = "UNIQUE INDEX"
		}
		indexSQL := fmt.Sprintf("CREATE %s ON %s (%s);",
			indexType,
			sb.dialect.FormatIdentifier(modelName),
			item.joinedAttrs)
		indexSQLs = append(indexSQLs, indexSQL)
	}
	slices.Sort(indexSQLs)
	return strings.Join(indexSQLs, "\n")
}
