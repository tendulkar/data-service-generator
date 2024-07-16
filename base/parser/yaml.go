package parser

import (
	"os"

	"gopkg.in/yaml.v3"
	"stellarsky.ai/platform/codegen/data-service-generator/base"
)

func ReadYamlTo[T any](filename string) (*T, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		base.LOG.Error("Error reading yaml file", "error", err, "filename", filename)
		return nil, err
	}
	var t T
	err = yaml.Unmarshal(data, &t)
	if err != nil {
		base.LOG.Error("Error parsing yaml", "error", err, "filename", filename)
		return &t, err
	}
	return &t, nil
}

func StringToYaml[T any](s string) (*T, error) {
	var t T
	err := yaml.Unmarshal([]byte(s), &t)
	if err != nil {
		base.LOG.Error("Error parsing yaml", "error", err, "string", s)
		return &t, err
	}
	return &t, nil
}
