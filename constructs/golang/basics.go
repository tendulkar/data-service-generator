package golang

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
)

// GoType represents a type in Go with a name and an optional import source.
type GoType struct {
	Name   string // Name of the type, could be a funciton too. example: func(logger *zap.Logger) (*zap.Logger, error)
	Source string // Import path, empty if it's a built-in type
}

var (
	GoBoolType                  = GoType{Name: "bool"}
	GoIntType                   = GoType{Name: "int"}
	GoInt8Type                  = GoType{Name: "int8"}
	GoInt16Type                 = GoType{Name: "int16"}
	GoInt32Type                 = GoType{Name: "int32"}
	GoInt64Type                 = GoType{Name: "int64"}
	GoUintType                  = GoType{Name: "uint"}
	GoUint8Type                 = GoType{Name: "uint8"}
	GoUint16Type                = GoType{Name: "uint16"}
	GoUint32Type                = GoType{Name: "uint32"}
	GoUint64Type                = GoType{Name: "uint64"}
	GoFloat32Type               = GoType{Name: "float32"}
	GoFloat64Type               = GoType{Name: "float64"}
	GoComplex64Type             = GoType{Name: "complex64"}
	GoStringType                = GoType{Name: "string"}
	GoErrorType                 = GoType{Name: "error"}
	GoInterfaceType             = GoType{Name: "interface{}"}
	GoAnyType                   = GoType{Name: "any"}
	GoInterfaceArrayType        = GoType{Name: "[]interface{}"}
	GoBoolArrayType             = GoType{Name: "[]bool"}
	GoIntArrayType              = GoType{Name: "[]int"}
	GoInt8ArrayType             = GoType{Name: "[]int8"}
	GoInt16ArrayType            = GoType{Name: "[]int16"}
	GoInt32ArrayType            = GoType{Name: "[]int32"}
	GoInt64ArrayType            = GoType{Name: "[]int64"}
	GoUintArrayType             = GoType{Name: "[]uint"}
	GoUint8ArrayType            = GoType{Name: "[]uint8"}
	GoUint16ArrayType           = GoType{Name: "[]uint16"}
	GoUint32ArrayType           = GoType{Name: "[]uint32"}
	GoUint64ArrayType           = GoType{Name: "[]uint64"}
	GoFloat32ArrayType          = GoType{Name: "[]float32"}
	GoFloat64ArrayType          = GoType{Name: "[]float64"}
	GoComplex64ArrayType        = GoType{Name: "[]complex64"}
	GoStringArrayType           = GoType{Name: "[]string"}
	GoErrorArrayType            = GoType{Name: "[]error"}
	GoStringBoolMapType         = GoType{Name: "map[string]bool"}
	GoStringIntMapType          = GoType{Name: "map[string]int"}
	GoStringInt8MapType         = GoType{Name: "map[string]int8"}
	GoStringInt16MapType        = GoType{Name: "map[string]int16"}
	GoStringInt32MapType        = GoType{Name: "map[string]int32"}
	GoStringInt64MapType        = GoType{Name: "map[string]int64"}
	GoStringUintMapType         = GoType{Name: "map[string]uint"}
	GoStringUint8MapType        = GoType{Name: "map[string]uint8"}
	GoStringUint16MapType       = GoType{Name: "map[string]uint16"}
	GoStringUint32MapType       = GoType{Name: "map[string]uint32"}
	GoStringUint64MapType       = GoType{Name: "map[string]uint64"}
	GoStringFloat64MapType      = GoType{Name: "map[string]float64"}
	GoStringStringMapType       = GoType{Name: "map[string]string"}
	GoStringInterfaceMapType    = GoType{Name: "map[string]interface{}"}
	GoInterfaceInterfaceMapType = GoType{Name: "map[interface{}]interface{}"}
)

