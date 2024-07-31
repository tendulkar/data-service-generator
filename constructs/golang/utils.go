package golang

import (
	"stellarsky.ai/platform/codegen/data-service-generator/base"
)

type NameWithType struct {
	Name string
	Type *GoType
}

func GenerateStructForDataModel(modelName string, nameWithTypes []NameWithType, addJsonTag bool, addYamlTag bool, addDBTag bool) *Struct {
	modelName = ToPascalCase(modelName)
	fields := make([]*Field, 0, len(nameWithTypes))
	for _, nameWithType := range nameWithTypes {
		fieldName := ToPascalCase(nameWithType.Name)
		fields = append(fields, &Field{Name: fieldName, Type: nameWithType.Type, AddJsonTag: addJsonTag, AddYamlTag: addYamlTag, AddDBTag: addDBTag})
	}
	base.LOG.Debug("Generating struct for JSON:", "model", modelName, "fields", fields)
	return &Struct{Name: modelName, Fields: fields}
}

func FunctionCallCE(newOutput, output interface{}, receiver, functionName string, args interface{},
	isAsync, isDefer bool, errHandler *ErrorHandler, cleanHandler *CleanningHandler) *CodeElement {
	return &CodeElement{
		FunctionCall: &FunctionCall{
			Output:           output,
			NewOutput:        newOutput,
			Receiver:         receiver,
			Function:         functionName,
			Args:             args,
			Async:            isAsync,
			Defer:            isDefer,
			ErrorHandler:     errHandler,
			CleanningHandler: cleanHandler,
		},
	}
}

func NewCleanningHandler(receiver string, funcName string, args interface{}, steps CodeElements) *CleanningHandler {
	return &CleanningHandler{
		Receiver: receiver,
		Function: funcName,
		Args:     args,
		Steps:    steps,
	}
}

func NewLits(values ...interface{}) []*Literal {
	literals := make([]*Literal, 0, len(values))
	for _, value := range values {
		literals = append(literals, &Literal{Value: value})
	}
	return literals
}
