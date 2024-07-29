package goutils

import (
	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang"
)

type FuncBuilder struct {
	funcName     string
	receiver     *golang.Receiver
	params       []*golang.Parameter
	returns      []*golang.Parameter
	body         []*golang.CodeElement
	imports      []string
	dependencies []golang.Dependency
}

func NewFunc(receiver *golang.Receiver, funcName string, params []*golang.Parameter, returns []*golang.Parameter) *FuncBuilder {
	return &FuncBuilder{
		funcName: funcName,
		receiver: receiver,
		params:   params,
		returns:  returns,
	}
}

func (fb *FuncBuilder) AddImportsAndDependencies(imports []string, dependencies []golang.Dependency) *FuncBuilder {
	fb.imports = append(fb.imports, imports...)
	fb.dependencies = append(fb.dependencies, dependencies...)
	return fb
}

func (fb *FuncBuilder) AddFunctionCall(funcCall *golang.CodeElement, errHandler *golang.ErrorHandler, cleanHandler *golang.CleanningHandler) *FuncBuilder {
	fb.body = append(fb.body, funcCall)
	fb.body = append(fb.body, &golang.CodeElement{ErrorHandler: errHandler})
	fb.body = append(fb.body, &golang.CodeElement{CleanningHandler: cleanHandler})
	return fb
}

func (fb *FuncBuilder) AddReturn(returns ...interface{}) *FuncBuilder {
	fb.body = append(fb.body, &golang.CodeElement{Return: returns})
	return fb
}

func (fb *FuncBuilder) Build() *golang.Function {
	return &golang.Function{
		Receiver:     fb.receiver,
		Name:         fb.funcName,
		Parameters:   fb.params,
		Returns:      fb.returns,
		Body:         fb.body,
		Imports:      fb.imports,
		Dependencies: fb.dependencies,
	}
}
