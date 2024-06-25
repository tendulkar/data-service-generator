package generator

import (
	"fmt"
	"os"
	"text/template"
)

type Model struct {
	ID         int    `yaml:"id"`
	Namespace  string `yaml:"namespace"`
	Family     string `yaml:"family"`
	Name       string `yaml:"name"`
	Attributes []struct {
		ID   int    `yaml:"id"`
		Name string `yaml:"name"`
		Type string `yaml:"type"`
	} `yaml:"attributes"`
	UniqueConstraints []struct {
		ConstraintName string `yaml:"constraint_name"`
		Attributes     []int  `yaml:"attributes"`
	} `yaml:"unique_constraints"`
}

type Access struct {
	Find         []AccessConfig `yaml:"find"`
	Update       []AccessConfig `yaml:"update"`
	Add          []AccessConfig `yaml:"add"`
	AddOrReplace []AccessConfig `yaml:"add_or_replace"`
	Delete       []AccessConfig `yaml:"delete"`
}

type ModelConfig struct {
	Model  `yaml:"model"`
	Access `yaml:"access"`
}

type Parameter struct {
	Attribute string `yaml:"attribute"`
}

type Request struct {
	Parameters []Parameter `yaml:"parameters"`
}

type AccessConfig struct {
	Name             string   `yaml:"name"`
	Request          Request  `yaml:"request,omitempty"`
	Attributes       []string `yaml:"attributes,omitempty"`
	Filter           []Filter `yaml:"filter,omitempty"`
	Set              []Update `yaml:"set,omitempty"`
	Autoincrement    []string `yaml:"autoincrement,omitempty"`
	CaptureTimestamp []string `yaml:"capture_timestamp,omitempty"`
	Values           []Update `yaml:"values,omitempty"`
}

type Update struct {
	Attribute string      `yaml:"attribute"`
	Value     interface{} `yaml:"value"`
}

type Attribute struct {
	Attribute string `yaml:"attribute"`
}

func Generate(config ModelConfig) error {
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
