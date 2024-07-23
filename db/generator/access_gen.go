package generator

import (
	"fmt"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"stellarsky.ai/platform/codegen/data-service-generator/base"
	"stellarsky.ai/platform/codegen/data-service-generator/config"
	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang"
	datahelpers "stellarsky.ai/platform/codegen/data-service-generator/db/generator/data-helpers"
	"stellarsky.ai/platform/codegen/data-service-generator/db/generator/defs"
	"stellarsky.ai/platform/codegen/data-service-generator/db/models"
)

func GenerateDB(model *models.Model, dataConfig *defs.DataConfig) ([]*golang.UnitModule, error) {

	unitModules := make([]*golang.UnitModule, 0)
	for _, config := range dataConfig.Models {
		srcFile, err := Generate(config)
		if err != nil {
			return nil, err
		}

		unitModule := &golang.UnitModule{
			Name:         golang.ToSnakeCase(config.Model.Name),
			Structs:      srcFile.Structs,
			Functions:    srcFile.Functions,
			Variables:    srcFile.Variables,
			Constants:    srcFile.Constants,
			Imports:      srcFile.Imports,
			Dependencies: srcFile.Dependencies,
		}
		unitModules = append(unitModules, unitModule)

	}

	return unitModules, nil

}

func readTypeAndValidations(attributeId int64) (string, *golang.GoType, []*models.Validation, error) {
	attribute, ok := config.Attributes[attributeId]
	if !ok {
		return "", nil, nil, fmt.Errorf("attribute %d not found", attributeId)
	}

	typeId := attribute.TypeId
	validationIds := attribute.ValidationIds

	postgresType := datahelpers.GetPostgresType(typeId)
	goTypeStr := datahelpers.PostgresToGoType(postgresType)
	goType, err := golang.TranslateToGoType(goTypeStr)
	if err != nil {
		return "", nil, nil, err
	}
	base.LOG.Info("ReadTypeValidations", "goType", goType, "validationIds", validationIds, "attribute", attribute,
		"attributeId", attributeId, "typeId", typeId, "postgresType", postgresType, "goTypeStr", goTypeStr)
	validations := datahelpers.GetValidations(validationIds)
	return attribute.Name, goType, validations, nil
}

func generateModel(config *defs.ModelConfig) ([]*golang.Struct, []*golang.Function, error) {

	models := make([]*golang.Struct, 0, 1)
	functions := make([]*golang.Function, 0, 1)

	nameWithTypes := make([]golang.NameWithType, 0, 1)
	for _, attribute := range config.Model.Attributes {
		attrName, goType, _, err := readTypeAndValidations(attribute)
		if err != nil {
			return nil, nil, err
		}
		base.LOG.Info("Attribute", "attrName", attrName, "goType", goType)

		nameWithTypes = append(nameWithTypes, golang.NameWithType{
			Name: attrName,
			Type: goType,
		})
	}

	modelStruct := golang.GenerateStructForJSON(config.Model.Name, nameWithTypes)
	models = append(models, modelStruct)
	return models, functions, nil
}

