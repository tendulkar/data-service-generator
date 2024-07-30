package goutils

import (
	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang"
)

func CHCall(receiverName, funcName string, args ...interface{}) *golang.CleanningHandler {
	return golang.NewCleanningHandler(receiverName, funcName, args, nil)
}

func CHCloseCall(receiverName string, args ...interface{}) *golang.CleanningHandler {
	return golang.NewCleanningHandler(receiverName, "Close", args, nil)
}

func CHStepsCall(steps []*golang.CodeElement, args ...interface{}) *golang.CleanningHandler {
	return golang.NewCleanningHandler("", "", args, steps)
}
