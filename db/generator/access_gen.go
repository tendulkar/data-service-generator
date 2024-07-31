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

// Generate Model struct for a given model
//
//	Example: User {
//	  ID int `db:"id"`
//	  FullName string `db:"full_name"`
//	  Email string `db:"email"`
//	}
//
//	Example 2: type Product struct {
//		Sku         string  `db:"sku"`
//		ProductName string  `db:"product_name"`
//		Description string  `db:"description"`
//		Price       float64 `db:"price"`
//	}

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

	modelStruct := golang.GenerateStructForDataModel(config.Model.Name, nameWithTypes, false, false, true)
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
	caser := cases.Title(language.English)
	modelName := caser.String(config.Model.Name)

	// Generate Model struct for a given model, for example `type User struct {<fields with db tags>}`
	models, fns, err := generateModel(&config)
	if err != nil {
		return nil, err
	}
	allStructs = append(allStructs, models...)
	allFunctions = append(allFunctions, fns...)

	// PrepareStmt function will prepare all queries for a given model
	prepareFn := PrepareStmtFunction(modelName, allQueries)
	allFunctions = append(allFunctions, prepareFn)

	// Generate methods for SELECT, UPDATE, INSERT, INSERT OR UPDATE, DELETE for a given model
	err = geneateAllAccessMethods(config, modelName, &allQueries, &allFunctions, &allStructs)
	if err != nil {
		return nil, err
	}

	goSrc := &golang.GoSourceFile{
		Package:      "database",
		Structs:      allStructs,
		Functions:    allFunctions,
		InitFunction: nil,
		Variables:    nil,
		Constants:    nil}

	return goSrc, nil
}

// All access methods for a given model (Find, Update, Add, AddOrReplace and Delete),
// will do query on database with above prepared statements (SELECT, UPDATE, INSERT, INSERT OR UPDATE, DELETE)
func geneateAllAccessMethods(config defs.ModelConfig, modelName string, allQueries *[]NamedQuery, allFunctions *[]*golang.Function, allStructs *[]*golang.Struct) error {
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
		err := genAccessFn(modelName, accessConfigs[i], accessMethod, allQueries, allFunctions, allStructs)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateParamsStruct(paramRefs []defs.ParameterRef, name string) *golang.Struct {
	nameWithTypes := make([]golang.NameWithType, 0, len(paramRefs))
	for _, param := range paramRefs {
		nameWithTypes = append(nameWithTypes, golang.NameWithType{
			Name: param.Name,
			Type: golang.GoInterfaceType,
		})
	}

	return golang.GenerateStructForDataModel(fmt.Sprintf("%sParams", name), nameWithTypes, true, false, false)
}

func generateRequestStruct(name string, paramStructName string) *golang.Struct {
	namedWithTypes := []golang.NameWithType{
		{Name: "Params", Type: &golang.GoType{Name: paramStructName}},
	}
	return golang.GenerateStructForDataModel(fmt.Sprintf("%sRequest", name), namedWithTypes, true, false, false)
}

func generateAccessStructs(paramRef []defs.ParameterRef, name string) []*golang.Struct {
	paramStruct := generateParamsStruct(paramRef, name)
	reqStruct := generateRequestStruct(name, paramStruct.Name)
	return []*golang.Struct{paramStruct, reqStruct}
}

// GenerateFindConfigs will generate all SELECT queries for a given model
// It generates find function, and one helper function for reading params from request to bind values to query
// It generates 2 structs for params and request, request is input (arg) to find function, and params is part of request
//
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
//	    func GetProductByID(ctx context.Context, db *sql.DB, requestParams GetProductByIDParams) (results []Product, err error) {
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

// GenerateUpdateConfigs will generate all UPDATE queries for a given model
// Similar to GenerateFindConfigs
// It generates Update function, and one helper function for reading params from request to bind values to query
// It generates 2 structs for params and request, request is input (arg) to Update function, and params is part of request
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

// GenerateAddConfigs will generate all INSERT queries for a given model
// Similar to GenerateFindConfigs
// It generates Add(INSERT) function, and one helper function for reading params from request to bind values to query
// It generates 2 structs for params and request, request is input (arg) to Add function, and params is part of request
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

// GenerateAddOrReplaceConfigs will generate all INSERT OR UPDATE queries for a given model
// Similar to GenerateFindConfigs
// It generates AddOrReplace(INSERT OR UPDATE) function, and one helper function for reading params from request to bind values to query
// It generates 2 structs for params and request, request is input (arg) to AddOrReplace function, and params is part of request
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