func Generate(config defs.ModelConfig) (*golang.GoSourceFile, error) {

	allQueries := make([]NamedQuery, 0)
	allFunctions := make([]*golang.Function, 0)
	allStructs := make([]*golang.Struct, 0)

	models, fns, err := generateModel(&config)
	if err != nil {
		return nil, err
	}
	allStructs = append(allStructs, models...)
	allFunctions = append(allFunctions, fns...)

	caser := cases.Title(language.English)
	modelName := caser.String(config.Model.Name)

	queries, accessFns, structs, err := GenerateFindConfigs(modelName, config.Access.Find)
	if err != nil {
		return nil, err
	}
	allQueries = append(allQueries, queries...)
	allFunctions = append(allFunctions, accessFns...)
	allStructs = append(allStructs, structs...)

	queries, accessFns, structs, err = GenerateUpdateConfigs(modelName, config.Access.Update)
	if err != nil {
		return nil, err
	}
	allQueries = append(allQueries, queries...)
	allFunctions = append(allFunctions, accessFns...)
	allStructs = append(allStructs, structs...)

	queries, accessFns, structs, err = GenerateAddConfigs(modelName, config.Access.Add)
	if err != nil {
		return nil, err
	}
	allQueries = append(allQueries, queries...)
	allFunctions = append(allFunctions, accessFns...)
	allStructs = append(allStructs, structs...)

	queries, accessFns, structs, err = GenerateAddOrReplaceConfigs(modelName, config.Access.AddOrReplace)
	if err != nil {
		return nil, err
	}
	allQueries = append(allQueries, queries...)
	allFunctions = append(allFunctions, accessFns...)
	allStructs = append(allStructs, structs...)

	queries, accessFns, structs, err = GenerateDeleteConfigs(modelName, config.Access.Delete)
	if err != nil {
		return nil, err
	}
	allQueries = append(allQueries, queries...)
	allFunctions = append(allFunctions, accessFns...)
	allStructs = append(allStructs, structs...)

	prepareFn := PrepareStmtFunction(modelName, allQueries)
	allFunctions = append(allFunctions, prepareFn)
	goSrc := &golang.GoSourceFile{
		Package:      "database",
		Structs:      allStructs,
		Functions:    allFunctions,
		InitFunction: nil,
		Variables:    nil,
		Constants:    nil}

	return goSrc, nil
}

func generateParamsStruct(paramRefs []defs.ParameterRef, name string) *golang.Struct {
	nameWithTypes := make([]golang.NameWithType, 0, len(paramRefs))
	for _, param := range paramRefs {
		nameWithTypes = append(nameWithTypes, golang.NameWithType{
			Name: param.Name,
			Type: golang.GoInterfaceType,
		})
	}

	return golang.GenerateStructForJSON(fmt.Sprintf("%sParams", name), nameWithTypes)
}

func generateRequestStruct(name string, paramStructName string) *golang.Struct {
	namedWithTypes := []golang.NameWithType{
		{Name: "Params", Type: &golang.GoType{Name: paramStructName}},
	}
	return golang.GenerateStructForJSON(fmt.Sprintf("%sRequest", name), namedWithTypes)
}

func generateAccessStructs(paramRef []defs.ParameterRef, name string) []*golang.Struct {
	paramStruct := generateParamsStruct(paramRef, name)
	reqStruct := generateRequestStruct(name, paramStruct.Name)
	return []*golang.Struct{paramStruct, reqStruct}
}

func GenerateFindConfigs(modelName string, findConfig []defs.AccessConfig) ([]NamedQuery, []*golang.Function, []*golang.Struct, error) {

	functions := make([]*golang.Function, 0, len(findConfig))
	reqs := make([]*golang.Struct, 0, len(findConfig))
	queries := make([]NamedQuery, 0, len(findConfig))

	for _, conf := range findConfig {

		query, paramRefs := datahelpers.MakeFindQuery(modelName, &conf)
		queries = append(queries, NamedQuery{Name: conf.Name, Query: query})
		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)

		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)

		fn := FindCodeFunction(modelName, conf.Name, conf.Attributes)
		functions = append(functions, fn)
	}

	return queries, functions, reqs, nil

}

func GenerateUpdateConfigs(modelName string, updateConfig []defs.AccessConfig) ([]NamedQuery, []*golang.Function, []*golang.Struct, error) {

	functions := make([]*golang.Function, 0, len(updateConfig))
	reqs := make([]*golang.Struct, 0, len(updateConfig))
	queries := make([]NamedQuery, 0, len(updateConfig))

	for _, conf := range updateConfig {

		query, paramRefs := datahelpers.MakeUpdateQuery(modelName, &conf)
		queries = append(queries, NamedQuery{Name: conf.Name, Query: query})

		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)
		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)

		fn := UpdateCodeFunction(conf.Name)
		functions = append(functions, fn)
	}

	return queries, functions, reqs, nil

}

