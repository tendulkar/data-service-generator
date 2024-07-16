package generator

import (
	"fmt"

	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang"
)

func lookupStmtCE(confName string, dbName, cacheName, stmtName string) *golang.MapLookup {
	return &golang.MapLookup{
		NewOutput: stmtName,
		Receiver:  dbName,
		Name:      cacheName,
		Key:       confName,
	}
}

func parseParamsCE(confName string, requestName string, valuesName string) *golang.FunctionCall {
	return &golang.FunctionCall{
		NewOutput: []string{valuesName, "err"},
		Function:  fmt.Sprintf("%sParseParams", confName),
		Args:      []string{requestName},
		ErrorHandler: golang.ErrorHandler{
			ErrorReturns: []string{"nil", "err"},
		},
	}
}

func queryStmtCE(stmtName, valuesName, rowsName string) *golang.FunctionCall {
	return &golang.FunctionCall{
		NewOutput: []string{rowsName, "err"},
		Receiver:  stmtName,
		Function:  "Query",
		Args:      []string{fmt.Sprintf("%s...", valuesName)},
		ErrorHandler: golang.ErrorHandler{
			ErrorReturns: []string{"nil", "err"},
		},
		CleanningHandler: golang.CleanningHandler{
			Receiver: rowsName,
			Function: "Close",
		},
	}
}

func queryRowStmtCE(stmtName string, args []string, resultName, zeroResultName string) *golang.FunctionCall {
	return &golang.FunctionCall{
		NewOutput: []string{resultName, "err"},
		Receiver:  fmt.Sprintf("%s.QueryRow()", stmtName),
		Function:  "Exec",
		Args:      args,
		ErrorHandler: golang.ErrorHandler{
			ErrorReturns: []string{zeroResultName, "err"},
		},
	}
}

func execStmtCE(stmtName, valuesName, resultName string) *golang.FunctionCall {
	return &golang.FunctionCall{
		NewOutput: []string{resultName, "err"},
		Receiver:  stmtName,
		Function:  "Exec",
		Args:      []string{fmt.Sprintf("%s...", valuesName)},
		ErrorHandler: golang.ErrorHandler{
			ErrorReturns: []string{"0", "err"},
		},
	}
}

func createVarCE(name string, typ string) *golang.VariableCreate {
	return &golang.VariableCreate{
		Name: name,
		Type: typ,
	}
}

func appendCE(arr string, value string) *golang.FunctionCall {
	return &golang.FunctionCall{
		Output:   arr,
		Function: "append",
		Args:     []string{arr, value},
	}
}

func returnResultNilCE(results string) *golang.CodeElement {
	return &golang.CodeElement{
		Return: []string{results, "nil"},
	}
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
		ErrorHandler: golang.ErrorHandler{
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
		Type: golang.GoType{
			Name: "context.Context",
		},
	}
}

func dbParamCE(dbName string) *golang.Parameter {
	return &golang.Parameter{
		Name: dbName,
		Type: golang.GoType{
			Name:   "*sql.DB",
			Source: "database/sql",
		},
	}
}

func requestParamCE(confName string, requestName string) *golang.Parameter {
	return &golang.Parameter{
		Name: requestName,
		Type: golang.GoType{
			Name: fmt.Sprintf("%sRequest", confName),
		},
	}
}

func dbRequestParamsCE(dbName, name, requestName string) []*golang.Parameter {
	return []*golang.Parameter{
		dbParamCE(dbName),
		requestParamCE(name, requestName),
	}
}

func ctxDBRequestParamsCE(ctxName, dbName, name, requestName string) []*golang.Parameter {
	params := []*golang.Parameter{
		ctxParamCE(ctxName),
	}
	params = append(params, dbRequestParamsCE(dbName, name, requestName)...)
	return params
}

func errorParamCE(errorName string) *golang.Parameter {
	return &golang.Parameter{
		Name: errorName,
		Type: golang.GoType{
			Name: "error",
		},
	}
}

func resultsParamCE(modelName, resultsName, sourceName string) *golang.Parameter {
	return &golang.Parameter{
		Name: resultsName,
		Type: golang.GoType{
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
			Type: golang.GoType{
				Name: typ,
			},
		})
	}
	return params
}

func callResultErrorCE(objName, funcName string, args []string, resultName string, errorName string) *golang.FunctionCall {
	return &golang.FunctionCall{
		NewOutput: []string{resultName, errorName},
		Receiver:  objName,
		Function:  funcName,
		Args:      args,
	}
}

