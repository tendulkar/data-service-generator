package golang

// import (
// 	"reflect"
// 	"testing"
// )

// func TestGoFunction_Code(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		function GoFunction
// 		want     string
// 	}{
// 		{
// 			name: "no params",
// 			function: GoFunction{
// 				Name:   "testFunc",
// 				Params: GoParams{},
// 			},
// 			want: "func testFunc() {\n}",
// 		},
// 		{
// 			name: "one param",
// 			function: GoFunction{
// 				Name: "testFunc",
// 				Params: GoParams{
// 					{Name: "param1", Type: GoType{TypeName: "int"}},
// 				},
// 			},
// 			want: "func testFunc(param1 int) {\n}",
// 		},
// 		{
// 			name: "multiple params",
// 			function: GoFunction{
// 				Name: "testFunc",
// 				Params: GoParams{
// 					{Name: "param1", Type: GoType{TypeName: "int"}},
// 					{Name: "param2", Type: GoType{TypeName: "string"}},
// 					{Name: "param3", Type: GoType{TypeName: "bool"}},
// 				},
// 			},
// 			want: "func testFunc(param1 int, param2 string, param3 bool) {\n}",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.function.Code(); got != tt.want {
// 				t.Errorf("GoFunction.Code() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestGoFunction_Imports(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		function GoFunction
// 		want     []string
// 	}{
// 		{
// 			name: "no imports",
// 			function: GoFunction{
// 				Name:   "testFunc",
// 				Params: GoParams{},
// 			},
// 			want: []string{},
// 		},
// 		{
// 			name: "one import",
// 			function: GoFunction{
// 				Name: "testFunc",
// 				Params: GoParams{
// 					{Name: "param1", Type: GoType{TypeName: "int", Source: "package1"}},
// 				},
// 			},
// 			want: []string{"package1"},
// 		},
// 		{
// 			name: "multiple imports",
// 			function: GoFunction{
// 				Name: "testFunc",
// 				Params: GoParams{
// 					{Name: "param1", Type: GoType{TypeName: "int", Source: "package1"}},
// 					{Name: "param2", Type: GoType{TypeName: "string", Source: "package2"}},
// 					{Name: "param3", Type: GoType{TypeName: "bool", Source: "package3"}},
// 				},
// 			},
// 			want: []string{"package1", "package2", "package3"},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.function.Imports(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GoFunction.Imports() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestGoStruct_Code(t *testing.T) {
// 	// Arrange
// 	fields := GoFields{
// 		{Name: "Field1", Type: GoType{TypeName: "string"}},
// 		{Name: "Field2", Type: GoType{TypeName: "int"}},
// 	}
// 	structName := "MyStruct"
// 	expected := "type MyStruct struct {\n  Field1 string\n  Field2 int\n}"

// 	// Act
// 	goStruct := &GoStruct{Name: structName, Fields: fields}
// 	result := goStruct.Code()

// 	// Assert
// 	if result != expected {
// 		t.Errorf("Expected '%s', but got '%s'", expected, result)
// 	}
// }

// func TestGoStruct_Imports(t *testing.T) {
// 	// Arrange
// 	fields := GoFields{
// 		{Name: "Field1", Type: GoType{TypeName: "string"}},
// 		{Name: "Field2", Type: GoType{TypeName: "time.Time", Source: "time"}},
// 	}
// 	expected := []string{"time"}

// 	// Act
// 	goStruct := &GoStruct{Name: "MyStruct", Fields: fields}
// 	result := goStruct.Imports()

// 	// Assert
// 	if !reflect.DeepEqual(result, expected) {
// 		t.Errorf("Expected '%v', but got '%v'", expected, result)
// 	}
// }
