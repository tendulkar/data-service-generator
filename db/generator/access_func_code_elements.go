package generator

import (
	"fmt"

	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang"
	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang/goutils"
	"stellarsky.ai/platform/codegen/data-service-generator/db/generator/defs"
)

func lookupStmtCE(confName string, dbName, cacheName, stmtName string) *golang.MapLookup {
	return &golang.MapLookup{
		NewOutput: stmtName,
		Receiver:  dbName,
		Name:      cacheName,
		Key:       confName,
	}
}

func parseParamsCE(confName string, requestParamsName string, valuesName string, returnParams []*golang.Parameter) *golang.FunctionCall {
	return &golang.FunctionCall{
		NewOutput: []string{valuesName, "err"},
		Function:  fmt.Sprintf("%sParseParams", confName),
		Args:      []string{requestParamsName},
		ErrorHandler: &golang.ErrorHandler{
			ErrorFunctionReturns: returnParams,
		},
	}
}

func queryStmtCE(stmtName, valuesName, rowsName string, returnParams []*golang.Parameter) *golang.FunctionCall {
	return &golang.FunctionCall{
		NewOutput: []string{rowsName, "err"},
		Receiver:  stmtName,
		Function:  "Query",
		Args:      []string{fmt.Sprintf("%s...", valuesName)},
		ErrorHandler: &golang.ErrorHandler{
			ErrorFunctionReturns: returnParams,
		},
		CleanningHandler: &golang.CleanningHandler{
			Receiver: rowsName,
			Function: "Close",
		},
	}
}

func queryRowStmtCE(stmtName string, args []string, valuesName string, returnParams []*golang.Parameter) *golang.FunctionCall {
	return &golang.FunctionCall{
		NewOutput: []string{"queryErr"},
		Receiver:  fmt.Sprintf("%s.QueryRow(%s...)", stmtName, valuesName),
		Function:  "Scan",
		Args:      args,
		ErrorHandler: &golang.ErrorHandler{
			Error:                "queryErr",
			ErrorFunctionReturns: returnParams,
		},
	}
}

func execStmtCE(stmtName, valuesName, resultName string, returnParams []*golang.Parameter) *golang.FunctionCall {
	return &golang.FunctionCall{
		NewOutput: []string{resultName, "err"},
		Receiver:  stmtName,
		Function:  "Exec",
		Args:      []string{fmt.Sprintf("%s...", valuesName)},
		ErrorHandler: &golang.ErrorHandler{
			ErrorFunctionReturns: returnParams,
		},
	}
}

func createVarCE(name string, typ string) *golang.Variable {
	return &golang.Variable{
		Names: name,
		Type:  typ,
	}
}

func appendCE(arr string, value string) *golang.FunctionCall {
	return &golang.FunctionCall{
		Output:   arr,
		Function: "append",
		Args:     []string{arr, value},
	}
}

func returnValuesCE(values ...string) *golang.CodeElement {
	return &golang.CodeElement{
		Return: values,
	}
}

func returnResultNilCE(results string) *golang.CodeElement {
	return returnValuesCE(results, "nil")
}

func attributeRefArgs(itemName string, attributes []string) []string {
	args := make([]string, 0, len(attributes))
	for _, attr := range attributes {
		args = append(args, fmt.Sprintf("&%s.%s", itemName, attr))
	}
	return args
}

func scanRowCE(itemName string, attributes []string, rowsName string) *golang.FunctionCall {
	return &golang.FunctionCall{
		NewOutput: []string{"scanErr"},
		Receiver:  rowsName,
		Function:  "Scan",
		Args:      attributeRefArgs(itemName, attributes),
		ErrorHandler: &golang.ErrorHandler{
			Error:        "scanErr",
			ErrorReturns: []string{"nil", "scanErr"},
		},
	}
}

func scanResultsCE(modelName string, attributes []string, resultsName string) *golang.RepeatByCondition {
	return &golang.RepeatByCondition{
		Condition: &golang.CodeElement{
			FunctionCall: &golang.FunctionCall{
				Function: "Next",
				Receiver: "rows",
			},
		},
		Body: golang.CodeElements{
			{
				Variable: createVarCE("item", modelName),
			},
			{
				FunctionCall: scanRowCE("item", attributes, "rows"),
			},
			{
				FunctionCall: appendCE(resultsName, "item"),
			},
		},
	}
}

