package golang

import (
	"fmt"
	"log"
	"os"
	"path"
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
	InitFunction CodeElements `yaml:"init_fn"`
	MainFunction CodeElements `yaml:"main"`
}

type Module struct {
	Name         string       `yaml:"name"`
	Units        []UnitModule `yaml:"units"`
	ChildModules []*Module    `yaml:"child_modules"`
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
func GenerateGoMod(project Project, cleanDeps map[string]string, basePath string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("module %s\n", project.Name))
	sb.WriteString(fmt.Sprintf("\ngo %s\n", project.GoVersion))

	allDeps := map[string]string{}
	for dep, version := range cleanDeps {
		allDeps[dep] = version
	}

	if len(project.Requirements) > 0 {
		sb.WriteString("\nrequire (\n")
		for _, req := range project.Requirements {
			allDeps[req.Name] = req.Version
		}
		for dep, version := range allDeps {
			sb.WriteString(fmt.Sprintf("\t%s %s\n", dep, version))
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

	srcFile := GoSourceFile{
		Package:      moduleName,
		Imports:      u.Imports,
		Structs:      u.Structs,
		Functions:    u.Functions,
		Variables:    u.Variables,
		Constants:    u.Constants,
		InitFunction: u.InitFunction,
		MainFunction: u.MainFunction,
	}

	srcCode, deps, err := srcFile.SourceCode()
	if err != nil {
		log.Fatalf("Unable to generate source code: %v", err)
	}

	goSrcPath := path.Join(filepath, u.Name+".go")
	writeFileFun(goSrcPath, srcCode)

	return deps
}

func (m *Module) GenerateModuleCode(modulePath string) (string, map[Dependency]bool, error) {
	filepath := filepath.Join(modulePath, m.Name)
	dependencies := make(map[Dependency]bool)
	for _, unit := range m.Units {
		deps := unit.GenerateUnitCode(filepath, m.Name)
		for dep := range deps {
			dependencies[dep] = true
		}
	}

	for _, child := range m.ChildModules {
		_, childDeps, err := child.GenerateModuleCode(filepath)
		if err != nil {
			return "", nil, err
		}

		for dep := range childDeps {
			dependencies[dep] = true
		}
	}

	return filepath, dependencies, nil
}

func GenerateProject(project *Project, localPath string) (string, error) {
	dependencies := make(map[Dependency]bool)

	for path, module := range project.ModuleMap {
		// Generate the go.mod file
		_, deps, err := module.GenerateModuleCode(filepath.Join(localPath, path))
		if err != nil {
			return "", err
		}

		for dep := range deps {
			dependencies[dep] = true
		}
	}

	cleanDeps := make(map[string]string)

	for dep := range dependencies {
		if _, ok := cleanDeps[dep.Source]; ok {
			return "", fmt.Errorf("duplicate dependency: %v with version %s", dep, cleanDeps[dep.Source])
		}
		cleanDeps[dep.Source] = dep.Version
	}

	// Generate the go.mod file
	GenerateGoMod(*project, cleanDeps, localPath)

	return "", nil
}
