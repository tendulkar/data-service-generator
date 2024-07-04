package datahelpers

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"stellarsky.ai/platform/codegen/data-service-generator/config"
	"stellarsky.ai/platform/codegen/data-service-generator/db/generator/defs"
	"stellarsky.ai/platform/codegen/data-service-generator/db/models"
)

type PreparedInfo struct {
	stmt     *sql.Stmt
	paramMap []defs.ParameterRef
}

type ModelDB struct {
	db         *sql.DB
	statements map[string]PreparedInfo
}

func NewModelDB(db *sql.DB) *ModelDB {
	return &ModelDB{
		db:         db,
		statements: make(map[string]PreparedInfo),
	}
}

func prepareAndCacheQuery(model *ModelDB, queryName string, query string, paramMap []defs.ParameterRef) error {
	stmt, err := model.db.Prepare(query)
	if err != nil {
		return err
	}

	model.statements[queryName] = PreparedInfo{
		stmt:     stmt,
		paramMap: paramMap,
	}
	return nil
}

// Given Access configuration we need to build prepared statements for it,
// Need to build different prepared statements for different access methods, find, update, add, add_or_replace, delete
// Each prepared statement will have different parameter mappings,
func PrepareStatements(model *ModelDB, access *defs.Access) error {

	for _, accessConfig := range access.Find {
		filter := accessConfig.Filter
		preparedQuery, paramMap := PrepareFilters(filter)
		err := prepareAndCacheQuery(model, accessConfig.Name, preparedQuery, paramMap)
		if err != nil {
			return err
		}
	}

	for _, accessConfig := range access.Update {
		filter := accessConfig.Filter
		preparedQuery, paramMap := PrepareFilters(filter)
		err := prepareAndCacheQuery(model, accessConfig.Name, preparedQuery, paramMap)
		if err != nil {
			return err
		}
	}

	for _, accessConfig := range access.Add {
		filter := accessConfig.Filter
		preparedQuery, paramMap := PrepareFilters(filter)
		err := prepareAndCacheQuery(model, accessConfig.Name, preparedQuery, paramMap)
		if err != nil {
			return err
		}
	}

	for _, accessConfig := range access.AddOrReplace {
		filter := accessConfig.Filter
		preparedQuery, paramMap := PrepareFilters(filter)
		err := prepareAndCacheQuery(model, accessConfig.Name, preparedQuery, paramMap)
		if err != nil {
			return err
		}
	}

	for _, accessConfig := range access.Delete {
		filter := accessConfig.Filter
		preparedQuery, paramMap := PrepareFilters(filter)
		err := prepareAndCacheQuery(model, accessConfig.Name, preparedQuery, paramMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadParamValues(params []defs.ParameterRef, request map[string]interface{}) ([]interface{}, error) {
	values := []interface{}{}
	for _, param := range params {
		paramName := param.Name
		paramIndex := param.Index

		if paramIndex == -1 {
			if value, ok := request[paramName]; ok {
				values = append(values, value)
			} else {
				return nil, fmt.Errorf("parameter %s not found in request", paramName)
			}
		} else {
			if value, ok := request[paramName].([]interface{}); ok {
				values = append(values, value[paramIndex])
			} else {
				return nil, fmt.Errorf("parameter %s not found in request", paramName)
			}
		}
	}
	return values, nil
}

func GetPostgresType(typeId int64) string {
	return config.PostgresTypeMaps[typeId].MappedType
}

func GetValidations(validationIds []int64) []*models.Validation {
	validations := []*models.Validation{}

	for _, validationId := range validationIds {
		validation, ok := config.Validations[validationId]
		if ok {
			validations = append(validations, &validation)
		}
	}
	return validations
}

func PostgresToGoType(pgType string) string {
	// Remove dimensions, e.g., converting varchar(255) to varchar
	dimensionRegex := regexp.MustCompile(`\(\d+(\s*,\s*\d+)*\)`)
	pgType = dimensionRegex.ReplaceAllString(pgType, "")
	// base.LOG.Info("Postgres type", "type", pgType)
	// Normalize the input to handle cases insensitively
	normalizedType := strings.ToLower(pgType)

	// Check for array types
	isArray := strings.HasPrefix(normalizedType, "[]")
	if isArray {
		normalizedType = normalizedType[2:] // Remove the array prefix
	}

	// Determine the Go type based on the PostgreSQL type
	var goType string
	switch normalizedType {
	case "bigint":
		goType = "int64"
	case "integer", "int":
		goType = "int"
	case "smallint":
		goType = "int16"
	case "boolean":
		goType = "bool"
	case "real":
		goType = "float32"
	case "double precision", "numeric", "decimal":
		goType = "float64"
	case "char", "varchar", "text", "uuid":
		goType = "string"
	case "bytea":
		goType = "[]byte"
	case "date", "time", "timestamp", "timestamptz":
		goType = "time.Time"
	case "json", "jsonb":
		goType = "interface{}"
	default:
		goType = "interface{}" // Default case for types that don't have a direct mapping
	}

	// Return the Go type, with slice notation if it's an array
	if isArray {
		return "[]" + goType
	}
	return goType
}