func ctxParamCE(ctxName string) *golang.Parameter {
	return &golang.Parameter{
		Name: ctxName,
		Type: &golang.GoType{
			Name: "context.Context",
		},
	}
}

func dbParamCE(dbName string) *golang.Parameter {
	return &golang.Parameter{
		Name: dbName,
		Type: &golang.GoType{
			Name:   "*sql.DB",
			Source: "database/sql",
		},
	}
}

func modelDBParamCE(dbName, modelDBName string) *golang.Parameter {
	return &golang.Parameter{
		Name: dbName,
		Type: &golang.GoType{
			Name: fmt.Sprintf("*%s", modelDBName),
		},
	}
}

func requestParamsCE(confName string, requestParamsName string) *golang.Parameter {
	return &golang.Parameter{
		Name: requestParamsName,
		Type: &golang.GoType{
			Name: fmt.Sprintf("%sParams", confName),
		},
	}
}

func dbRequestParamsCE(dbName, modelDBName, name, requestParamsName string) []*golang.Parameter {
	return []*golang.Parameter{
		modelDBParamCE(dbName, modelDBName),
		requestParamsCE(name, requestParamsName),
	}
}

func ctxDBRequestParamsCE(ctxName, dbName, modelDBName, name, requestParamsName string) []*golang.Parameter {
	params := []*golang.Parameter{
		ctxParamCE(ctxName),
	}
	params = append(params, dbRequestParamsCE(dbName, modelDBName, name, requestParamsName)...)
	return params
}

func errorParamCE(errorName string) *golang.Parameter {
	return &golang.Parameter{
		Name: errorName,
		Type: &golang.GoType{
			Name: "error",
		},
	}
}

func resultsParamCE(modelName, resultsName, sourceName string) *golang.Parameter {
	return &golang.Parameter{
		Name: resultsName,
		Type: &golang.GoType{
			Name:   fmt.Sprintf("[]%s", modelName),
			Source: sourceName,
		},
	}
}

func resultsErrorParamsCE(modelName, resultsName, errorName, sourceName string) []*golang.Parameter {
	return []*golang.Parameter{
		resultsParamCE(modelName, resultsName, sourceName),
		errorParamCE(errorName),
	}
}

func typeOnlyParamsCE(types ...string) []*golang.Parameter {
	params := make([]*golang.Parameter, 0, len(types))
	for _, typ := range types {
		params = append(params, &golang.Parameter{
			Type: &golang.GoType{
				Name: typ,
			},
		})
	}
	return params
}

func callResultErrorCE(objName, funcName string, args []string, resultName string, errorName string, returnParams []*golang.Parameter) *golang.FunctionCall {
	return &golang.FunctionCall{
		NewOutput: []string{resultName, errorName},
		Receiver:  objName,
		Function:  funcName,
		Args:      args,
		ErrorHandler: &golang.ErrorHandler{
			Error:                errorName,
			ErrorFunctionReturns: returnParams,
		},
	}
}

func callRowsAffectedCE(objName, resultName, errorName string, returnParams []*golang.Parameter) *golang.FunctionCall {
	return callResultErrorCE(objName, "RowsAffected", []string{}, resultName, errorName, returnParams)
}

func FindCodeFunction(modelName, modelDBName, name string, attributes []string) *golang.FunctionDef {
	fnReturns := resultsErrorParamsCE(modelName, "results", "err", "")
	codeElems := golang.CodeElements{
		{
			MapLookup: lookupStmtCE(name, "db", "preparedCache", "stmt"),
		},
		{
			FunctionCall: parseParamsCE(name, "requestParams", "values", fnReturns),
		},
		{
			FunctionCall: queryStmtCE("stmt", "values", "rows", fnReturns),
		},
		{
			Variable: createVarCE("results", fmt.Sprintf("[]%s", modelName)),
		},
		{
			RepeatCond: scanResultsCE(modelName, attributes, "results"),
		},
		{
			Return: []string{"results", "nil"},
		},
	}

	params := ctxDBRequestParamsCE("ctx", "db", modelDBName, name, "requestParams")
	dependencies := []golang.Dependency{}

	fn := &golang.FunctionDef{
		Name:         name,
		Parameters:   params,
		Body:         codeElems,
		Returns:      fnReturns,
		Dependencies: dependencies,
	}

	return fn
}

