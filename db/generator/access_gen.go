package generator

import (
	"fmt"
	"maps"
	"os"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"stellarsky.ai/platform/codegen/data-service-generator/base"
	"stellarsky.ai/platform/codegen/data-service-generator/config"
	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang"
	datahelpers "stellarsky.ai/platform/codegen/data-service-generator/db/generator/data-helpers"
	"stellarsky.ai/platform/codegen/data-service-generator/db/generator/defs"
	"stellarsky.ai/platform/codegen/data-service-generator/db/models"
)

func Generate(config defs.ModelConfig) error {
	tmpl, err := template.New("model").Funcs(template.FuncMap{
		"Args": Args,
		"Join": Join,
		// "WhereClause":         WhereClause,
		"AttributeNames":      AttributeNames,
		"AttributeValues":     AttributeValues,
		"SetClause":           SetClause,
		"ScanArgs":            ScanArgs,
		"ApplyTransformation": ApplyTransformation,
	}).Parse(modelTemplate)
	if err != nil {
		return err
	}

	// Create output file
	file, err := os.Create(fmt.Sprintf("generated/%s_gen.go", config.Model.Name))
	if err != nil {
		return err
	}
	defer file.Close()

	// Execute template with config data
	err = tmpl.Execute(file, config)
	if err != nil {
		return err
	}

	return nil

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
	base.LOG.Info("ReadTypeValidations", "goType", goType, "validationIds", validationIds,
		"attributeId", attributeId, "typeId", typeId, "postgresType", postgresType, "goTypeStr", goTypeStr)
	validations := datahelpers.GetValidations(validationIds)
	return attribute.Name, goType, validations, nil
}

func generateModel(config *defs.ModelConfig) ([]*golang.Struct, []*golang.Function, error) {

	models := make([]*golang.Struct, 0, 1)
	functions := make([]*golang.Function, 0, 1)

	modelFields := make([]*golang.Field, 0, 1)
	titleCaser := cases.Title(language.English)
	for _, attribute := range config.Model.Attributes {
		attrName, goType, _, err := readTypeAndValidations(attribute)
		if err != nil {
			return nil, nil, err
		}
		base.LOG.Info("Attribute", "attrName", attrName, "goType", goType)
		modelFields = append(modelFields, &golang.Field{
			Name: titleCaser.String(attrName),
			Type: *goType,
			Tag:  fmt.Sprintf(`json:"%s"`, attrName),
		})
	}

	modelStruct := &golang.Struct{
		Name:   fmt.Sprintf("%sModel", config.Model.Name),
		Fields: modelFields,
	}

	models = append(models, modelStruct)

	return models, functions, nil

}

func GenerateV2(config defs.ModelConfig) error {
	// for i := 0; i < len(conf); i++ {
	// findConfigs := config.Access[i].Find
	// updateConfigs := config.Access[i].Update
	// addConfigs := config.Access[i].Add
	// addOrReplaceConfigs := config.Access[i].AddOrReplace
	// deleteConfigs := config.Access[i].Delete
	// GenerateFindConfigs(findConfigs)
	// GenerateUpdateConfigs(updateConfigs)
	// GenerateAddConfigs(addConfigs)
	// GenerateAddOrReplaceConfigs(addOrReplaceConfigs)
	// GenerateDeleteConfigs(deleteConfigs)
	// }

	allQueries := make(map[string]string)
	allFunctions := make([]*golang.Function, 0)
	allStructs := make([]*golang.Struct, 0)

	models, fns, err := generateModel(&config)
	if err != nil {
		return err
	}
	allStructs = append(allStructs, models...)
	allFunctions = append(allFunctions, fns...)

	caser := cases.Title(language.English)
	modelName := caser.String(config.Model.Name)

	queries, accessFns, structs, err := GenerateFindConfigs(modelName, config.Access.Find)
	if err != nil {
		return err
	}
	maps.Copy(allQueries, queries)
	allFunctions = append(allFunctions, accessFns...)
	allStructs = append(allStructs, structs...)

	queries, accessFns, structs, err = GenerateUpdateConfigs(modelName, config.Access.Update)
	if err != nil {
		return err
	}
	maps.Copy(allQueries, queries)
	allFunctions = append(allFunctions, accessFns...)
	allStructs = append(allStructs, structs...)

	queries, accessFns, structs, err = GenerateAddConfigs(modelName, config.Access.Add)
	if err != nil {
		return err
	}
	maps.Copy(allQueries, queries)
	allFunctions = append(allFunctions, accessFns...)
	allStructs = append(allStructs, structs...)

	queries, accessFns, structs, err = GenerateAddOrReplaceConfigs(modelName, config.Access.AddOrReplace)
	if err != nil {
		return err
	}
	maps.Copy(allQueries, queries)
	allFunctions = append(allFunctions, accessFns...)
	allStructs = append(allStructs, structs...)

	queries, accessFns, structs, err = GenerateDeleteConfigs(modelName, config.Access.Delete)
	if err != nil {
		return err
	}
	maps.Copy(allQueries, queries)
	allFunctions = append(allFunctions, accessFns...)
	allStructs = append(allStructs, structs...)

	goSrc := golang.GoSourceFile{
		Package:      "database",
		Structs:      allStructs,
		Functions:    allFunctions,
		InitFunction: nil,
		Variables:    nil,
		Constants:    nil}
	srcCode, deps, err := goSrc.SourceCode()
	if err != nil {
		return err
	}
	fmt.Println("GenerateV2 Source code", srcCode, deps)
	// base.LOG.Info("Source code", "source", goSrc.SourceCode())

	return nil
}

