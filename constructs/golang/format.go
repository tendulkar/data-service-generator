package golang

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

const (
	Indent = "\t"
)

func IndentCode(code string, indentSize int) string {
	var lines []string
	indent := strings.Repeat(Indent, indentSize)
	for _, line := range strings.Split(code, "\n") {
		trimmedLine := strings.Trim(line, " ")
		if trimmedLine == "" {
			lines = append(lines, trimmedLine)
			continue
		}
		lines = append(lines, fmt.Sprintf("%s%s", indent, trimmedLine))
	}
	return strings.Join(lines, "\n")
}

// go:inline
func ToCamelCase(input string) string {
	return strcase.ToLowerCamel(input)
}

// go:inline
func ToCamelCaseArray(input []string) []string {
	var result []string
	for _, item := range input {
		result = append(result, ToCamelCase(item))
	}
	return result
}

// Converts a string to pascal case
func ToPascalCase(input string) string {
	return strcase.ToCamel(input)
}

func ToPascalCaseArray(input []string) []string {
	var result []string
	for _, item := range input {
		result = append(result, ToPascalCase(item))
	}
	return result
}

// Converts a string to snake case
func ToSnakeCase(input string) string {
	return strcase.ToSnake(input)
}

func ToSnakeCaseArray(input []string) []string {
	var result []string
	for _, item := range input {
		result = append(result, ToSnakeCase(item))
	}
	return result
}
