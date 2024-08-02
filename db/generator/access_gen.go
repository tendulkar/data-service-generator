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

type modelNameMapping struct {
	ModelName         string
	ModelStructName   string
	ModelDBStructName string
}

type modelNameMappings []*modelNameMapping

func GenerateDB(dataConfig *defs.DataConfig) ([]*golang.UnitModule, error) {

	if dataConfig.FamilyName == "" {
		return nil, fmt.Errorf("dataconf is missing family name")
	}

	if dataConfig.Models == nil {
		return nil, fmt.Errorf("dataconf is missing models config")
	}

	if dataConfig.DatabaseConfig == nil {
		return nil, fmt.Errorf("dataconf is missing database config")
	}

	unitModules := make([]*golang.UnitModule, 0)
	modelNameMaps := make(modelNameMappings, 0)
	for _, config := range dataConfig.Models {
		srcFile, modelNameMap, err := Generate(config)
		if err != nil {
			base.LOG.Error("GenerateDB::Error generating code for model %s: %v", config.Model.Name, err)
			return nil, err
		}
		modelNameMaps = append(modelNameMaps, modelNameMap)

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

	st, fn, dbvar, err := GenerateFamily(dataConfig.FamilyName, modelNameMaps)
	if err != nil {
		base.LOG.Error("GenerateDB::Error generating code for family %s: %v", dataConfig.FamilyName, err)
		return nil, err
	}
	unitModules = append(unitModules, &golang.UnitModule{
		Name:         dataConfig.FamilyName,
		Structs:      st,
		Functions:    fn,
		Variables:    dbvar,
		Constants:    nil,
		Imports:      nil,
		Dependencies: nil,
	})

	return unitModules, nil

}

func GenerateFamily(familyName string, modelNameMaps modelNameMappings) ([]*golang.StructDef, []*golang.FunctionDef, []*golang.Variable, error) {

	structs := make([]*golang.StructDef, 0)
	functions := make([]*golang.FunctionDef, 0)

	structName := golang.ToCamelCase(familyName)
	varName := golang.ToPascalCase(familyName)
	nameWithTypes := make([]golang.NameWithType, 0, 1)

	for _, nameMap := range modelNameMaps {
		// TODO: read modelName from actual model struct, rather than deriving it, it's dangerous and can break at any time.
		modelName := nameMap.ModelStructName
		nameWithTypes = append(nameWithTypes, golang.NameWithType{
			Name: modelName,
			Type: &golang.GoType{Name: modelName},
		})
	}

	st := golang.GenStructForDataModel(structName, nameWithTypes, false, false, false)
	fn, err := GenerateInitFamilyFunction(modelNameMaps, varName, structName)
	if err != nil {
		return nil, nil, nil, err
	}
	structs = append(structs, st)
	functions = append(functions, fn)

	varDeclare := []*golang.Variable{{
		Names:       varName,
		Type:        structName,
		IsReference: true,
	},
	}
	return structs, functions, varDeclare, nil
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
	base.LOG.Debug("ReadTypeValidations", "goType", goType, "validationIds", validationIds, "attribute", attribute,
		"attributeId", attributeId, "typeId", typeId, "postgresType", postgresType, "goTypeStr", goTypeStr)
	validations := datahelpers.GetValidations(validationIds)
	return attribute.Name, goType, validations, nil
}

// Generate Model struct for a given model
//
//	Example: type Product struct {
//		Sku         string  `db:"sku"`
//		ProductName string  `db:"product_name"`
//		Description string  `db:"description"`
//		Price       float64 `db:"price"`
//	}

func generateModel(config *defs.ModelConfig) (*modelNameMapping, []*golang.StructDef, []*golang.FunctionDef, error) {

	models := make([]*golang.StructDef, 0, 1)
	functions := make([]*golang.FunctionDef, 0, 1)
	modelNameMap := &modelNameMapping{
		ModelName:         config.Model.Name,
		ModelStructName:   golang.ToPascalCase(config.Model.Name),
		ModelDBStructName: golang.ToPascalCase(config.Model.Name) + "_DB",
	}

	nameWithTypes := make([]golang.NameWithType, 0, 1)
	for _, attribute := range config.Model.Attributes {
		attrName, goType, _, err := readTypeAndValidations(attribute)
		if err != nil {
			return nil, nil, nil, err
		}
		base.LOG.Debug("Attribute", "attrName", attrName, "goType", goType)

		nameWithTypes = append(nameWithTypes, golang.NameWithType{
			Name: attrName,
			Type: goType,
		})
	}

	modelStruct := golang.GenStructForDataModel(modelNameMap.ModelStructName, nameWithTypes, false, false, true)

	dbNameWithTypes := []golang.NameWithType{
		{Name: "db", Type: &golang.GoType{Name: "*sql.DB"}},
		{Name: "preparedCache", Type: &golang.GoType{Name: "map[string]*sql.Stmt"}},
	}

	modelDBStruct, modelDBNewFn := golang.GenStructWithNewFunction(modelNameMap.ModelDBStructName, dbNameWithTypes, true, false, false, false)
	models = append(models, modelStruct, modelDBStruct)
	functions = append(functions, modelDBNewFn)

	return modelNameMap, models, functions, nil
}

type AccessFnGenerator func(modelName string, modelDBName string, config []defs.AccessConfig) ([]NamedQuery, []*golang.FunctionDef, []*golang.StructDef, error)

func genAccessFn(modelName string, modelDBName string, config []defs.AccessConfig, accessFn AccessFnGenerator,
	allQueries *[]NamedQuery, allFunctions *[]*golang.FunctionDef, allStructs *[]*golang.StructDef) error {

	queries, accessFns, structs, err := accessFn(modelName, modelDBName, config)
	if err != nil {
		return err
	}
	*allQueries = append(*allQueries, queries...)
	*allFunctions = append(*allFunctions, accessFns...)
	*allStructs = append(*allStructs, structs...)
	return nil
}

func Generate(config defs.ModelConfig) (*golang.GoSourceFile, *modelNameMapping, error) {

	allQueries := make([]NamedQuery, 0)
	allFunctions := make([]*golang.FunctionDef, 0)
	allStructs := make([]*golang.StructDef, 0)
	caser := cases.Title(language.English)
	modelName := caser.String(config.Model.Name)

	// Generate Model struct for a given model, for example `type User struct {<fields with db tags>}`
	modelNameMap, models, fns, err := generateModel(&config)
	if err != nil {
		return nil, nil, err
	}
	allStructs = append(allStructs, models...)
	allFunctions = append(allFunctions, fns...)

	// Generate methods for SELECT, UPDATE, INSERT, INSERT OR UPDATE, DELETE for a given model
	err = geneateAllAccessMethods(config, modelNameMap.ModelStructName, modelNameMap.ModelDBStructName,
		&allQueries, &allFunctions, &allStructs)
	if err != nil {
		base.LOG.Error("Generate::geneateAllAccessMethods", "err", err, "model", modelName, "modelMap", *modelNameMap)
		return nil, nil, err
	}

	// PrepareStmt function will prepare all queries for a given model
	// Make sure allQueries have been populated,
	// so don't move the PrepareStmtFunction call to the above geneateAllAccessMethods
	prepareFn := PrepareStmtFunction(modelName, allQueries)
	allFunctions = append(allFunctions, prepareFn)

	goSrc := &golang.GoSourceFile{
		Package:      "database",
		Structs:      allStructs,
		Functions:    allFunctions,
		InitFunction: nil,
		Variables:    nil,
		Constants:    nil}

	return goSrc, modelNameMap, nil
}

// All access methods for a given model (Find, Update, Add, AddOrReplace and Delete),
// will do query on database with above prepared statements (SELECT, UPDATE, INSERT, INSERT OR UPDATE, DELETE)
func geneateAllAccessMethods(config defs.ModelConfig, modelName string, modelDBName string,
	allQueries *[]NamedQuery, allFunctions *[]*golang.FunctionDef, allStructs *[]*golang.StructDef) error {
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
		err := genAccessFn(modelName, modelDBName, accessConfigs[i], accessMethod, allQueries, allFunctions, allStructs)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateParamsStruct(paramRefs []defs.ParameterRef, name string) *golang.StructDef {
	nameWithTypes := make([]golang.NameWithType, 0, len(paramRefs))
	for _, param := range paramRefs {
		nameWithTypes = append(nameWithTypes, golang.NameWithType{
			Name: param.Name,
			Type: golang.GoInterfaceType,
		})
	}

	return golang.GenStructForDataModel(fmt.Sprintf("%sParams", name), nameWithTypes, true, false, false)
}

func generateRequestStruct(name string, paramStructName string) *golang.StructDef {
	namedWithTypes := []golang.NameWithType{
		{Name: "Params", Type: &golang.GoType{Name: paramStructName}},
	}
	return golang.GenStructForDataModel(fmt.Sprintf("%sRequest", name), namedWithTypes, true, false, false)
}

func generateAccessStructs(paramRef []defs.ParameterRef, name string) []*golang.StructDef {
	paramStruct := generateParamsStruct(paramRef, name)
	reqStruct := generateRequestStruct(name, paramStruct.Name)
	return []*golang.StructDef{paramStruct, reqStruct}
}

// GenerateFindConfigs will generate all SELECT queries for a given model
// It generates find function, and one helper function for reading params from request to bind values to query
// It generates 2 structs for params and request, request is input (arg) to find function, and params is part of request
// Note that we already have a prepared statement cache for each query, Product{<fields>} struct, Product_DB struct {db, preparedCache fields}
// Example:
//
//		type GetProductByIdparams struct {
//			Id interface{} `json:"id"`
//		}
//
//		type GetProductByIdrequest struct {
//			Params GetProductByIdparams `json:"params"`
//		}
//
//		func GetProductByIDReadParams(params GetProductByIDParams) ([]interface{}, error) {
//	        var values []interface{}
//	        values = append(values, params.id)
//	        return values, nil
//		}
//
//	    func GetProductByID(ctx context.Context, db *Product_DB, requestParams GetProductByIDParams) (results []Product, err error) {
//			stmt := db.preparedCache["GetProductByID"]
//			values, err := GetProductByIDParseParams(requestParams)
//			if err != nil {
//					return nil, err
//			}
//			rows, err := stmt.Query(values...)
//			if err != nil {
//					return nil, err
//			}
//			defer rows.Close()
//			var results []Product
//			for rows.Next() {
//					var item Product
//					scanErr := rows.Scan(&item.id, &item.name, &item.price, &item.quantity)
//					if scanErr != nil {
//							return nil, scanErr
//					}
//					results = append(results, item)
//			}
//			return results, nil
//	}
func GenerateFindConfigs(modelName string, modelDBName string, findConfig []defs.AccessConfig) ([]NamedQuery, []*golang.FunctionDef, []*golang.StructDef, error) {

	functions := make([]*golang.FunctionDef, 0, len(findConfig))
	reqs := make([]*golang.StructDef, 0, len(findConfig))
	queries := make([]NamedQuery, 0, len(findConfig))

	for _, conf := range findConfig {

		query, paramRefs := datahelpers.MakeFindQuery(modelName, &conf)
		queries = append(queries, NamedQuery{Name: conf.Name, Query: query})
		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)

		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)

		fn := FindCodeFunction(modelName, modelDBName, conf.Name, conf.Attributes)
		functions = append(functions, fn)
	}

	return queries, functions, reqs, nil

}