// typeMap maps type names to their corresponding GoType structs
var typeMap = map[string]GoType{
	"bool":                        GoBoolType,
	"int":                         GoIntType,
	"int8":                        GoInt8Type,
	"int16":                       GoInt16Type,
	"int32":                       GoInt32Type,
	"int64":                       GoInt64Type,
	"uint":                        GoUintType,
	"uint8":                       GoUint8Type,
	"uint16":                      GoUint16Type,
	"uint32":                      GoUint32Type,
	"uint64":                      GoUint64Type,
	"float32":                     GoFloat32Type,
	"float64":                     GoFloat64Type,
	"complex64":                   GoComplex64Type,
	"string":                      GoStringType,
	"error":                       GoErrorType,
	"interface{}":                 GoInterfaceType,
	"any":                         GoAnyType,
	"[]interface{}":               GoInterfaceArrayType,
	"[]bool":                      GoBoolArrayType,
	"[]int":                       GoIntArrayType,
	"[]int8":                      GoInt8ArrayType,
	"[]int16":                     GoInt16ArrayType,
	"[]int32":                     GoInt32ArrayType,
	"[]int64":                     GoInt64ArrayType,
	"[]uint":                      GoUintArrayType,
	"[]uint8":                     GoUint8ArrayType,
	"[]uint16":                    GoUint16ArrayType,
	"[]uint32":                    GoUint32ArrayType,
	"[]uint64":                    GoUint64ArrayType,
	"[]float32":                   GoFloat32ArrayType,
	"[]float64":                   GoFloat64ArrayType,
	"[]complex64":                 GoComplex64ArrayType,
	"[]string":                    GoStringArrayType,
	"[]error":                     GoErrorArrayType,
	"map[string]bool":             GoStringBoolMapType,
	"map[string]int":              GoStringIntMapType,
	"map[string]int8":             GoStringInt8MapType,
	"map[string]int16":            GoStringInt16MapType,
	"map[string]int32":            GoStringInt32MapType,
	"map[string]int64":            GoStringInt64MapType,
	"map[string]uint":             GoStringUintMapType,
	"map[string]uint8":            GoStringUint8MapType,
	"map[string]uint16":           GoStringUint16MapType,
	"map[string]uint32":           GoStringUint32MapType,
	"map[string]uint64":           GoStringUint64MapType,
	"map[string]float64":          GoStringFloat64MapType,
	"map[string]string":           GoStringStringMapType,
	"map[string]interface{}":      GoStringInterfaceMapType,
	"map[interface{}]interface{}": GoInterfaceInterfaceMapType,
}

// findGoType takes a Go type name as a string and returns the corresponding GoType.
func TranslateToGoType(typeName string) (*GoType, error) {
	if goType, ok := typeMap[typeName]; ok {
		return &goType, nil
	}
	return &GoType{}, fmt.Errorf("type %s not found", typeName)
}

// Parameter represents a function parameter, using GoType.
type Parameter struct {
	Name string
	Type GoType
}

// GoCodeBlock represents a block of code with optional additional imports.
type GoCodeBlock struct {
	CodeBlock string
	Sources   []string // Additional import paths required by the code block
}

type Variable struct {
	Name  string
	Type  GoType
	Value string
}

type Constant struct {
	Name  string
	Type  GoType
	Value string
}

// Function represents a Go function definition.
type Function struct {
	Name       string
	Parameters []*Parameter
	Returns    []*GoType
	Body       GoCodeBlock
	Receiver   *Receiver // Nil if not a member function
}

// Receiver represents the receiver of a Go method.
type Receiver struct {
	Name string
	Type GoType
}

type Field struct {
	Name string
	Type GoType
	Tag  string
}

// Struct represents a Go struct with member functions.
type Struct struct {
	Name      string
	Fields    []*Field // Using Parameter as it has both name and type
	Functions []*Function
}

