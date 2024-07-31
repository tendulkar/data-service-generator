package generator

import (
	"fmt"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"stellarsky.ai/platform/codegen/data-service-generator/base"
	"stellarsky.ai/platform/codegen/data-service-generator/config"
	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang"
	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang/goutils"
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

type AccessFnGenerator func(modelName string, config []defs.AccessConfig) ([]NamedQuery, []*golang.Function, []*golang.Struct, error)

func genAccessFn(modelName string, config []defs.AccessConfig, accessFn AccessFnGenerator,
	allQueries *[]NamedQuery, allFunctions *[]*golang.Function, allStructs *[]*golang.Struct) error {

	queries, accessFns, structs, err := accessFn(modelName, config)
	if err != nil {
		return err
	}
	*allQueries = append(*allQueries, queries...)
	*allFunctions = append(*allFunctions, accessFns...)
	*allStructs = append(*allStructs, structs...)
	return nil
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
	accessMethods := []AccessFnGenerator{
		GenerateFindConfigs,
		GenerateUpdateConfigs,
		GenerateAddConfigs,
		GenerateAddOrReplaceConfigs,
		GenerateDeleteConfigs,
	}

	accessConfigs := [][]defs.AccessConfig{
		config.Access.Find,
		config.Access.Update,
		config.Access.Add,
		config.Access.AddOrReplace,
		config.Access.Delete,
	}

	for i, accessMethod := range accessMethods {
		err = genAccessFn(modelName, accessConfigs[i], accessMethod, &allQueries, &allFunctions, &allStructs)
		if err != nil {
			return nil, err
		}
	}

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

func SetupDatabaseFunction(dataConf defs.DataConfig) (*golang.Function, error) {

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
	fn := golang.Function{}
	fn.FunctionCode()
	returnParams := typeOnlyParamsCE("error")
	return &golang.Function{
		Name:    "SetupDatabase",
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
				[]string{"driverName", "dsn"}, goutils.EHError("err")),

			goutils.FCEHOutReceiverArgsCE([]string{"err"}, "db", "Ping", nil, goutils.EHError("err")),
			goutils.FCReceiverArgsCE("db", "SetMaxIdleConns", "idleConns"),
			goutils.FCReceiverArgsCE("db", "SetConnMaxLifetime", &golang.Mul{BinaryOp: golang.BinaryOp{
				Left: &golang.Literal{Value: "time", Attribute: "Minute"}, Right: "connMaxLifetime"}}),
			returnValuesCE("nil"),
		},
	}, nil
}