func generateParamsStruct(paramRefs []defs.ParameterRef, name string) *golang.Struct {

	paramFileds := make([]*golang.Field, 0, len(paramRefs))
	for _, param := range paramRefs {
		paramFileds = append(paramFileds, &golang.Field{
			Name: param.Name,
			Type: golang.GoInterfaceType,
			Tag:  fmt.Sprintf(`json:"%s"`, param.Name),
		})
	}

	paramsStruct := &golang.Struct{
		Name:   fmt.Sprintf("%sParams", name),
		Fields: paramFileds,
	}

	return paramsStruct
}

func generateRequestStruct(name string, paramStructName string) *golang.Struct {
	return &golang.Struct{
		Name: name,
		Fields: []*golang.Field{
			{Name: "Params", Type: golang.GoType{Name: paramStructName}, Tag: "json:\"params\""},
		},
	}
}

func generateAccessStructs(paramRef []defs.ParameterRef, name string) []*golang.Struct {
	paramStruct := generateParamsStruct(paramRef, name)
	reqStruct := generateRequestStruct(name, paramStruct.Name)
	return []*golang.Struct{paramStruct, reqStruct}
}

func GenerateFindConfigs(modelName string, findConfig []defs.AccessConfig) (map[string]string, []*golang.Function, []*golang.Struct, error) {

	functions := make([]*golang.Function, 0, len(findConfig))
	reqs := make([]*golang.Struct, 0, len(findConfig))
	queries := make(map[string]string)

	for _, conf := range findConfig {

		query, paramRefs := datahelpers.MakeFindQuery(modelName, &conf)
		queries[conf.Name] = query
		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)

		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)

		fn := FindCodeFunction(modelName, conf.Name, conf.Attributes)
		functions = append(functions, fn)
	}

	return queries, functions, reqs, nil

}

func GenerateUpdateConfigs(modelName string, updateConfig []defs.AccessConfig) (map[string]string, []*golang.Function, []*golang.Struct, error) {

	functions := make([]*golang.Function, 0, len(updateConfig))
	reqs := make([]*golang.Struct, 0, len(updateConfig))
	queries := make(map[string]string)

	for _, conf := range updateConfig {

		query, paramRefs := datahelpers.MakeUpdateQuery(modelName, &conf)
		queries[conf.Name] = query

		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)
		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)

		fn := UpdateCodeFunction(conf.Name)
		functions = append(functions, fn)
	}

	return queries, functions, reqs, nil

}

func GenerateAddConfigs(modelName string, addConfig []defs.AccessConfig) (map[string]string, []*golang.Function, []*golang.Struct, error) {
	functions := make([]*golang.Function, 0, len(addConfig))
	reqs := make([]*golang.Struct, 0, len(addConfig))
	queries := make(map[string]string)

	for _, conf := range addConfig {
		query, paramRefs := datahelpers.MakeAddQuery(modelName, &conf)
		queries[conf.Name] = query
		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)
		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)
		fn := AddCodeFunction(conf.Name)
		functions = append(functions, fn)
	}
	return queries, functions, reqs, nil

}

func GenerateAddOrReplaceConfigs(modelName string, addOrReplaceConfig []defs.AccessConfig) (map[string]string, []*golang.Function, []*golang.Struct, error) {
	functions := make([]*golang.Function, 0, len(addOrReplaceConfig))
	reqs := make([]*golang.Struct, 0, len(addOrReplaceConfig))
	queries := make(map[string]string)

	for _, conf := range addOrReplaceConfig {
		query, paramRefs := datahelpers.MakeAddOrReplaceQuery(modelName, &conf)
		queries[conf.Name] = query
		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)
		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)
		fn := AddOrReplaceCodeFunction(conf.Name)
		functions = append(functions, fn)
	}
	return queries, functions, reqs, nil
}

func GenerateDeleteConfigs(modelName string, deleteConfig []defs.AccessConfig) (map[string]string, []*golang.Function, []*golang.Struct, error) {
	functions := make([]*golang.Function, 0, len(deleteConfig))
	reqs := make([]*golang.Struct, 0, len(deleteConfig))
	queries := make(map[string]string)

	for _, conf := range deleteConfig {
		query, paramRefs := datahelpers.MakeDeleteQuery(modelName, &conf)
		queries[conf.Name] = query
		accessStructs := generateAccessStructs(paramRefs, conf.Name)
		reqs = append(reqs, accessStructs...)
		paramFn := ReadParamsFunction(paramRefs, conf.Name, "values", "params")
		functions = append(functions, paramFn)
		fn := DeleteCodeFunction(conf.Name)
		functions = append(functions, fn)
	}
	return queries, functions, reqs, nil
}
