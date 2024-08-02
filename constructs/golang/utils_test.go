package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateStructWithNewFunction(t *testing.T) {
	testCases := []struct {
		name               string
		structName         string
		structSuffix       string
		nameWithTypes      []NameWithType
		addJsonTag         bool
		addYamlTag         bool
		addDBTag           bool
		expectedStructCode string
		expectedFuncCode   string
	}{
		{
			name:         "Test case 1: Simple struct with two fields",
			structName:   "User",
			structSuffix: "",
			nameWithTypes: []NameWithType{
				{Name: "full_name", Type: &GoType{Name: "string", Source: ""}},
				{Name: "age", Type: &GoType{Name: "int", Source: ""}},
				{Name: "is_married", Type: &GoType{Name: "bool", Source: ""}},
				{Name: "height", Type: &GoType{Name: "float64", Source: ""}},
			},
			addJsonTag: true,
			addYamlTag: false,
			addDBTag:   true,
			expectedStructCode: `type User struct {
	FullName string	` + "`json:\"full_name\" db:\"full_name\"`" + `
	Age int	` + "`json:\"age\" db:\"age\"`" + `
	IsMarried bool	` + "`json:\"is_married\" db:\"is_married\"`" + `
	Height float64	` + "`json:\"height\" db:\"height\"`" + `
}
`,
			expectedFuncCode: `func NewUser(fullName string, age int, isMarried bool, height float64) *User {
	return &User{
		FullName: fullName,
		Age: age,
		IsMarried: isMarried,
		Height: height,
	}
}`,
		},
		// Add more test cases here...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			st, fn := GenStructWithNewFunction(tc.structName, tc.nameWithTypes, false, tc.addJsonTag, tc.addYamlTag, tc.addDBTag)
			// t.Log(st.StructCode())
			// t.Log(fn.FunctionCode())
			stCode, _ := st.StructCode()
			fnCode, _ := fn.FunctionCode()
			assert.Equal(t, tc.expectedStructCode, stCode)
			assert.Equal(t, tc.expectedFuncCode, fnCode)
		})
	}
}
