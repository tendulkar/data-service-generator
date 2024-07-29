package goutils

import (
	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang"
)

func CHCall(receiverName, funcName string, args ...interface{}) *golang.CleanningHandler {
	return &golang.CleanningHandler{
		Receiver: receiverName,
		Function: funcName,
		Args:     args,
	}
}

func CHStepsCall(steps []*golang.CodeElement, args ...interface{}) *golang.CleanningHandler {
	return &golang.CleanningHandler{
		Steps: steps,
		Args:  args,
	}
}
