package golang

import (
	"bytes"
	"fmt"
	"strings"
)

// GoType represents a type in Go with a name and an optional import source.
type GoType struct {
	Name   string // Name of the type, could be a funciton too. example: func(logger *zap.Logger) (*zap.Logger, error)
	Source string // Import path, empty if it's a built-in type
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
	Parameters []Parameter
	Returns    []GoType
	Body       GoCodeBlock
	Receiver   *Receiver // Nil if not a member function
}

// Receiver represents the receiver of a Go method.
type Receiver struct {
	Name string
	Type GoType
}

// Struct represents a Go struct with member functions.
type Struct struct {
	Name      string
	Fields    []Parameter // Using Parameter as it has both name and type
	Functions []Function
}

// StructCode generates the Go code for the struct, including its member functions.
func (s Struct) StructCode() (string, map[string]bool) {
	fieldStrs := make([]string, len(s.Fields))
	for i, field := range s.Fields {
		typeName := field.Type.Name
		if !strings.HasPrefix(typeName, "*") && field.Type.Source != "" { // Assume non-primitive types need pointers
			typeName = "*" + typeName
		}
		fieldStrs[i] = fmt.Sprintf("%s%s %s", Indent, field.Name, typeName)
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
	for src := range gatherSources(s.Fields, nil) {
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
		if !strings.HasPrefix(paramType, "*") && param.Type.Source != "" { // Assume non-primitive types need pointers
			paramType = "*" + paramType
		}
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
	allImports := gatherSources(f.Parameters, f.Returns)
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

// gatherSources collects import sources from parameters and returns.
func gatherSources(params []Parameter, returns []GoType) map[string]bool {
	uniqueSources := make(map[string]bool)
	for _, param := range params {
		if param.Type.Source != "" {
			uniqueSources[param.Type.Source] = true
		}
	}
	for _, ret := range returns {
		if ret.Source != "" {
			uniqueSources[ret.Source] = true
		}
	}

	return uniqueSources
}

func deduplicate(sources []string) []string {
	set := make(map[string]bool)
	for _, source := range sources {
		set[source] = true
	}
	var dedup []string
	for source := range set {
		dedup = append(dedup, source)
	}
	return dedup
}

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
func GenerateGoFile(packageName string, structs []Struct, functions []Function, variables []Variable, constants []Constant, initFunction *GoCodeBlock) string {
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
		collectImports(&allImports, s.Fields, nil)
		for _, function := range s.Functions {
			collectFunctionImports(&allImports, function)
		}
		structCode, structSources := s.StructCode()
		for source := range structSources {
			allImports[source] = true
		}
		structDefinitions = append(structDefinitions, structCode)
	}

	for _, f := range functions {
		collectFunctionImports(&allImports, f)
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

	return buffer.String()
}

func collectImports(allImports *map[string]bool, fields []Parameter, returns []GoType) {
	for _, field := range fields {
		if field.Type.Source != "" {
			(*allImports)[field.Type.Source] = true
		}
	}
	for _, ret := range returns {
		if ret.Source != "" {
			(*allImports)[ret.Source] = true
		}
	}
}

func collectFunctionImports(allImports *map[string]bool, function Function) {
	collectImports(allImports, function.Parameters, function.Returns)
	for _, src := range function.Body.Sources {
		(*allImports)[src] = true
	}
	if function.Receiver != nil && function.Receiver.Type.Source != "" {
		(*allImports)[function.Receiver.Type.Source] = true
	}
}
