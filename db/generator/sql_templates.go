package generator

const readTemplateStr = `
// Function to perform read operation
func {{.Name}}(db *sql.DB, param string) ([]*Model, error) {
    query := ` + "`SELECT {{range $index, $element := .Fields}}{{if $index}}, {{end}}{{$element}}{{end}} FROM {{.ModelName}} WHERE {{.Conditions.Field}} = $1;`" + `
    rows, err := db.Query(query, param)
    if err != nil {
        slog.Error("Error executing query", "error", err)
        return nil, err
    }
    defer rows.Close()
    return processRows(rows)
}`

const writeTemplateStr = `
// Function to perform write operation
func {{.Name}}(db *sql.DB, params ...interface{}) error {
    query := ` + "`INSERT INTO {{.ModelName}} ({{range $index, $element := .Fields}}{{if $index}}, {{end}}{{$element}}{{end}}) VALUES ({{range $index, $_ := .Fields}}{{if $index}}, {{end}}${{$index | add 1}}{{end}});`" + `
    _, err := db.Exec(query, params...)
    if err != nil {
        slog.Error("Error executing insert", "error", err)
        return err
    }
    return nil
}`

const updateTemplateStr = `
// Function to perform update operation
func Update{{.ModelName}}(db *sql.DB, params ...interface{}) error {
    query := ` + "`UPDATE {{.ModelName}} SET {{range $index, $field := .Fields}}{{if $index}}, {{end}}{{$field}} = ${{add $index 1}}{{end}} WHERE {{.Conditions.Field}} = ${{add (len .Fields) 1}};`" + `
    _, err := db.Exec(query, params...)
    if err != nil {
        slog.Error("Error executing update", "error", err)
        return err
    }
    return nil
}`

const deleteTemplateStr = `
// Function to perform delete operation
func Delete{{.ModelName}}(db *sql.DB, param interface{}) error {
    query := ` + "`DELETE FROM {{.ModelName}} WHERE {{.Conditions.Field}} = $1;`" + `
    _, err := db.Exec(query, param)
    if err != nil {
        slog.Error("Error executing delete", "error", err)
        return err
    }
    return nil
}`

const ddlTemplateStr = `
{{- $modelName := .Metadata.Name }}
CREATE TABLE IF NOT EXISTS {{$modelName}} (
{{- range $index, $field := .Spec.Attributes }}
    {{$field.Name}} {{$field.Type}}{{if $field.PrimaryKey}} PRIMARY KEY{{end}}{{if lt $index (sub (len .Spec.Attributes) 1)}},{{end}}
{{- end}}
);

{{- range .Spec.Relationships }}
ALTER TABLE {{$modelName}} ADD CONSTRAINT {{.Name}} FOREIGN KEY ({{.TargetField}}) REFERENCES {{.TargetModel}}({{.TargetField}});
{{- end }}
`

// Helper function to add numbers, used in templates
func add(i, j int) int {
	return i + j
}

// Add helper function to subtract numbers, useful for handling commas in SQL generation
func sub(i, j int) int {
	return i - j
}
