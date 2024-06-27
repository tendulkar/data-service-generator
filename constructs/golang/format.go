package golang

import (
	"fmt"
	"strings"
)

const (
	Indent = "  "
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