// StructCode generates the Go code for the struct, including its member functions.
func (s Struct) StructCode() (string, map[string]bool) {
	fieldStrs := make([]string, len(s.Fields))
	for i, field := range s.Fields {
		typeName := field.Type.Name
		if !strings.HasPrefix(typeName, "*") && field.Type.Source != "" { // Assume non-primitive types need pointers
			typeName = "*" + typeName
		}
		tagSpace := Indent
		tagName := fmt.Sprintf("`json:\"%s\"`", field.Name)
		if field.Tag == "" {
			tagSpace = ""
			tagName = ""
		}
		fieldStrs[i] = fmt.Sprintf("%s%s %s%s%s", Indent, field.Name, typeName, tagSpace, tagName)
	}

	structDef := fmt.Sprintf("type %s struct {\n%s\n}", s.Name, strings.Join(fieldStrs, "\n"))

	var funcDefs []string
	allSources := make(map[string]bool)
	for _, fn := range s.Functions {
		fnCode, fnSources := fn.FunctionCode()
		for source := range fnSources {
			allSources[source] = true
		}
		funcDefs = append(funcDefs, fnCode)
	}
	for src := range gatherSources(nil, s.Fields, nil) {
		allSources[src] = true
	}

	return fmt.Sprintf("%s\n%s", structDef, strings.Join(funcDefs, "\n\n")), allSources
}

// FunctionCode generates the Go code for the function, including necessary imports.
func (f Function) FunctionCode() (string, map[string]bool) {
	// Generate parameters string
	var params []string
	for _, param := range f.Parameters {
		paramType := param.Type.Name
		// if !strings.HasPrefix(paramType, "*") && param.Type.Source != "" { // Assume non-primitive types need pointers
		// 	paramType = "*" + paramType
		// }
		params = append(params, fmt.Sprintf("%s %s", param.Name, paramType))
	}
	paramStr := strings.Join(params, ", ")

	// Generate return types string
	var returns []string
	for _, ret := range f.Returns {
		retType := ret.Name
		if !strings.HasPrefix(retType, "*") && ret.Source != "" {
			retType = "*" + retType
		}
		returns = append(returns, retType)
	}
	returnStr := formatReturnTypes(returns)

	// Generate all required import statements
	allImports := gatherSources(f.Parameters, nil, f.Returns)
	for _, src := range f.Body.Sources {
		allImports[src] = true
	}

	receiver := ""
	if f.Receiver != nil {
		receiverType := f.Receiver.Type.Name
		if !strings.HasPrefix(receiverType, "*") { // Member functions always on pointer receivers
			receiverType = "*" + receiverType
		}
		receiver = fmt.Sprintf("(%s %s) ", f.Receiver.Name, receiverType)
	}

	body := IndentCode(f.Body.CodeBlock, 1)

	returnSpace := " "
	if returnStr == "" {
		returnSpace = ""
	}
	return fmt.Sprintf("func %s%s(%s) %s%s{\n%s\n}", receiver, f.Name, paramStr, returnStr, returnSpace, body), allImports
}

// Imports are automatically derived from each block,
// Variables, Constants, Structs, Functions, and InitFunction have their own sources, automatically gets added to allImports
type GoSourceFile struct {
	Package      string
	Variables    []*Variable
	Constants    []*Constant
	Structs      []*Struct
	Functions    []*Function
	InitFunction *GoCodeBlock
}

func (s *GoSourceFile) SourceCode() (string, error) {
	return GenerateGoFile(s.Package, s.Structs, s.Functions, s.Variables, s.Constants, s.InitFunction)
}

// gatherSources collects import sources from parameters and returns.
func gatherSources(params []*Parameter, fields []*Field, returns []*GoType) map[string]bool {
	uniqueSources := make(map[string]bool)
	for _, param := range params {
		if param.Type.Source != "" {
			uniqueSources[param.Type.Source] = true
		}
	}
	for _, field := range fields {
		if field.Type.Source != "" {
			uniqueSources[field.Type.Source] = true
		}
	}
	for _, ret := range returns {
		if ret.Source != "" {
			uniqueSources[ret.Source] = true
		}
	}

	return uniqueSources
}

// func deduplicate(sources []string) []string {
// 	set := make(map[string]bool)
// 	for _, source := range sources {
// 		set[source] = true
// 	}
// 	var dedup []string
// 	for source := range set {
// 		dedup = append(dedup, source)
// 	}
// 	return dedup
// }

// generateImports creates unique import statements from a list of import paths.
func generateImports(sources map[string]bool) string {
	if len(sources) == 0 {
		return ""
	}

	var importLines []string
	for source := range sources {
		importLines = append(importLines, fmt.Sprintf("%s\"%s\"", Indent, source))
	}

	importCode := fmt.Sprintf("import (\n%s\n)\n", strings.Join(importLines, "\n"))
	return importCode
}

