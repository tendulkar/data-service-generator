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

	project := Project{
		Name:         "example.com/testproject",
		GoVersion:    "1.15",
		Requirements: []*ProjectRequirement{{Name: "github.com/example/library", Version: "v1.0.0"}},
		Replacements: []*ProjectReplacement{{SourceName: "github.com/old/library", SourceVersion: "v1.0.0",
			ReplacementName: "github.com/new/library", ReplacementVersion: "v1.1.0"}},
		Exclusions: []*ProjectExclusion{{Name: "github.com/unused/library", Version: "v1.0.0"}},
		Retracts:   []*ProjectRetract{{Version: "v1.2.3", Comment: "Critical bug"}},
	}

	GenerateGoMod(project, basePath)

	expectedGoMod := `module example.com/testproject

go 1.15

require (
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