// GenerateUpdateConfigs will generate all UPDATE queries for a given model
// Similar to GenerateFindConfigs
// It generates Update function, and one helper function for reading params from request to bind values to query
// It generates 2 structs for params and request, request is input (arg) to Update function, and params is part of request
func GenerateUpdateConfigs(modelName string, modelDBName string, updateConfig []defs.AccessConfig) ([]NamedQuery, []*golang.FunctionDef, []*golang.StructDef, error) {

	functions := make([]*golang.FunctionDef, 0, len(updateConfig))
	reqs := make([]*golang.StructDef, 0, len(updateConfig))
	queries := make([]NamedQuery, 0, len(updateConfig))

	for _, conf := range updateConfig {

		query, paramRefs := datahelpers.MakeUpdateQuery(modelName, &conf)
		queries = append(queries, NamedQuery{Name: conf.Name, Query: query})

		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)
		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)

		fn := UpdateCodeFunction(conf.Name, modelDBName)
		functions = append(functions, fn)
	}

	return queries, functions, reqs, nil

}

// GenerateAddConfigs will generate all INSERT queries for a given model
// Similar to GenerateFindConfigs
// It generates Add(INSERT) function, and one helper function for reading params from request to bind values to query
// It generates 2 structs for params and request, request is input (arg) to Add function, and params is part of request
func GenerateAddConfigs(modelName string, modelDBName string, addConfig []defs.AccessConfig) ([]NamedQuery, []*golang.FunctionDef, []*golang.StructDef, error) {
	functions := make([]*golang.FunctionDef, 0, len(addConfig))
	reqs := make([]*golang.StructDef, 0, len(addConfig))
	queries := make([]NamedQuery, 0, len(addConfig))

	for _, conf := range addConfig {
		query, paramRefs := datahelpers.MakeAddQuery(modelName, &conf)
		queries = append(queries, NamedQuery{Name: conf.Name, Query: query})
		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)
		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)
		fn := AddCodeFunction(conf.Name, modelDBName)
		functions = append(functions, fn)
	}
	return queries, functions, reqs, nil

}

