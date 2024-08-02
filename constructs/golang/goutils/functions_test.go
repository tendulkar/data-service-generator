package goutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang"
)

func TestFuncBuilder_Build(t *testing.T) {
	receiver := &golang.Receiver{
		Name: "receiver",
	}
	funcName := "funcName"
	params := []*golang.Parameter{
		{Name: "name", Type: &golang.GoType{Name: "string"}},
		{Name: "age", Type: &golang.GoType{Name: "int"}},
	}
	returns := []*golang.Parameter{
		{Name: "err", Type: &golang.GoType{Name: "error"}},
	}

	imports := []string{"fmt"}
	dependencies := []golang.Dependency{
		{Source: "github.com/pkg/errors", Version: "v0.8.1"},
	}

	builder := NewFunc(receiver, funcName, params, returns)
	builder.AddImportsAndDependencies(imports, dependencies)
	builder.AddFunctionCall(FCNewOutArgsCE([]string{"db", "err"}, "Open", &golang.Literal{Value: "postgres"}), EHReturnParams(returns), CHCall("db", "Close"))
	builder.AddReturn("db", "err")

	expectedBody := golang.CodeElements{
		{
			FunctionCall: &golang.FunctionCall{
				Function:  "Open",
				NewOutput: []string{"db", "err"},
				Args:      &golang.Literal{Value: "postgres"},
				ErrorHandler: &golang.ErrorHandler{
					ErrorFunctionReturns: returns,
				},
				CleanningHandler: &golang.CleanningHandler{
					Receiver: "db",
					Function: "Close",
				},
			},
		},
		{
			Return: []string{"db", "err"},
		},
	}

	expectedFunction := &golang.FunctionDef{
		Receiver:     receiver,
		Name:         funcName,
		Parameters:   params,
		Returns:      returns,
		Body:         expectedBody,
		Imports:      imports,
		Dependencies: dependencies,
	}

	function := builder.Build()
	funcCode, funcImports := function.FunctionCode()
	expectedFunctionCode, expectedFunctionImports := expectedFunction.FunctionCode()
	t.Log(function)
	t.Log(function.FunctionCode())
	t.Log(expectedFunction.FunctionCode())
	assert.Equal(t, expectedFunctionCode, funcCode)
	assert.Equal(t, expectedFunctionImports, funcImports)
	// assert.Equal(t, expectedFunction, function)
}
