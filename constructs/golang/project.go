package golang

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type UnitModule struct {
	Name         string       `yaml:"name"`
	Imports      []string     `yaml:"imports"`
	Structs      []*Struct    `yaml:"structs"`
	Functions    []*Function  `yaml:"functions"`
	Variables    []*Variable  `yaml:"variables"`
	Constants    []*Constant  `yaml:"constants"`
	InitFunction *CodeElement `yaml:"init_fn"`
	MainFunction *CodeElement `yaml:"main"`
}

type Module struct {
	Name  string       `yaml:"name"`
	Units []UnitModule `yaml:"units"`
}

type ProjectRequirement struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type ProjectExclusion struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type ProjectReplacement struct {
	SourceName         string `yaml:"source_name"`
	SourceVersion      string `yaml:"source_version"`
	ReplacementName    string `yaml:"replacement_name"`
	ReplacementVersion string `yaml:"replacement_version"`
}

type ProjectRetract struct {
	Version string `yaml:"version"`
	Comment string `yaml:"comment"`
}

type ProjectReposity struct {
	URL            string `yaml:"url"`
	VersionControl string `yaml:"version_control"`
}

type Project struct {
	Name            string                `yaml:"name"`
	ProjectLocation string                `yaml:"location"`
	ModuleMap       map[string]*Module    `yaml:"modules"`
	GoVersion       string                `yaml:"go_version"`
	Requirements    []*ProjectRequirement `yaml:"requirements"`
	Replacements    []*ProjectReplacement `yaml:"replacements"`
	Exclusions      []*ProjectExclusion   `yaml:"exclusions"`
	Retracts        []*ProjectRetract     `yaml:"retracts"`
	Repositories    []*ProjectReposity    `yaml:"repositories"`
}

var writeFileFun func(string, string) = writeFile

// Helper function to write data to a file, creating directories as needed
func writeFile(path, data string) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Unable to create directory %s: %v", dir, err)
	}

	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		log.Fatalf("Unable to write file %s: %v", path, err)
	}
}

// GenerateGoMod creates the go.mod file based on project configurations
func GenerateGoMod(project Project, basePath string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("module %s\n", project.Name))
	sb.WriteString(fmt.Sprintf("\ngo %s\n", project.GoVersion))

	if len(project.Requirements) > 0 {
		sb.WriteString("\nrequire (\n")
		for _, req := range project.Requirements {
			sb.WriteString(fmt.Sprintf("\t%s %s\n", req.Name, req.Version))
		}
		sb.WriteString(")\n")
	}

	if len(project.Replacements) > 0 {
		sb.WriteString("\nreplace (\n")
		for _, rep := range project.Replacements {
			sb.WriteString(fmt.Sprintf("\t%s %s => %s %s\n", rep.SourceName, rep.SourceVersion, rep.ReplacementName, rep.ReplacementVersion))
		}
		sb.WriteString(")\n")
	}

	if len(project.Exclusions) > 0 {
		sb.WriteString("\nexclude (\n")
		for _, exc := range project.Exclusions {
			sb.WriteString(fmt.Sprintf("\t%s %s\n", exc.Name, exc.Version))
		}
		sb.WriteString(")\n")
	}

	if len(project.Retracts) > 0 {
		sb.WriteString("\nretract (\n")
		for _, ret := range project.Retracts {
			sb.WriteString(fmt.Sprintf("\t%s // %s\n", ret.Version, ret.Comment))
		}
		sb.WriteString(")\n")
	}

	// Write the go.mod file
	writeFileFun(filepath.Join(basePath, "go.mod"), sb.String())
}

func (u *UnitModule) GenerateUnitCode(filepath string, moduleName string) map[Dependency]bool {
	// srcFile := GoSourceFile{
	// 	Package:      moduleName,
	// 	Structs:      u.Structs,
	// 	Functions:    u.Functions,
	// 	Variables:    u.Variables,
	// 	Constants:    u.Constants,
	// 	InitFunction: u.InitFunction,
	// 	MainFunction: u.MainFunction,
	// }

	srcFile := GoSourceFile{
		Package:   moduleName,
		Structs:   u.Structs,
		Functions: u.Functions,
		Variables: u.Variables,
		Constants: u.Constants,
	}

	srcCode, deps, err := srcFile.SourceCode()
	if err != nil {
		log.Fatalf("Unable to generate source code: %v", err)
	}
	writeFile(filepath, srcCode)

	return deps
}

func (m *Module) GenerateModuleCode(modulePath string) (string, error) {
	filepath := filepath.Join(modulePath, m.Name)

	for _, unit := range m.Units {
		unit.GenerateUnitCode(filepath, m.Name)
	}

	return filepath, nil
}

func GenerateProject(project *Project, localPath string) (string, error) {
	// Generate the go.mod file
	GenerateGoMod(*project, localPath)

	for path, module := range project.ModuleMap {
		// Generate the go.mod file
		module.GenerateModuleCode(filepath.Join(localPath, path))
	}

	return "", nil
}
