package generator

import (
	"fmt"
	"os"
	"text/template"

	"stellarsky.ai/platform/codegen/data-service-generator/db/generator/defs"
)

func Generate(config defs.ModelConfig) error {
	tmpl, err := template.New("model").Funcs(template.FuncMap{
		"Args":                Args,
		"Join":                Join,
		"WhereClause":         WhereClause,
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
