package generator

import (
	"bytes"
	"fmt"
	"os"
	"strings"
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

func readTypeValidations(attributeId int64) (string, *golang.GoType, []*models.Validation, error) {
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
		attrName, goType, _, err := readTypeValidations(attribute)
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

	models, _, err := generateModel(&config)
	if err != nil {
		return err
	}
	goSrc := golang.GoSourceFile{
		Package:      "database",
		Structs:      models,
		Functions:    nil,
		InitFunction: nil,
		Variables:    nil,
		Constants:    nil}
	srcCode, err := goSrc.SourceCode()
	if err != nil {
		return err
	}
	fmt.Println("GenerateV2 Source code", srcCode)
	// base.LOG.Info("Source code", "source", goSrc.SourceCode())
	GenerateFindConfigs(strings.Title(config.Model.Name), config.Access.Find)

	return nil
}

func generateContextDBFunction(data any, tmpl *template.Template, modelName string, name string) (*golang.Function, error) {
	buff := new(bytes.Buffer)
	err := tmpl.Execute(buff, data)
	if err != nil {
		return nil, err
	}
	fn := &golang.Function{
		Name: name,
		Parameters: []*golang.Parameter{
			{Name: "ctx", Type: golang.GoType{Name: "context.Context", Source: "context"}},
			{Name: "db", Type: golang.GoType{Name: fmt.Sprintf("*%sDB", modelName)}},
			{Name: "request", Type: golang.GoType{Name: fmt.Sprintf("*%sRequest", name)}},
		},
		Returns: []*golang.GoType{
			{Name: fmt.Sprintf("*%sResponse", name)},
			{Name: "error"},
		},
		Body: golang.GoCodeBlock{CodeBlock: buff.String(),
			Sources: []string{},
		},
	}

	fmt.Printf("Function code generated: %v\n", buff.String())
	return fn, nil
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

func GenerateFindConfigs(modelName string, findConfig []defs.AccessConfig) ([]*golang.Function, []*golang.Struct, error) {
	tmpl, err := template.New("find").Funcs(template.FuncMap{
		"Args":                Args,
		"Join":                Join,
		"AttributeNames":      AttributeNames,
		"AttributeValues":     AttributeValues,
		"SetClause":           SetClause,
		"ScanArgs":            ScanArgs,
		"ApplyTransformation": ApplyTransformation,
	}).Parse(findTemplate)

	if err != nil {
		return nil, nil, err
	}

	readParamsTmpl, err := template.New("params").Parse(readParamsToValues)

	if err != nil {
		return nil, nil, err
	}

	functions := make([]*golang.Function, 0, len(findConfig))
	reqs := make([]*golang.Struct, 0, len(findConfig))

	for _, findConf := range findConfig {

		query, paramRefs := datahelpers.MakeFindQuery(modelName, &findConf)
		paramsStruct := generateParamsStruct(paramRefs, findConf.Name)

		fmt.Println("query", query)
		req := &golang.Struct{
			Name: fmt.Sprintf("%sRequest", findConf.Name),
			Fields: []*golang.Field{
				{Name: "Params", Type: golang.GoType{Name: paramsStruct.Name}, Tag: "json:\"params\""},
			},
		}

		reqs = append(reqs, paramsStruct)
		reqs = append(reqs, req)

		data := struct {
			Name           string
			ScanAttributes string
			ModelName      string
		}{
			ModelName:      modelName,
			Name:           findConf.Name,
			ScanAttributes: ScanArgs(findConf.Attributes),
		}

		fn, err := generateContextDBFunction(data, tmpl, modelName, findConf.Name)
		if err != nil {
			return functions, reqs, err
		}
		functions = append(functions, fn)

		paramsData := struct {
			ParameterRefs []defs.ParameterRef
		}{
			ParameterRefs: paramRefs,
		}
		readParamsCode := new(bytes.Buffer)
		err = readParamsTmpl.Execute(readParamsCode, paramsData)
		if err != nil {
			return functions, reqs, err
		}

		paramFn := &golang.Function{
			Name: fmt.Sprintf("%sReadParams", findConf.Name),
			Parameters: []*golang.Parameter{
				{Name: "request", Type: golang.GoType{Name: fmt.Sprintf("*%sRequest", findConf.Name)}},
			},
			Returns: []*golang.GoType{&golang.GoInterfaceArrayType},
			Body: golang.GoCodeBlock{CodeBlock: readParamsCode.String(),
				Sources: []string{},
			},
		}

		functions = append(functions, paramFn)
		fmt.Println(readParamsCode.String())
	}

	return functions, reqs, nil

}

func GenerateUpdateConfigs(modelName string, updateConfig []defs.AccessConfig) ([]*golang.Function, error) {
	template, err := template.New("update").Funcs(template.FuncMap{
		"Args":                Args,
		"Join":                Join,
		"AttributeNames":      AttributeNames,
		"AttributeValues":     AttributeValues,
		"SetClause":           SetClause,
		"ScanArgs":            ScanArgs,
		"ApplyTransformation": ApplyTransformation,
	}).Parse(updateTemplate)

	if err != nil {
		return nil, err
	}

	functions := make([]*golang.Function, 0, len(updateConfig))
	for _, updateConf := range updateConfig {

		data := struct {
			Name           string
			ScanAttributes string
			ModelName      string
		}{
			ModelName:      modelName,
			Name:           updateConf.Name,
			ScanAttributes: ScanArgs(updateConf.Attributes),
		}

		buff := new(bytes.Buffer)
		err := template.Execute(buff, data)
		if err != nil {
			return nil, err
		}
		fn := &golang.Function{
			Name: updateConf.Name,
			Parameters: []*golang.Parameter{
				{Name: "ctx", Type: golang.GoType{Name: "context.Context", Source: "context"}},
				{Name: "db", Type: golang.GoType{Name: fmt.Sprintf("*%sDB", modelName)}},
				{Name: "request", Type: golang.GoType{Name: fmt.Sprintf("*%sRequest", updateConf.Name)}},
			},
			Returns: []*golang.GoType{
				&golang.GoErrorType,
			},
			Body: golang.GoCodeBlock{
				CodeBlock: buff.String(),
				Sources:   []string{},
			},
		}

		functions = append(functions, fn)
		fmt.Println(buff.String())

	}

	return functions, nil
}