// GenerateAddOrReplaceConfigs will generate all INSERT OR UPDATE queries for a given model
// Similar to GenerateFindConfigs
// It generates AddOrReplace(INSERT OR UPDATE) function, and one helper function for reading params from request to bind values to query
// It generates 2 structs for params and request, request is input (arg) to AddOrReplace function, and params is part of request
func GenerateAddOrReplaceConfigs(modelName string, modelDBName string, addOrReplaceConfig []defs.AccessConfig) ([]NamedQuery, []*golang.FunctionDef, []*golang.StructDef, error) {
	functions := make([]*golang.FunctionDef, 0, len(addOrReplaceConfig))
	reqs := make([]*golang.StructDef, 0, len(addOrReplaceConfig))
	queries := make([]NamedQuery, 0, len(addOrReplaceConfig))

	for _, conf := range addOrReplaceConfig {
		query, paramRefs := datahelpers.MakeAddOrReplaceQuery(modelName, &conf)
		queries = append(queries, NamedQuery{Name: conf.Name, Query: query})
		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)
		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)
		fn := AddOrReplaceCodeFunction(conf.Name, modelDBName)
		functions = append(functions, fn)
	}
	return queries, functions, reqs, nil
}

// GenerateDeleteConfigs will generate all DELETE queries for a given model
// Similar to GenerateFindConfigs
// It generates Delete function, and one helper function for reading params from request to bind values to query
// It generates 2 structs for params and request, request is input (arg) to Delete function, and params is part of request
//
//	 type DeleteProductParams struct {
//		Id interface{} `json:"id"`
//	}
//
//	type DeleteProductRequest struct {
//			Params DeleteProductParams `json:"params"`
//	}
//
//	func DeleteProductReadParams(params DeleteProductParams) ([]interface{}, error) {
//		var values []interface{}
//		values = append(values, params.id)
//		return values, nil
//	}
//
//	func DeleteProduct(ctx context.Context, db *sql.DB, requestParams DeleteProductParams) (int64, error) {
//		stmt := db.preparedCache["DeleteProduct"]
//		values, err := DeleteProductParseParams(requestParams)
//		if err != nil {
//				return int64(0), err
//		}
//		result, err := stmt.Exec(values...)
//		if err != nil {
//				return int64(0), err
//		}
//		rowsAffected, err := result.RowsAffected()
//		if err != nil {
//				return int64(0), err
//		}
//		return rowsAffected, nil
//	}
func GenerateDeleteConfigs(modelName string, modelDBName string, deleteConfig []defs.AccessConfig) ([]NamedQuery, []*golang.FunctionDef, []*golang.StructDef, error) {
	functions := make([]*golang.FunctionDef, 0, len(deleteConfig))
	reqs := make([]*golang.StructDef, 0, len(deleteConfig))
	queries := make([]NamedQuery, 0, len(deleteConfig))

	for _, conf := range deleteConfig {
		query, paramRefs := datahelpers.MakeDeleteQuery(modelName, &conf)
		queries = append(queries, NamedQuery{Name: conf.Name, Query: query})
		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)
		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)
		fn := DeleteCodeFunction(conf.Name, modelDBName)
		functions = append(functions, fn)
	}
	return queries, functions, reqs, nil
}

func SetupDatabaseFunction(dataConf defs.DataConfig) (*golang.FunctionDef, error) {

	returnParams := typeOnlyParamsCE("error")
	fn := &golang.FunctionDef{
		Name:    "SetupDatabase",
		Imports: []string{"database/sql", "github.com/lib/pq", "time"},
		Returns: returnParams,
		Body: golang.CodeElements{
			goutils.FCEHNewOutCE([]string{"db", "err"}, "SetupDBConnection", goutils.EHError("err")),
			returnValuesCE("nil"),
		},
	}
	return fn, nil
}