func GenerateAddConfigs(modelName string, addConfig []defs.AccessConfig) ([]NamedQuery, []*golang.Function, []*golang.Struct, error) {
	functions := make([]*golang.Function, 0, len(addConfig))
	reqs := make([]*golang.Struct, 0, len(addConfig))
	queries := make([]NamedQuery, 0, len(addConfig))

	for _, conf := range addConfig {
		query, paramRefs := datahelpers.MakeAddQuery(modelName, &conf)
		queries = append(queries, NamedQuery{Name: conf.Name, Query: query})
		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)
		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)
		fn := AddCodeFunction(conf.Name)
		functions = append(functions, fn)
	}
	return queries, functions, reqs, nil

}

func GenerateAddOrReplaceConfigs(modelName string, addOrReplaceConfig []defs.AccessConfig) ([]NamedQuery, []*golang.Function, []*golang.Struct, error) {
	functions := make([]*golang.Function, 0, len(addOrReplaceConfig))
	reqs := make([]*golang.Struct, 0, len(addOrReplaceConfig))
	queries := make([]NamedQuery, 0, len(addOrReplaceConfig))

	for _, conf := range addOrReplaceConfig {
		query, paramRefs := datahelpers.MakeAddOrReplaceQuery(modelName, &conf)
		queries = append(queries, NamedQuery{Name: conf.Name, Query: query})
		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)
		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)
		fn := AddOrReplaceCodeFunction(conf.Name)
		functions = append(functions, fn)
	}
	return queries, functions, reqs, nil
}

func GenerateDeleteConfigs(modelName string, deleteConfig []defs.AccessConfig) ([]NamedQuery, []*golang.Function, []*golang.Struct, error) {
	functions := make([]*golang.Function, 0, len(deleteConfig))
	reqs := make([]*golang.Struct, 0, len(deleteConfig))
	queries := make([]NamedQuery, 0, len(deleteConfig))

	for _, conf := range deleteConfig {
		query, paramRefs := datahelpers.MakeDeleteQuery(modelName, &conf)
		queries = append(queries, NamedQuery{Name: conf.Name, Query: query})
		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)
		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)
		fn := DeleteCodeFunction(conf.Name)
		functions = append(functions, fn)
	}
	return queries, functions, reqs, nil
}

func SetupDatabaseFunction(dataConf []defs.DataConfig) *golang.Function {

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%d host=%s", "postgres", "postgres", "postgres", 5432, "localhost")
	fmt.Println(dsn)
	return &golang.Function{
		Name:    "SetupDatabase",
		Imports: []string{"database/sql", "github.com/lib/pq"},
		Returns: typeOnlyParamsCE("error"),
		Body: golang.CodeElements{
			{
				StructCreation: &golang.StructCreation{
					NewOutput:  "cfg",
					ModuleName: "pg",
					StructType: "Config",
					KeyValues: golang.KeyValues{
						{Key: "User", Value: "postgres"},
						{Key: "Password", Value: "<PASSWORD>"},
						{Key: "Database", Value: "postgres"},
						{Key: "Port", Value: 5432},
						{Key: "Host", Value: "localhost"},
					},
				},
			},
			{
				FunctionCall: &golang.FunctionCall{
					NewOutput: []string{"db", "err"},
					Function:  "Open",
					Receiver:  "db",
					Args: []interface{}{
						"postgres",
						&golang.FunctionCall{
							Function: "FormatDSN",
							Receiver: "cfg",
						},
					},
					ErrorHandler: golang.ErrorHandler{
						ErrorFunctionReturns: typeOnlyParamsCE("error"),
					},
				},
			},
			{
				FunctionCall: &golang.FunctionCall{
					Output:   "err",
					Function: "Ping",
					Receiver: "db",
					ErrorHandler: golang.ErrorHandler{
						ErrorFunctionReturns: typeOnlyParamsCE("error"),
					},
				},
			},
			{
				FunctionCall: &golang.FunctionCall{
					Function: "SetMaxIdleConns",
					Receiver: "db",
					Args:     10,
				},
			},
			{
				FunctionCall: &golang.FunctionCall{
					Function: "SetConnMaxLifetime",
					Receiver: "db",
					Args: &golang.Mul{BinaryOp: golang.BinaryOp{
						Left: &golang.Literal{
							Value:     "time",
							Attribute: "Minute",
						},
						Right: 30}},
				},
			},
			returnValuesCE("nil"),
		},
	}
}
