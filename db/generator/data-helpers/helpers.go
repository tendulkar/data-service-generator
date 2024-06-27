package datahelpers

import (
	"database/sql"

	"stellarsky.ai/platform/codegen/data-service-generator/db/generator/defs"
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