// formatReturnTypes formats the return types into a string.
func formatReturnTypes(returns []string) string {
	if len(returns) > 1 {
		return fmt.Sprintf("(%s)", strings.Join(returns, ", "))
	} else if len(returns) == 1 {
		return returns[0]
	}
	return ""
}

// GenerateGoFile generates a complete Go source file including the specified package name,
// structs, functions, and standalone functions.
func GenerateGoFile(packageName string, structs []*Struct, functions []*Function, variables []*Variable, constants []*Constant, initFunction *GoCodeBlock) (string, error) {
	var buffer bytes.Buffer

	// Write the package declaration
	buffer.WriteString(fmt.Sprintf("package %s\n\n", packageName))

	// Collect all imports from structs, functions, variables, and constants
	var allImports map[string]bool = make(map[string]bool)
	var structDefinitions []string
	var functionDefinitions []string
	var variableDefinitions []string
	var constantDefinitions []string

	for _, s := range structs {
		collectImports(allImports, nil, s.Fields, nil)
		for _, function := range s.Functions {
			collectFunctionImports(allImports, function)
		}
		structCode, structSources := s.StructCode()
		for source := range structSources {
			allImports[source] = true
		}
		structDefinitions = append(structDefinitions, structCode)
	}

	for _, f := range functions {
		collectFunctionImports(allImports, f)
		fnCode, fnSources := f.FunctionCode()
		for source := range fnSources {
			allImports[source] = true
		}
		functionDefinitions = append(functionDefinitions, fnCode)
	}

	for _, v := range variables {
		variableDefinitions = append(variableDefinitions, fmt.Sprintf("var %s %s = %s", v.Name, v.Type.Name, v.Value))
		if v.Type.Source != "" {
			allImports[v.Type.Source] = true
		}
	}

	for _, c := range constants {
		constantDefinitions = append(constantDefinitions, fmt.Sprintf("const %s %s = %s", c.Name, c.Type.Name, c.Value))
		if c.Type.Source != "" {
			allImports[c.Type.Source] = true
		}
	}

	if initFunction != nil {
		for _, source := range initFunction.Sources {
			allImports[source] = true
		}
	}

	// Write the import statements
	buffer.WriteString(generateImports(allImports))

	// Write global variables and constants
	if len(variableDefinitions) > 0 {
		buffer.WriteString(strings.Join(variableDefinitions, "\n") + "\n\n")
	}
	if len(constantDefinitions) > 0 {
		buffer.WriteString(strings.Join(constantDefinitions, "\n") + "\n\n")
	}

	// Write each struct and its methods
	for _, def := range structDefinitions {
		buffer.WriteString(def)
		buffer.WriteString("\n")
	}

	// Write standalone functions
	for _, def := range functionDefinitions {
		buffer.WriteString(def)
		buffer.WriteString("\n")
	}

	// Write init function if exists
	if initFunction != nil {
		buffer.WriteString("\nfunc init() {\n")
		buffer.WriteString(IndentCode(initFunction.CodeBlock, 1))
		buffer.WriteString("\n}\n")
	}

	srcCode, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", err
	}
	return string(srcCode), nil
}

func collectImports(allImports map[string]bool, parameters []*Parameter, fields []*Field, returns []*GoType) {
	for _, param := range parameters {
		if param.Type.Source != "" {
			allImports[param.Type.Source] = true
		}
	}
	for _, field := range fields {
		if field.Type.Source != "" {
			allImports[field.Type.Source] = true
		}
	}
	for _, ret := range returns {
		if ret.Source != "" {
			allImports[ret.Source] = true
		}
	}
}

func collectFunctionImports(allImports map[string]bool, function *Function) {
	collectImports(allImports, function.Parameters, nil, function.Returns)
	for _, src := range function.Body.Sources {
		allImports[src] = true
	}
	if function.Receiver != nil && function.Receiver.Type.Source != "" {
		allImports[function.Receiver.Type.Source] = true
	}
}
