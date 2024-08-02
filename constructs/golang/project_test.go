package golang

import (
	"path/filepath"
	"reflect"
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
		Structs: []*StructDef{
			{
				Name: "TestStruct",
				Fields: []*Field{
					{Name: "Field1", Type: GoStringType},
					{Name: "Field2", Type: GoIntType},
				},
			},
		},
		Functions: []*FunctionDef{
			{
				Name: "TestFunction",
				Body: CodeElements{
					{FunctionCall: &FunctionCall{Receiver: "fmt", Function: "Println", Args: &Literal{Value: "Hello, world!"}}},
				},
			},
		},
		Variables: []*Variable{
			{Names: "TestVariable", Values: "\"test value\""},
		},
		Constants: []*Constant{
			{Name: "TestConstant", Value: 123},
		},
		InitFunction: CodeElements{
			{FunctionCall: &FunctionCall{Receiver: "fmt", Function: "Println", Args: &Literal{Value: "Initializing module"}}},
		},
		MainFunction: CodeElements{
			{FunctionCall: &FunctionCall{Receiver: "fmt", Function: "Println", Args: &Literal{Value: "Running main function"}}},
		},
	}

	modulePath := "path/to/module"
	deps := um.GenerateAndWriteCode(modulePath, "testModule")
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

	deps2 := um.GenerateAndWriteCode(modulePath, "testModule")
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

func TestModule_GenerateModuleCode(t *testing.T) {
	mockFileData = make(map[string]string) // Initialize the mock file data store
	actualWriteFileFun := writeFile
	writeFileFun = mockWriteFile
	defer func() {
		writeFileFun = actualWriteFileFun
	}()

	// Test case 1: Module with no units or child modules
	m := &Module{
		Name: "testModule",
	}

	modulePath := "path/to/module"
	fPath, dependencies, err := m.GenerateModuleCode(modulePath)
	if err != nil {
		t.Errorf("Error generating module code: %v", err)
	}

	expectedFilepath := filepath.Join(modulePath, m.Name)
	if fPath != expectedFilepath {
		t.Errorf("Expected filepath %s, but got %s", expectedFilepath, fPath)
	}

	if len(dependencies) != 0 {
		t.Errorf("Expected no dependencies, but got %d", len(dependencies))
	}

	// Test case 2: Module with units and child modules
	m = &Module{
		Name: "testModule",

		Dependencies: []Dependency{
			{Source: "github.com/google/uuid", Version: "v1.0.0"},
			{Source: "github.com/stretchr/testify", Version: "v1.7.0"},
		},
		Units: []UnitModule{
			{
				Name: "testUnit1",
				Functions: []*FunctionDef{
					{
						Name: "TestFunction1",
						Body: CodeElements{
							{FunctionCall: &FunctionCall{Receiver: "fmt", Function: "Println", Args: &Literal{Value: "Hello, world!"}}},
						},
					},
				},
				Imports: []string{"fmt"},
			},
			{
				Name: "testUnit2",
				Functions: []*FunctionDef{
					{
						Name: "TestFunction2",
						Body: CodeElements{
							{FunctionCall: &FunctionCall{Receiver: "fmt", Function: "Println", Args: &Literal{Value: "Hello, world 2.0!"}}},
						}},
				},
				Imports: []string{"fmt"},
			},
		},
		ChildModules: []*Module{
			{
				Name: "testChildModule1",
				Units: []UnitModule{
					{
						Name: "testChildUnit1",
						Functions: []*FunctionDef{
							{
								Name: "TestChildFunction1",
								Body: CodeElements{
									{FunctionCall: &FunctionCall{Receiver: "fmt", Function: "Println", Args: &Literal{Value: "Hello, world child!"}}},
								},
							},
						},
						Imports: []string{"fmt"},
					},
				},
				Dependencies: []Dependency{
					{Source: "yourcompany.com/yourproject", Version: "v2.1.0"},
					{Source: "yourcompany.com/yourproject2", Version: "v2.2.0"},
				},
			},
			{
				Name: "testChildModule2",
				Units: []UnitModule{
					{
						Name: "testChildUnit2",
						Functions: []*FunctionDef{
							{
								Name: "TestChildFunction2",
								Body: CodeElements{
									{FunctionCall: &FunctionCall{Receiver: "fmt", Function: "Println", Args: &Literal{Value: "Hello, world child 2.0!"}}},
								}},
						},
						Imports: []string{"fmt"},
					},
				},
			},
		},
	}

	modulePath = "path/to/module"
	fPath, dependencies, err = m.GenerateModuleCode(modulePath)
	if err != nil {
		t.Errorf("Error generating module code: %v", err)
	}

	expectedFilepath = filepath.Join(modulePath, m.Name)
	if fPath != expectedFilepath {
		t.Errorf("Expected filepath %s, but got %s", expectedFilepath, fPath)
	}

	expectedDependencies := map[Dependency]bool{
		{
			Source:  "github.com/google/uuid",
			Version: "v1.0.0",
		}: true,
		{
			Source:  "github.com/stretchr/testify",
			Version: "v1.7.0",
		}: true,
		{
			Source:  "yourcompany.com/yourproject",
			Version: "v2.1.0",
		}: true,
		{
			Source:  "yourcompany.com/yourproject2",
			Version: "v2.2.0",
		}: true,
	}
	if !reflect.DeepEqual(dependencies, expectedDependencies) {
		t.Errorf("Expected dependencies %v, but got %v", expectedDependencies, dependencies)
	}

	testChildModule1Code := `package testChildModule1

import (
	"fmt"
)

func TestChildFunction1() {
	fmt.Println("Hello, world child!")
}
`
	testChildModule2Code := `package testChildModule2

import (
	"fmt"
)

func TestChildFunction2() {
	fmt.Println("Hello, world child 2.0!")
}
`

	testModuleCode := `package testModule

import (
	"fmt"
)

func TestFunction1() {
	fmt.Println("Hello, world!")
}
`

	testModule2Code := `package testModule

import (
	"fmt"
)

func TestFunction2() {
	fmt.Println("Hello, world 2.0!")
}
`

	// Verify the generated code for the module
	expectedCodeMap := map[string]string{
		"path/to/module/testModule/testChildModule1/testChildUnit1.go": testChildModule1Code,
		"path/to/module/testModule/testChildModule2/testChildUnit2.go": testChildModule2Code,
		"path/to/module/testModule/testUnit1.go":                       testModuleCode,
		"path/to/module/testModule/testUnit2.go":                       testModule2Code,
	}

	// t.Log(mockFileData)

	for p, expectedCode := range expectedCodeMap {
		if mockFileData[p] != expectedCode {
			t.Errorf("At path %s, Expected code:\n%s\nbut got:\n%s", p, expectedCode, mockFileData[p])
		}
	}

}