func UpdateCodeFunction(name string, modelDBName string) *golang.FunctionDef {
	fnReturns := typeOnlyParamsCE("int64", "error")

	codeElems := golang.CodeElements{
		{
			MapLookup: lookupStmtCE(name, "db", "preparedCache", "stmt"),
		},
		{
			FunctionCall: parseParamsCE(name, "requestParams", "values", fnReturns),
		},
		{
			FunctionCall: execStmtCE("stmt", "values", "result", fnReturns),
		},
		{
			FunctionCall: callRowsAffectedCE("result", "rowsAffected", "err", fnReturns),
		},
		{
			Return: []string{"rowsAffected", "nil"},
		},
	}
	params := ctxDBRequestParamsCE("ctx", "db", modelDBName, name, "requestParams")
	dependencies := make([]golang.Dependency, 0)
	fn := &golang.FunctionDef{
		Name:         name,
		Parameters:   params,
		Body:         codeElems,
		Returns:      fnReturns,
		Dependencies: dependencies,
	}
	return fn
}

func AddCodeFunction(name string, modelDBName string) *golang.FunctionDef {
	fnReturns := typeOnlyParamsCE("int64", "error")
	codeElems := golang.CodeElements{
		{
			MapLookup: lookupStmtCE(name, "db", "preparedCache", "stmt"),
		},
		{
			FunctionCall: parseParamsCE(name, "requestParams", "values", fnReturns),
		},
		{
			Variable: createVarCE("id", "int64"),
		},
		{
			FunctionCall: queryRowStmtCE("stmt", []string{"&id"}, "values", fnReturns),
		},
		{
			Return: []string{"id", "nil"},
		},
	}
	params := ctxDBRequestParamsCE("ctx", "db", modelDBName, name, "requestParams")
	fn := &golang.FunctionDef{
		Name:         name,
		Parameters:   params,
		Body:         codeElems,
		Returns:      fnReturns,
		Dependencies: nil,
	}
	return fn
}

func AddOrReplaceCodeFunction(name string, modelDBName string) *golang.FunctionDef {
	fnReturns := typeOnlyParamsCE("int64", "bool", "error")
	codeElems := golang.CodeElements{
		{
			MapLookup: lookupStmtCE(name, "db", "preparedCache", "stmt"),
		},
		{
			FunctionCall: parseParamsCE(name, "requestParams", "values", fnReturns),
		},
		{
			Variable: createVarCE("id", "int64"),
		},
		{
			Variable: createVarCE("inserted", "bool"),
		},
		{
			FunctionCall: queryRowStmtCE("stmt", []string{"&id", "&inserted"}, "values", fnReturns),
		},
		{
			Return: []string{"id", "inserted", "nil"},
		},
	}
	params := ctxDBRequestParamsCE("ctx", "db", modelDBName, name, "requestParams")
	fn := &golang.FunctionDef{
		Name:         name,
		Parameters:   params,
		Body:         codeElems,
		Returns:      fnReturns,
		Dependencies: nil,
	}
	return fn
}

func DeleteCodeFunction(name string, modelDBName string) *golang.FunctionDef {
	fnReturns := typeOnlyParamsCE("int64", "error")
	codeElems := golang.CodeElements{
		{
			MapLookup: lookupStmtCE(name, "db", "preparedCache", "stmt"),
		},
		{
			FunctionCall: parseParamsCE(name, "requestParams", "values", fnReturns),
		},
		{
			FunctionCall: execStmtCE("stmt", "values", "result", fnReturns),
		},
		{
			FunctionCall: callRowsAffectedCE("result", "rowsAffected", "err", fnReturns),
		},
		{
			Return: []string{"rowsAffected", "nil"},
		},
	}
	params := ctxDBRequestParamsCE("ctx", "db", modelDBName, name, "requestParams")
	fn := &golang.FunctionDef{
		Name:         name,
		Parameters:   params,
		Body:         codeElems,
		Returns:      fnReturns,
		Dependencies: nil,
	}
	return fn
}