func findCodeBody(modelName, name string, attributes []string) *golang.Function {
	codeElems := golang.CodeElements{
		{
			MapLookup: lookupStmtCE(name, "db", "preparedCache", "stmt"),
		},
		{
			FunctionCall: parseParamsCE(name, "request", "values"),
		},
		{
			FunctionCall: queryStmtCE("stmt", "values", "rows"),
		},
		{
			Variable: createVarCE("results", fmt.Sprintf("[]%s", modelName)),
		},
		{
			RepeatCond: scanResultsCE(name, attributes, "results"),
		},
		{
			Return: []string{"results", "nil"},
		},
	}

	params := ctxDBRequestParamsCE("ctx", "db", name, "request")
	returns := resultsErrorParamsCE(modelName, "results", "err", "")
	dependencies := []golang.Dependency{}

	fn := &golang.Function{
		Name:         name,
		Parameters:   params,
		Body:         codeElems,
		Returns:      returns,
		Dependencies: dependencies,
	}

	return fn
}

func updateCodeBody(name string) *golang.Function {
	codeElems := golang.CodeElements{
		{
			MapLookup: lookupStmtCE(name, "db", "preparedCache", "stmt"),
		},
		{
			FunctionCall: parseParamsCE(name, "request", "values"),
		},
		{
			FunctionCall: execStmtCE("stmt", "values", "result"),
		},
		{
			FunctionCall: callResultErrorCE("result", "RowsAffected", []string{}, "rowsEffected", "err"),
		},
		{
			Return: []string{"rowsEffected", "err"},
		},
	}
	params := ctxDBRequestParamsCE("ctx", "db", name, "request")
	returns := typeOnlyParamsCE("int64", "error")
	dependencies := make([]golang.Dependency, 0, 0)
	fn := &golang.Function{
		Name:         name,
		Parameters:   params,
		Body:         codeElems,
		Returns:      returns,
		Dependencies: dependencies,
	}
	return fn
}

func addCodeFunction(name string) *golang.Function {
	codeElems := golang.CodeElements{
		{
			MapLookup: lookupStmtCE(name, "db", "preparedCache", "stmt"),
		},
		{
			FunctionCall: parseParamsCE(name, "request", "values"),
		},
		{
			Variable: createVarCE("id", "int64"),
		},
		{
			FunctionCall: queryRowStmtCE("stmt", []string{"&id"}, "result"),
		},
		{
			Return: []string{"id", "err"},
		},
	}
	params := ctxDBRequestParamsCE("ctx", "db", name, "request")
	returns := typeOnlyParamsCE("int64", "error")
	dependencies := make([]golang.Dependency, 0, 0)
	fn := &golang.Function{
		Name:         name,
		Parameters:   params,
		Body:         codeElems,
		Returns:      returns,
		Dependencies: dependencies,
	}
	return fn
}

func addOrReplaceCodeFunction(name string) *golang.Function {
	codeElems := golang.CodeElements{
		{
			MapLookup: lookupStmtCE(name, "db", "preparedCache", "stmt"),
		},
		{
			FunctionCall: parseParamsCE(name, "request", "values"),
		},
		{
			Variable: createVarCE("id", "int64"),
		},
		{
			Variable: createVarCE("inserted", "bool"),
		},
		{
			FunctionCall: queryRowStmtCE("stmt", []string{"&id", "&inserted"}, "result", "err"),
		},
		{
			Return: []string{"id", "inserted", "err"},
		},
	}
	params := ctxDBRequestParamsCE("ctx", "db", name, "request")
	returns := typeOnlyParamsCE("int64", "bool", "error")
	dependencies := make([]golang.Dependency, 0, 0)
	fn := &golang.Function{
		Name:         name,
		Parameters:   params,
		Body:         codeElems,
		Returns:      returns,
		Dependencies: dependencies,
	}
	return fn
}

func deleteCodeFunction(name string) *golang.Function {
	codeElems := golang.CodeElements{
		{
			MapLookup: lookupStmtCE(name, "db", "preparedCache", "stmt"),
		},
		{
			FunctionCall: parseParamsCE(name, "request", "values"),
		},
		{
			FunctionCall: execStmtCE("stmt", "values", "result"),
		},
		{
			FunctionCall: callResultErrorCE("result", "RowsAffected", []string{}, "rowsEffected", "err"),
		},
		{
			Return: []string{"rowsEffected", "err"},
		},
	}
	params := ctxDBRequestParamsCE("ctx", "db", name, "request")
	returns := typeOnlyParamsCE("int64", "error")
	dependencies := make([]golang.Dependency, 0, 0)
	fn := &golang.Function{
		Name:         name,
		Parameters:   params,
		Body:         codeElems,
		Returns:      returns,
		Dependencies: dependencies,
	}
	return fn
}
