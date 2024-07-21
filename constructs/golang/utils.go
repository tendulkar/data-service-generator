package golang

import (
	"stellarsky.ai/platform/codegen/data-service-generator/base"
)

type NameWithType struct {
	Name string
	Type *GoType
}

func GenerateStructForJSON(modelName string, nameWithTypes []NameWithType) *Struct {
	modelName = ToPascalCase(modelName)
	fields := make([]*Field, 0, len(nameWithTypes))
	for _, nameWithType := range nameWithTypes {
		fieldName := ToPascalCase(nameWithType.Name)
		fields = append(fields, &Field{Name: fieldName, Type: nameWithType.Type, AddJsonTag: true})
	}
	base.LOG.Info("Generating struct for JSON:", "model", modelName, "fields", fields)
	return &Struct{Name: modelName, Fields: fields}

}
