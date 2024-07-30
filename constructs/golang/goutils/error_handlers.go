package goutils

import (
	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang"
)

func EHReturnParams(params []*golang.Parameter) *golang.ErrorHandler {
	return &golang.ErrorHandler{
		ErrorFunctionReturns: params,
	}
}

func EHError(errName string) *golang.ErrorHandler {
	return &golang.ErrorHandler{
		Error:        errName,
		ErrorReturns: []string{errName},
	}
}

func EHNilError(errName string) *golang.ErrorHandler {
	return &golang.ErrorHandler{
		Error:        errName,
		ErrorReturns: []string{"nil", errName},
	}
}

func EHFalseError(errName string) *golang.ErrorHandler {
	return &golang.ErrorHandler{
		Error:        errName,
		ErrorReturns: []string{"false", errName},
	}
}

func EHZeroError(errName string) *golang.ErrorHandler {
	return &golang.ErrorHandler{
		Error:        errName,
		ErrorReturns: []string{"0", errName},
	}
}
