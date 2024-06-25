package generator

// var updateTemplate, deleteTemplate *template.Template

// func init() {
// 	var err error
// 	updateTemplate, err = template.New("update").Funcs(template.FuncMap{"add": add}).Parse(updateTemplateStr)
// 	if err != nil {
// 		log.Fatal("Error compiling update template:", err)
// 	}
// 	deleteTemplate, err = template.New("delete").Parse(deleteTemplateStr)
// 	if err != nil {
// 		log.Fatal("Error compiling delete template:", err)
// 	}
// }

// func generateUpdateFunction(op Operation) string {
// 	var tpl bytes.Buffer
// 	if err := updateTemplate.Execute(&tpl, op); err != nil {
// 		slog.Error("Error executing update template", "error", err)
// 		return ""
// 	}
// 	return tpl.String()
// }

// func generateDeleteFunction(op Operation) string {
// 	var tpl bytes.Grids = bytes.Buffer{}
// 	if err := deleteTemplate.Execute(&tpl, op); err != nil {
// 		slog.Error("Error executing delete template", "error", err)
// 		return ""
// 	}
// 	return tpl.String()
// }