func ReadParamsFunction(paramRefs []defs.ParameterRef, confName string,
	valuesName string, paramsName string) *golang.FunctionDef {
	body := golang.CodeElements{
		{
			Variable: createVarCE(valuesName, "[]interface{}"),
		},
	}

	for _, paramRef := range paramRefs {
		paramArg := fmt.Sprintf("%s.%s", paramsName, paramRef.Name)
		if paramRef.Index != -1 {
			paramArg = fmt.Sprintf("%s.([]interface{})[%d]", paramArg, paramRef.Index)
		}
		body = append(body, &golang.CodeElement{
			FunctionCall: appendCE(valuesName, paramArg),
		})
	}

	body = append(body, returnResultNilCE(valuesName))

	fnParams := []*golang.Parameter{
		{
			Name: paramsName,
			Type: &golang.GoType{
				Name: fmt.Sprintf("%sParams", confName),
			},
		},
	}

	fnName := fmt.Sprintf("%sReadParams", confName)
	fnReturns := typeOnlyParamsCE("[]interface{}", "error")
	fn := &golang.FunctionDef{
		Name:         fnName,
		Parameters:   fnParams,
		Body:         body,
		Returns:      fnReturns,
		Dependencies: nil,
	}

	return fn
}

func makeNewMapCE(name string, mapType string) *golang.FunctionCall {
	return &golang.FunctionCall{
		NewOutput: name,
		Function:  "make",
		Args:      []string{mapType},
	}
}

func makeNewArrayCE(name string, arrayType string) *golang.CodeElement {
	return &golang.CodeElement{
		Variable: createVarCE(name, arrayType),
	}
}

func prepareStmtCE(dbVar string, query string, mapName string, indexName string, returnParams []*golang.Parameter) *golang.FunctionCall {
	return &golang.FunctionCall{
		Receiver: dbVar,
		Function: "Prepare",
		Args:     []interface{}{&golang.Literal{Value: query}},
		Output:   []interface{}{&golang.Literal{Value: mapName, Indexes: indexName}, "err"},
		ErrorHandler: &golang.ErrorHandler{
			ErrorFunctionReturns: returnParams,
		},
	}
}

type NamedQuery struct {
	Name  string
	Query string
}

// Generates function to prepare statements
// Signature: func(db *sql.DB, queries map[string]string) (map[string]*sql.Stmt, error) {...}
// The function will prepare all queries passed with NamedQuery
func PrepareStmtFunction(modelName string, queries []NamedQuery) *golang.FunctionDef {
	returnFn := typeOnlyParamsCE("map[string]*sql.Stmt", "error")
	body := golang.CodeElements{
		{
			FunctionCall: makeNewMapCE("preparedCache", "map[string]*sql.Stmt"),
		},
		{
			Variable: createVarCE("err", "error"),
		},
	}

	for _, namedQuery := range queries {
		body = append(body, &golang.CodeElement{
			FunctionCall: prepareStmtCE("db", namedQuery.Query, "preparedCache", namedQuery.Name, returnFn),
		})
	}

	body = append(body, returnResultNilCE("preparedCache"))

	fnParams := []*golang.Parameter{
		{
			Name: "db",
			Type: &golang.GoType{
				Name: "*sql.DB",
			},
		},
		{
			Name: "queries",
			Type: &golang.GoType{
				Name: "map[string]string",
			},
		},
	}

	fnName := fmt.Sprintf("%sPrepareStmt", modelName)
	fn := &golang.FunctionDef{
		Name:         fnName,
		Parameters:   fnParams,
		Body:         body,
		Returns:      returnFn,
		Dependencies: nil,
	}
	return fn

}

