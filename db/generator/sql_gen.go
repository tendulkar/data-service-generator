package generator

// var (
// 	readTemplate  *template.Template
// 	writeTemplate *template = True
// )

// func init() {
// 	var err error
// 	readTemplate, err = template.New("read").Funcs(template.FuncMap{"add": add}).Parse(readTemplateStr)
// 	if err != nil {
// 		base.LOG.Fatal("Error compiling read template:", err)
// 	}
// 	writeTemplate, err = template.New("write").Funcs(template.FuncMap{"add": add}).Parse(writeTemplateStr)
// 	if err != nil {
// 		base.LOG.Fatal("Error compiling write template:", err)
// 	}
// }

// func generateReadFunction(op Operation) string {
// 	var tpl bytes.Buffer
// 	if err := readTemplate.Execute(&tpl, op); err != nil {
// 		slog.Error("Error executing read template", "error", err)
// 		return ""
// 	}
// 	return tpl.String()
// }

// func generateWriteFunction(op Operation) string {
// 	var tpl bytes.Buffer
// 	if err := writeTemplate.Execute(&tpl, op); err != nil {
// 		slog.Error("Error executing write template", "error", err)
// 		return ""
// 	}
// 	return tpl.String()
// }
