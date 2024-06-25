package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/google/uuid"
)

func GenerateSQLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var data RequestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate SQL from the received data
	sqlFiles, err := generateSQL(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with generated SQL files (example just returns DDL SQL)
	fmt.Fprintf(w, "DDL SQL:\n%s", sqlFiles.DDL)
}

// generateSQL generates DDL, DML, and Query SQL based on the provided data
func generateSQL(data RequestData) (struct {
	DDL, DML, Query string
}, error) {
	// Placeholder for SQL generation logic
	return struct {
		DDL, DML, Query string
	}{
		DDL:   "CREATE TABLE example (...);",
		DML:   "INSERT INTO example (...);",
		Query: "SELECT * FROM example;",
	}, nil
}

// Function to generate DDL SQL for a given model using text/template
func GenerateDDL(model Model, attributes []Attribute) (string, error) {
	attrMap := make(map[int]Attribute)
	for _, attr := range attributes {
		attrMap[attr.ID] = attr
	}

	id, _ := uuid.NewV7()
	id.ID()
	type TableColumn struct {
		ColumnName string
		Type       string
	}

	var columns []TableColumn
	columns = append(columns, TableColumn{ColumnName: "id", Type: "UUID PRIMARY KEY"})
	for _, attrID := range model.Attributes {
		attr, exists := attrMap[attrID]
		if !exists {
			return "", fmt.Errorf("attribute ID %d not found", attrID)
		}
		columns = append(columns, TableColumn{ColumnName: attr.Name, Type: resolvePostgresType(attr.TypeID)})
	}
	columns = append(columns, TableColumn{ColumnName: "created_at", Type: "TIMESTAMP WITH TIME ZONE"})
	columns = append(columns, TableColumn{ColumnName: "updated_at", Type: "TIMESTAMP WITH TIME ZONE"})
	columns = append(columns, TableColumn{ColumnName: "updated_id", Type: "BIGINT"})
	columns = append(columns, TableColumn{ColumnName: "deleted_at", Type: "TIMESTAMP WITH TIME ZONE"})
	columns = append(columns, TableColumn{ColumnName: "is_deleted", Type: "BOOLEAN"})
	columns = append(columns, TableColumn{ColumnName: "version", Type: "INT"})

	// Setup the template
	const tableTemplate = `CREATE TABLE {{.Namespace}}.{{.Name}} (
    {{- range $index, $col := .Columns}}
    {{if $index}},{{end}}{{$col.ColumnName}} {{$col.Type}}
	{{- end}}
    {{- if .Constraints}}
    ,{{- range $index, $constraint := .Constraints}}{{if $index}},{{end}}UNIQUE ({{join $constraint.Columns ","}})
    {{- end}}
    {{- end}}
);`

	tmpl, err := template.New("table").Funcs(template.FuncMap{
		"join": strings.Join,
	}).Parse(tableTemplate)
	if err != nil {
		return "", err
	}

	var constraints []struct {
		Name    string
		Columns []string
	}

	for _, uc := range model.UniqueConstraints {
		var ucColumns []string
		for _, attrID := range uc.Attributes {
			attr, exists := attrMap[attrID]
			if !exists {
				return "", fmt.Errorf("attribute ID %d not found for unique constraint", attrID)
			}
			ucColumns = append(ucColumns, attr.Name)
		}
		constraints = append(constraints, struct {
			Name    string
			Columns []string
		}{Name: uc.ConstraintName, Columns: ucColumns})
	}

	data := struct {
		Namespace   string
		Name        string
		Columns     []TableColumn
		Constraints []struct {
			Name    string
			Columns []string
		}
	}{
		Namespace:   model.Namespace,
		Name:        model.Name,
		Columns:     columns,
		Constraints: constraints,
	}

	var sql bytes.Buffer
	if err = tmpl.Execute(&sql, data); err != nil {
		return "", err
	}

	return sql.String(), nil
}

// Dummy function to resolve PostgreSQL data types based on TypeID
func resolvePostgresType(typeID int) string {
	// Implement actual logic based on your type mappings
	switch typeID {
	case 1:
		return "VARCHAR(255)"
	case 2:
		return "TEXT"
	case 3:
		return "MONEY"
	default:
		return "VARCHAR(255)"
	}
}