func TestGenerateProject(t *testing.T) {
	mockFileData = make(map[string]string) // Initialize the mock file data store
	actualWriteFileFun := writeFile
	writeFileFun = mockWriteFile
	defer func() {
		writeFileFun = actualWriteFileFun
	}()

	// Test case 1: Generate project with no modules
	project := &Project{
		Name:      "testProject",
		GoVersion: "1.22.0",
		Modules:   nil,
	}

	projectPath := "path/to/project"
	err := project.GenerateProject(projectPath)
	if err != nil {
		t.Errorf("Error generating project: %v", err)
	}

	// Test case 2: Generate project with modules
	project = &Project{
		Name:      "testProject",
		GoVersion: "1.22.0",
		Modules: []*Module{
			{
				Name: "testModule1",
				Units: []UnitModule{
					{
						Name: "testUnit1",
						Functions: []*FunctionDef{
							{
								Name: "TestFunction1",
								Body: CodeElements{
									{FunctionCall: &FunctionCall{Receiver: "fmt", Function: "Println", Args: &Literal{Value: "Hello, world!"}}},
								},
							},
						},
						Imports: []string{"fmt"},
					},
				},
			},
			{
				Name: "testModule2",
				Units: []UnitModule{
					{
						Name: "testUnit2",
						Functions: []*FunctionDef{
							{
								Name: "TestFunction2",
								Body: CodeElements{
									{FunctionCall: &FunctionCall{Receiver: "fmt", Function: "Println", Args: &Literal{Value: "Hello, world 2.0!"}}},
								},
							},
						},
						Imports: []string{"fmt"},
					},
				},
			},
		},
	}

	goModCode := `module testProject

go 1.22.0
`
	testUnit1Code := `package testModule1

import (
	"fmt"
)

func TestFunction1() {
	fmt.Println("Hello, world!")
}
`
	testUnit2Code := `package testModule2

import (
	"fmt"
)

func TestFunction2() {
	fmt.Println("Hello, world 2.0!")
}
`
	expectedCodeMap := map[string]string{
		"path/to/project/go.mod":                   goModCode,
		"path/to/project/testModule1/testUnit1.go": testUnit1Code,
		"path/to/project/testModule2/testUnit2.go": testUnit2Code,
	}
	projectPath = "path/to/project"
	err = project.GenerateProject(projectPath)
	if err != nil {
		t.Errorf("Error generating project: %v", err)
	}

	// t.Log(mockFileData)

	for p, expectedCode := range expectedCodeMap {
		if mockFileData[p] != expectedCode {
			t.Errorf("At path %s, Expected code:\n%s\nbut got:\n%s", p, expectedCode, mockFileData[p])
		}
	}
}
