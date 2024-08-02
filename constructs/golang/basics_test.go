package golang

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestTranslateToGoType(t *testing.T) {
	// Define test cases with input type names and expected GoTypes and errors
	testCases := []struct {
		typeName          string
		expectedGoType    *GoType
		expectedErrorText string
	}{
		{
			"int16",
			GoInt16Type,
			"",
		},
		{
			"float32",
			GoFloat32Type,
			"",
		},
		{
			"[]complex64",
			GoComplex64ArrayType,
			"",
		},
		{
			"map[interface{}]interface{}",
			GoInterfaceInterfaceMapType,
			"",
		},
		{
			"unknown",
			&GoType{},
			"type unknown not found",
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		actualGoType, err := TranslateToGoType(tc.typeName)

		// Check if the actual GoType matches the expected GoType
		if !reflect.DeepEqual(actualGoType, tc.expectedGoType) {
			t.Errorf("For type %s: expected GoType %v, got %v", tc.typeName, tc.expectedGoType, actualGoType)
		}

		// Check if the actual error message matches the expected error message
		if err != nil && err.Error() != tc.expectedErrorText {
			t.Errorf("For type %s: expected error message %s, got %s", tc.typeName, tc.expectedErrorText, err.Error())
		}
	}
}

func TestFunctionCodeGeneration(t *testing.T) {
	f := FunctionDef{
		Name: "Update",
		Parameters: []*Parameter{
			{Name: "data", Type: &GoType{Name: "*Data", Source: "github.com/example/data"}},
		},
		Returns: []*Parameter{{Type: &GoType{Name: "error", Source: ""}}},
		Body: CodeElements{{
			Imports: []string{"github.com/example/data"},
			FunctionCall: &FunctionCall{
				Receiver: "this.Data",
				Function: "Merge",
				Args:     "data",
			},
		},
			{
				Return: "nil",
			},
		},
		Receiver: &Receiver{Name: "this", Type: &GoType{Name: "Processor", Source: ""}},
	}

	expectedImports := fmt.Sprintf("import (\n%s\"github.com/example/data\"\n)", Indent)
	expectedSignature := "func (this *Processor) Update(data *Data) error {"
	expectedBody := fmt.Sprintf("%sthis.Data.Merge(data)\n%sreturn nil", Indent, Indent)
	result, resultImports := f.FunctionCode()
	resultImportCode := generateImports(resultImports)

	t.Log("Testing FunctionCodeGeneration, generated code is: \n", result)
	// fmt.Println("Testing FunctionCodeGeneration, generated code is: \n", result)

	if !strings.Contains(resultImportCode, expectedImports) {
		t.Errorf("Expected to find import statement: %s, got %s", expectedImports, resultImportCode)
	}
	if !strings.Contains(result, expectedSignature) {
		t.Errorf("Expected function signature: %s, got %s", expectedSignature, result)
	}
	if !strings.Contains(result, expectedBody) {
		t.Errorf("Expected function body to contain: %s, got %s", expectedBody, result)
	}
}

func TestStructCodeGeneration(t *testing.T) {

	s := StructDef{
		Name: "Processor",
		Fields: []*Field{
			{Name: "Data", Type: &GoType{Name: "Data", Source: "github.com/example/data"}},
			{Name: "Logger", Type: &GoType{Name: "Logger", Source: "github.com/sirupsen/logrus"}},
		},
		Functions: []*FunctionDef{
			{
				Name: "Process",
				Parameters: []*Parameter{
					{Name: "input", Type: &GoType{Name: "[]byte", Source: ""}},
				},
				Returns: []*Parameter{{Type: &GoType{Name: "error", Source: ""}}},
				Body: []*CodeElement{
					{
						FunctionCall: &FunctionCall{
							Receiver: "this.Logger",
							Function: "Info",
							Args: []*Literal{
								{Value: "Processing"},
							},
						},
						Imports: []string{"github.com/sirupsen/logrus"},
					},
					{
						Return: "nil",
					},
				},
				Receiver: &Receiver{Name: "this", Type: &GoType{Name: "Processor", Source: ""}},
			},
		},
	}

	// expectedImports := "import (\n  \"github.com/example/data\"\n  \"github.com/sirupsen/logrus\"\n)"
	expectedImports := map[string]bool{"github.com/example/data": true, "github.com/sirupsen/logrus": true}
	expectedStruct := fmt.Sprintf("type Processor struct {\n%sData *Data\n%sLogger *Logger\n}", Indent, Indent)
	expectedMethod := "func (this *Processor) Process(input []byte) error {"
	result, imports := s.StructCode()

	t.Log("Testing StructCodeGeneration, generated code is: \n", result)
	// fmt.Println("Testing StructCodeGeneration, generated code is: \n", result, "\n", imports)
	if !reflect.DeepEqual(expectedImports, imports) {
		t.Errorf("Expected to find import statements: %v, got %v", expectedImports, imports)
	}
	if !strings.Contains(result, expectedStruct) {
		t.Errorf("Expected struct definition to contain: %s, got %s", expectedStruct, result)
	}
	if !strings.Contains(result, expectedMethod) {
		t.Errorf("Expected method definition to contain: %s, got %s", expectedMethod, result)
	}
}
func TestGenerateGoFile(t *testing.T) {
	// Define structs, functions, variables, constants, and init function
	structs := []*StructDef{
		{
			Name: "Logger",
			Fields: []*Field{
				{Name: "Level", Type: &GoType{Name: "int", Source: ""}},
			},
			Functions: []*FunctionDef{
				{
					Name:       "SetLevel",
					Parameters: []*Parameter{{Name: "level", Type: &GoType{Name: "int", Source: ""}}},
					Returns:    []*Parameter{},
					Body: []*CodeElement{
						{
							Assign: &Assignment{
								Left:  "l.Level",
								Right: "level",
							},
						},
					},
					Receiver: &Receiver{Name: "l", Type: &GoType{Name: "Logger", Source: ""}},
				},
			},
		},
	}
	functions := []*FunctionDef{
		{
			Name:       "NewLogger",
			Parameters: []*Parameter{},
			Returns:    []*Parameter{{Type: &GoType{Name: "*Logger", Source: ""}}},
			Body: []*CodeElement{
				{
					Return: &CodeElement{
						StructCreation: &StructCreation{StructType: "Logger", KeyValues: KeyValues{{Key: "Level", Value: "0"}}},
					},
				},
			},
		},
	}
	variables := []*Variable{
		{Names: "defaultLogger", Type: "Logger", IsReference: true, Values: "NewLogger()"},
	}
	constants := []*Constant{
		{Name: "DefaultLevel", Type: "int", Value: 1},
	}
	initFunction := CodeElements{
		{
			FunctionCall: &FunctionCall{
				Receiver: "defaultLogger",
				Function: "SetLevel",
				Args:     "DefaultLevel",
			},
		},
	}
	generatedCode, _, err := GenerateGoFile("main", structs, functions, variables, constants, initFunction, nil, nil, nil)

	if err != nil {
		t.Error(err)
	}

	t.Log("Testing GenerateGoFile, generated code is: \n", generatedCode)
	// formattedCode, err := format.Source([]byte(generatedCode))
	// fmt.Printf("Testing GenerateGoFile, generated code is: \n%s\nFormatted code is: \n%s\nError: %v\n", generatedCode, string(formattedCode), err)
	// Check package declaration
	if !strings.Contains(generatedCode, "package main") {
		t.Error("Package declaration is missing or incorrect")
	}

	// Check for import statements (if any)
	if strings.Contains(generatedCode, "import (") {
		t.Error("Unexpected imports detected")
	}

	// Check for struct declaration
	if !strings.Contains(generatedCode, "type Logger struct {") || !strings.Contains(generatedCode, "Level int") {
		t.Error("Struct Logger declaration is missing or incorrect")
	}

	// Check for member function
	if !strings.Contains(generatedCode, "func (l *Logger) SetLevel(level int) {") {
		t.Error("Member function SetLevel is missing or incorrect")
	}

	// Check for standalone function
	if !strings.Contains(generatedCode, "func NewLogger() *Logger {") {
		t.Error("Standalone function NewLogger is missing or incorrect")
	}

	// Check for variable declaration
	if !strings.Contains(generatedCode, "var defaultLogger *Logger = NewLogger()") {
		t.Error("Variable defaultLogger declaration is missing or incorrect")
	}

	// Check for constant declaration
	if !strings.Contains(generatedCode, "const DefaultLevel int = 1") {
		t.Error("Constant DefaultLevel declaration is missing or incorrect")
	}

	// Check for init function
	if !strings.Contains(generatedCode, "func init() {") || !strings.Contains(generatedCode, "defaultLogger.SetLevel(DefaultLevel)") {
		t.Error("Init function is missing or incorrect")
	}
}
