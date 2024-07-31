package golang

import (
	"fmt"

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

func GenerateStructWithNewFunction(structName string, structSuffix string, nameWithTypes []NameWithType, addJsonTag bool, addYamlTag bool, addDBTag bool) (*Struct, *Function) {

	fields := make([]*Field, 0, len(nameWithTypes))
	keyValues := make(KeyValues, 0, len(nameWithTypes))
	fnParams := make([]*Parameter, 0, len(nameWithTypes))
	for _, nameWithType := range nameWithTypes {
		fieldName := ToPascalCase(nameWithType.Name)
		fieldCamelName := ToCamelCase(nameWithType.Name)
		fields = append(fields, &Field{Name: fieldName, Type: nameWithType.Type, AddJsonTag: addJsonTag, AddYamlTag: addYamlTag, AddDBTag: addDBTag})
		keyValues = append(keyValues, &KeyValue{Key: fieldName, Variable: fieldCamelName})
		fnParams = append(fnParams, &Parameter{Name: fieldCamelName, Type: nameWithType.Type})
	}

	structFormatName := ToPascalCase(structName) + structSuffix
	st := &Struct{
		Name:   structFormatName,
		Fields: fields,
	}

	stCreate := &StructCreation{
		StructType: st.Name,
		KeyValues:  keyValues,
	}
	fn := &Function{
		Name:       "New" + structFormatName,
		Parameters: fnParams,
		Returns:    NewReturnTypes(fmt.Sprintf("*%s", structFormatName)),
		Body: CodeElements{
			NewReturnStatement(stCreate),
		},
	}
	return st, fn
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

func NewField(name string, typeName string) *Field {
	return &Field{Name: name, Type: &GoType{Name: typeName, Source: ""}}
}

func NewFieldWithSource(name string, typeName string, source string) *Field {
	return &Field{Name: name, Type: &GoType{Name: typeName, Source: source}}
}

func NewParameter(name string, typeName string) *Parameter {
	return &Parameter{Name: name, Type: &GoType{Name: typeName, Source: ""}}
}

func NewParameterWithSource(name string, typeName string, source string) *Parameter {
	return &Parameter{Name: name, Type: &GoType{Name: typeName, Source: source}}
}

func NewReturnTypes(types ...string) []*Parameter {
	if len(types) == 0 {
		return nil
	}
	returns := make([]*Parameter, 0, len(types))
	for _, typeName := range types {
		returns = append(returns, &Parameter{Type: &GoType{Name: typeName, Source: ""}})
	}
	return returns
}

func NewReturnStatement(code CodeBlock) *CodeElement {
	return &CodeElement{Return: code}
}