// Generates function to setup database connection
// func SetupDBConnection() (*sql.DB, error) {...}
// It'll
// 1. read config,
// 2. do sql.Open() to host,
// 3. db.Ping() for connection,
// 4. configure pool and connection,
// 5. and return db or error
func SetupDBConnectionFunction(dataConf defs.DataConfig) (*golang.FunctionDef, error) {

	if dataConf.DatabaseConfig == nil {
		return nil, fmt.Errorf("dataconf is missing connection config")
	}
	dbConf := dataConf.DatabaseConfig
	if dbConf.DriverName == "" {
		return nil, fmt.Errorf("dataconf is missing driver")
	}

	if dbConf.DBName == "" || dbConf.Host == "" || dbConf.Port == 0 || dbConf.UserName == "" || dbConf.Password == "" {
		return nil, fmt.Errorf("dataconf is missing connection details dbname: [%s], host: [%s], port: [%d], username: [%s], password: [%s]",
			dbConf.DBName, dbConf.Host, dbConf.Port, dbConf.UserName, dbConf.Password)
	}

	if dbConf.ConnectionConfig == nil {
		return nil, fmt.Errorf("dataconf is missing connection config")
	}

	if dbConf.ConnectionPoolConfig == nil {
		return nil, fmt.Errorf("dataconf is missing connection pool config")
	}

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%d host=%s",
		dbConf.UserName, dbConf.Password, dbConf.DBName, dbConf.Port, dbConf.Host)
	fn := golang.FunctionDef{}
	fn.FunctionCode()
	returnParams := typeOnlyParamsCE("*sql.DB", "error")
	return &golang.FunctionDef{
		Name:    "SetupDBConnection",
		Imports: []string{"database/sql", "github.com/lib/pq", "time"},
		Returns: returnParams,
		Body: golang.CodeElements{
			{
				NewAssign: &golang.NewAssignment{
					Left:  []string{"driverName", "dsn", "idleConns", "connMaxLifetime"},
					Right: golang.NewLits(dbConf.DriverName, dsn, dbConf.ConnectionPoolConfig.MaxIdleConns, dbConf.ConnectionConfig.MaxLifetimeMins),
				},
			},
			goutils.FCEHNewOutReceiverArgsCE([]string{"db", "err"}, "sql", "Open",
				[]string{"driverName", "dsn"}, goutils.EHNilError("err")),

			goutils.FCEHOutReceiverArgsCE([]string{"err"}, "db", "Ping", nil, goutils.EHNilError("err")),
			goutils.FCReceiverArgsCE("db", "SetMaxIdleConns", "idleConns"),
			goutils.FCReceiverArgsCE("db", "SetConnMaxLifetime", &golang.Mul{BinaryOp: golang.BinaryOp{
				Left: &golang.Literal{Value: "time", Attribute: "Minute"}, Right: "connMaxLifetime"}}),
			returnValuesCE("db", "nil"),
		},
	}, nil
}

// Generate function, which can be used to create a singleton database instance
// Function signature is func NewRetailDB() error {...}
// return values &RetailDB{User: user, Product: product, Order: order}, nil
// For example, user is User_DB struct {db, preparedCache}
// this function will build one db instance, and prepared cache for each model.
// Since it's a singleton better to keep return type as lower case(private).
func GenerateInitFamilyFunction(modelNameMaps modelNameMappings, varName, familyTypeName string) (*golang.FunctionDef, error) {
	body := golang.CodeElements{
		goutils.FCEHOutCE([]string{"db", "err"}, "SetupDBConnection", goutils.EHError("err")), // SetupDBConnection()
	}

	for _, modelNameMap := range modelNameMaps {
		body = append(body, goutils.FCEHOutCE([]string{fmt.Sprintf("stmtMap%s", modelNameMap.ModelStructName), "err"},
			fmt.Sprintf("%sPrepareStmt", modelNameMap.ModelStructName), goutils.EHError("err"))) // PrepareCache())

		modelDBStruct := &golang.StructCreation{
			NewOutput:  golang.ToCamelCase(modelNameMap.ModelStructName),
			StructType: modelNameMap.ModelDBStructName,
			KeyValues: golang.KeyValues{
				{Key: "db", Variable: "db"},
				{Key: "preparedCache", Variable: fmt.Sprintf("stmtMap%s", modelNameMap.ModelStructName)},
			},
		}
		body = append(body, &golang.CodeElement{StructCreation: modelDBStruct})
	}

	familyKeyValues := golang.KeyValues{}
	for _, modelNameMap := range modelNameMaps {
		familyKeyValues = append(familyKeyValues,
			&golang.KeyValue{Key: modelNameMap.ModelStructName, Variable: golang.ToCamelCase(modelNameMap.ModelStructName)})
	}

	stCreate := &golang.StructCreation{
		Output:     varName,
		StructType: familyTypeName,
		KeyValues:  familyKeyValues,
	}
	body = append(body, &golang.CodeElement{StructCreation: stCreate})
	body = append(body, returnValuesCE("nil"))

	fn := &golang.FunctionDef{
		Name:    "Init" + varName,
		Returns: typeOnlyParamsCE("error"),
		Body:    body,
	}

	return fn, nil
}
