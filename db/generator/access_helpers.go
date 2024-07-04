package generator

// import (
// 	"context"
// 	"database/sql"

// 	"stellarsky.ai/platform/codegen/data-service-generator/db/generator/defs"
// )

// type ModelDB struct {
// 	modelName     string
// 	dbConn        *sql.DB
// 	preparedCache map[string]*sql.Stmt
// 	preparedMap   map[string][]defs.ParameterRef
// }

// func findExec[T any, M any, RQ any, RS any](name string, ctx context.Context, db *ModelDB, request *RQ) ([]M, error) {
// 	stmt := db.preparedCache[name]
// 	rows, err := stmt.Query(QueryArgs(request))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	stmt.Exec()
// 	var results []M
// 	for rows.Next() {
// 		var item M
// 		err := rows.Scan(&item)
// 		if err != nil {
// 			return nil, err
// 		}
// 		results = append(results, item)
// 	}
// 	return results, nil
// }
