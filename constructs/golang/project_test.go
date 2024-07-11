package golang

import (
	"path/filepath"
	"testing"
)

var mockFileData map[string]string

func mockWriteFile(path, data string) {
	mockFileData[path] = data
}

func TestGenerateGoMod(t *testing.T) {
	basePath := "/tmp/project"
	mockFileData = make(map[string]string) // Initialize the mock file data store

	actualWriteFileFun := writeFile
	writeFileFun = mockWriteFile
	defer func() {
		writeFileFun = actualWriteFileFun
	}()

	cleanDeps := make(map[string]string)
	cleanDeps["github.com/stretchr/testify"] = "v1.0.2"
	project := Project{
		Name:         "example.com/testproject",
		GoVersion:    "1.15",
		Requirements: []*ProjectRequirement{{Name: "github.com/example/library", Version: "v1.0.0"}},
		Replacements: []*ProjectReplacement{{SourceName: "github.com/old/library", SourceVersion: "v1.0.0",
			ReplacementName: "github.com/new/library", ReplacementVersion: "v1.1.0"}},
		Exclusions: []*ProjectExclusion{{Name: "github.com/unused/library", Version: "v1.0.0"}},
		Retracts:   []*ProjectRetract{{Version: "v1.2.3", Comment: "Critical bug"}},
	}

	GenerateGoMod(project, cleanDeps, basePath)

	expectedGoMod := `module example.com/testproject

go 1.15

require (
	github.com/stretchr/testify v1.0.2
	github.com/example/library v1.0.0
)

replace (
	github.com/old/library v1.0.0 => github.com/new/library v1.1.0
)

exclude (
	github.com/unused/library v1.0.0
)

retract (
	v1.2.3 // Critical bug
)
`
	goModPath := filepath.Join(basePath, "go.mod")
	if got, exists := mockFileData[goModPath]; !exists || got != expectedGoMod {
		t.Errorf("GenerateGoMod() = %v, want %v", got, expectedGoMod)
	}
}

func TestUnitModule_GenerateCode(t *testing.T) {
	mockFileData = make(map[string]string) // Initialize the mock file data store
	actualWriteFileFun := writeFile
	writeFileFun = mockWriteFile
	defer func() {
		writeFileFun = actualWriteFileFun
	}()

	// Test case 1: Generate code for unit module with all details
	um := &UnitModule{
		Name:    "testModule",
		Imports: []string{"fmt", "strings"},
		Structs: []*Struct{
			{
				Name: "TestStruct",
				Fields: []*Field{
					{Name: "Field1", Type: GoStringType},
					{Name: "Field2", Type: GoIntType},
				},
			},
		},
		Functions: []*Function{
			{
				Name: "TestFunction",
				Body: CodeElements{
					{MemberFunctionCall: &MemberFunctionCall{Receiver: "fmt", Function: "Println", Params: &Literal{Value: "Hello, world!"}}},
				},
			},
		},
		Variables: []*Variable{
			{Name: "TestVariable", Value: "\"test value\""},
		},
		Constants: []*Constant{
			{Name: "TestConstant", Value: "123"},
		},
		InitFunction: CodeElements{
			{MemberFunctionCall: &MemberFunctionCall{Receiver: "fmt", Function: "Println", Params: &Literal{Value: "Initializing module"}}},
		},
		MainFunction: CodeElements{
			{MemberFunctionCall: &MemberFunctionCall{Receiver: "fmt", Function: "Println", Params: &Literal{Value: "Running main function"}}},
		},
	}

	modulePath := "path/to/module"
	deps := um.GenerateUnitCode(modulePath, "testModule")
	if len(deps) != 0 {
		t.Errorf("Error generating code, should have empty deps, but got: %v", deps)
	} // Verify the generated code
	expectedCode := `package testModule

import (
	"fmt"
	"strings"
)

var TestVariable = "test value"

const TestConstant = 123

type TestStruct struct {
	Field1 string
	Field2 int
}

func TestFunction() {
	fmt.Println("Hello, world!")
}

func init() {
	fmt.Println("Initializing module")
}

func main() {
	fmt.Println("Running main function")
}
`
	goSrcPath := filepath.Join(modulePath, "testModule.go")
	if got, exists := mockFileData[goSrcPath]; !exists || got != expectedCode {
		t.Errorf("GenerateGoMod() = %v, want %v", got, expectedCode)
	}

	// Test case 2: Generate code for unit module with no details
	um = &UnitModule{
		Name: "testModule",
	}

	deps2 := um.GenerateUnitCode(modulePath, "testModule")
	if len(deps2) != 0 {
		t.Errorf("Error generating code, should have empty deps, but got: %v", deps2)
	}

	// Verify the generated code
	expectedCode = `package testModule
`
	goSrcPath2 := filepath.Join(modulePath, "testModule.go")
	if got, exists := mockFileData[goSrcPath2]; !exists || got != expectedCode {
		t.Errorf("GenerateGoMod() = %v, want %v", got, expectedCode)
	}
}
