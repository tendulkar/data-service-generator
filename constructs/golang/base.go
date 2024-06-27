package golang

// import (
// 	"fmt"
// 	"strings"
// )

// type GoConstruct interface {
// 	Code() string
// 	Imports() []string
// }

// type GoConstructs []GoConstruct

// type GoType struct {
// 	TypeName string
// 	Source   string
// }

// func (t *GoType) Code() string {
// 	return t.TypeName
// }

// func (t *GoType) Imports() []string {
// 	if t.Source == "" {
// 		return []string{}
// 	}
// 	return []string{t.Source}
// }

// var (
// 	GoIntType    = GoType{TypeName: "int", Source: ""}
// 	GoInt64Type  = GoType{TypeName: "int64", Source: ""}
// 	GoInt32Type  = GoType{TypeName: "int32", Source: ""}
// 	GoInt16Type  = GoType{TypeName: "int16", Source: ""}
// 	GoInt8Type   = GoType{TypeName: "int8", Source: ""}
// 	GoUIntType   = GoType{TypeName: "uint", Source: ""}
// 	GoUInt64Type = GoType{TypeName: "uint64", Source: ""}
// 	GoUInt32Type = GoType{TypeName: "uint32", Source: ""}
// 	GoUInt16Type = GoType{TypeName: "uint16", Source: ""}
// 	GoUInt8Type  = GoType{TypeName: "uint8", Source: ""}
// 	GoByteType   = GoType{TypeName: "byte", Source: ""}
// 	GoStringType = GoType{TypeName: "string", Source: ""}
// 	GoFloatType  = GoType{TypeName: "float64", Source: ""}
// 	GoBoolType   = GoType{TypeName: "bool", Source: ""}
// 	GoUUIDType   = GoType{TypeName: "uuid.UUID", Source: "github.com/gofrs/uuid"}
// )

// type GoCodeBlock struct {
// 	CodeBlock string
// 	Sources   []string
// }

// func (c *GoCodeBlock) Code() string {
// 	return c.CodeBlock
// }

// func (c *GoCodeBlock) Imports() []string {
// 	imports := make(map[string]bool)
// 	for _, source := range c.Sources {
// 		imports[source] = true
// 	}
// 	importsList := make([]string, 0, len(imports))
// 	for imp := range imports {
// 		importsList = append(importsList, imp)
// 	}
// 	return importsList
// }

// type GoParam struct {
// 	Name string
// 	Type GoType
// }

// type GoParams []GoParam

// func (p *GoParams) Code() string {
// 	var code []string
// 	for _, param := range *p {
// 		code = append(code, fmt.Sprintf("%s %s", param.Name, param.Type.Code()))
// 	}
// 	return strings.Join(code, ", ")
// }

// func (p *GoParams) Imports() []string {
// 	imports := make(map[string]bool)
// 	for _, param := range *p {
// 		imports[param.Type.Source] = true
// 	}
// 	importsList := make([]string, 0, len(imports))
// 	for imp := range imports {
// 		importsList = append(importsList, imp)
// 	}
// 	return importsList
// }

// type GoFunction struct {
// 	Name       string
// 	Params     GoParams
// 	ReturnType GoType
// 	Body       GoCodeBlock
// }

// func (f *GoFunction) Code() string {
// 	code := IndentCode(f.Body.Code(), 1)
// 	returnTypeCode := f.ReturnType.Code()
// 	returnTypeSpace := ""
// 	if returnTypeCode != "" {
// 		returnTypeSpace = " "
// 	}
// 	return fmt.Sprintf("func %s(%s) %s%s{\n%s}", f.Name, f.Params.Code(), f.ReturnType.Code(), returnTypeSpace, code)
// }

// func (f *GoFunction) Imports() []string {
// 	allImports := f.ReturnType.Imports()
// 	for _, param := range f.Params {
// 		allImports = append(allImports, param.Type.Imports()...)
// 	}
// 	allImports = append(allImports, f.Body.Imports()...)
// 	return allImports
// }

// type GoFunctions []GoFunction

// func (f *GoFunctions) Code() string {
// 	var code []string
// 	for _, function := range *f {
// 		code = append(code, function.Code())
// 	}
// 	return strings.Join(code, "\n\n")
// }

// func (f *GoFunctions) Imports() []string {
// 	imports := make(map[string]bool)
// 	for _, function := range *f {
// 		for _, imp := range function.Imports() {
// 			imports[imp] = true
// 		}
// 	}
// 	importsList := make([]string, 0, len(imports))
// 	for imp := range imports {
// 		importsList = append(importsList, imp)
// 	}
// 	return importsList
// }

// type GoField struct {
// 	Name string
// 	Type GoType
// }

// type GoFields []GoField

// func (f *GoFields) Code() string {
// 	var code []string
// 	for _, field := range *f {
// 		code = append(code, fmt.Sprintf("%s %s", field.Name, field.Type.Code()))
// 	}
// 	return strings.Join(code, "\n")
// }

// func (f *GoFields) Imports() []string {
// 	imports := make(map[string]bool)
// 	for _, field := range *f {
// 		if field.Type.Source == "" {
// 			continue
// 		}
// 		imports[field.Type.Source] = true
// 	}
// 	importsList := make([]string, 0, len(imports))
// 	for imp := range imports {
// 		importsList = append(importsList, imp)
// 	}
// 	return importsList
// }

// type GoStruct struct {
// 	Name   string
// 	Fields GoFields
// }

// func (s *GoStruct) Code() string {
// 	fieldsCode := s.Fields.Code()
// 	indentedFieldsCode := IndentCode(fieldsCode, 1)
// 	return fmt.Sprintf("type %s struct {\n%s\n}", s.Name, indentedFieldsCode)
// }

// func (s *GoStruct) Imports() []string {
// 	allImports := s.Fields.Imports()
// 	return allImports
// }

// type GoMemberFunction struct {
// 	*GoFunction
// 	Receiver GoType
// }

// func (m *GoMemberFunction) Code() string {
// 	code := IndentCode(m.Body.Code(), 1)
// 	returnTypeCode := m.ReturnType.Code()
// 	returnTypeSpace := ""
// 	if returnTypeCode != "" {
// 		returnTypeSpace = " "
// 	}
// 	receiver := fmt.Sprintf("(m *%s)", m.Receiver.TypeName)
// 	return fmt.Sprintf("func %s(%s) %s%s{\n%s}", m.Name, m.Params.Code(), m.ReturnType.Code(), returnTypeSpace, code)
// }

// type GoStructMembers struct {
// 	*GoStruct
// 	Members GoFunctions
// }

// func (s *GoStructMembers) Code() string {
// 	membersCode := s.Members.Code()
// 	indentedMembersCode := IndentCode(membersCode, 1)
// 	return fmt.Sprintf("%s\n\n%s", s.GoStruct.Code(), indentedMembersCode)
// }

// func (s *GoStructMembers) Imports() []string {
// 	allImports := s.GoStruct.Imports()
// 	for _, member := range s.Members {
// 		allImports = append(allImports, member.Imports()...)
// 	}
// 	return allImports
// }
