package goutils

import (
	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang"
)

func FCCE(funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, nil, false, false,
		nil, nil)
}

func FCArgsCE(funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, args, false, false,
		nil, nil)
}

func FCReceiverArgsCE(receiver string, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, args, false, false,
		nil, nil)
}

func FCReceiverCE(receiver string, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, nil, false, false,
		nil, nil)
}

func FCAsyncCE(funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, nil, true, false,
		nil, nil)
}

func FCArgsAsyncCE(funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, args, true, false,
		nil, nil)
}

func FCReceiverArgsAsyncCE(receiver string, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, args, true, false,
		nil, nil)
}

func FCReceiverAsyncCE(receiver string, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, nil, true, false,
		nil, nil)
}

func FCDeferCE(funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, nil, false, true,
		nil, nil)
}

func FCArgsDeferCE(funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, args, false, true,
		nil, nil)
}

func FCReceiverArgsDeferCE(receiver string, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, args, false, true,
		nil, nil)
}

func FCReceiverDeferCE(receiver string, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, nil, false, true,
		nil, nil)
}

func FCOutCE(output interface{}, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, nil, false, false,
		nil, nil)
}

func FCOutArgsCE(output interface{}, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, args, false, false,
		nil, nil)
}

func FCOutReceiverArgsCE(output interface{}, receiver string, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, args, false, false,
		nil, nil)
}

func FCOutReceiverCE(output interface{}, receiver string, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, nil, false, false,
		nil, nil)
}

func FCOutAsyncCE(output interface{}, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, nil, true, false,
		nil, nil)
}

func FCOutArgsAsyncCE(output interface{}, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, args, true, false,
		nil, nil)
}

func FCOutReceiverArgsAsyncCE(output interface{}, receiver string, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, args, true, false,
		nil, nil)
}

func FCOutReceiverAsyncCE(output interface{}, receiver string, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, nil, true, false,
		nil, nil)
}

func FCOutDeferCE(output interface{}, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, nil, false, true,
		nil, nil)
}

func FCOutArgsDeferCE(output interface{}, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, args, false, true,
		nil, nil)
}

func FCOutReceiverArgsDeferCE(output interface{}, receiver string, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, args, false, true,
		nil, nil)
}

func FCOutReceiverDeferCE(output interface{}, receiver string, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, nil, false, true,
		nil, nil)
}

func FCNewOutCE(newOutput interface{}, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, nil, false, false,
		nil, nil)
}

func FCNewOutArgsCE(newOutput interface{}, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, args, false, false,
		nil, nil)
}

func FCNewOutReceiverArgsCE(newOutput interface{}, receiver string, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, args, false, false,
		nil, nil)
}

func FCNewOutReceiverCE(newOutput interface{}, receiver string, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, nil, false, false,
		nil, nil)
}

func FCNewOutAsyncCE(newOutput interface{}, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, nil, true, false,
		nil, nil)
}

func FCNewOutArgsAsyncCE(newOutput interface{}, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, args, true, false,
		nil, nil)
}

func FCNewOutReceiverArgsAsyncCE(newOutput interface{}, receiver string, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, args, true, false,
		nil, nil)
}

func FCNewOutReceiverAsyncCE(newOutput interface{}, receiver string, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, nil, true, false,
		nil, nil)
}

func FCNewOutDeferCE(newOutput interface{}, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, nil, false, true,
		nil, nil)
}

func FCNewOutArgsDeferCE(newOutput interface{}, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, args, false, true,
		nil, nil)
}

func FCNewOutReceiverArgsDeferCE(newOutput interface{}, receiver string, funcName string, args interface{}) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, args, false, true,
		nil, nil)
}

func FCNewOutReceiverDeferCE(newOutput interface{}, receiver string, funcName string) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, nil, false, true,
		nil, nil)
}

func FCEHCE(output interface{}, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, nil, false, false,
		errorHandler, nil)
}

func FCEHArgsCE(output interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, args, false, false,
		errorHandler, nil)
}

func FCEHReceiverArgsCE(output interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, args, false, false,
		errorHandler, nil)
}

func FCEHReceiverCE(output interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, nil, false, false,
		errorHandler, nil)
}

func FCEHAsyncCE(output interface{}, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, nil, true, false,
		errorHandler, nil)
}

func FCEHArgsAsyncCE(output interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, args, true, false,
		errorHandler, nil)
}

func FCEHReceiverArgsAsyncCE(output interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, args, true, false,
		errorHandler, nil)
}

func FCEHReceiverAsyncCE(output interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, nil, true, false,
		errorHandler, nil)
}

func FCEHDeferCE(output interface{}, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, nil, false, true,
		errorHandler, nil)
}

func FCEHArgsDeferCE(output interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, args, false, true,
		errorHandler, nil)
}

func FCEHReceiverArgsDeferCE(output interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, args, false, true,
		errorHandler, nil)
}

func FCEHReceiverDeferCE(output interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, nil, false, true,
		errorHandler, nil)
}

func FCEHOutCE(output interface{}, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, nil, false, false,
		errorHandler, nil)
}

func FCEHOutArgsCE(output interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, args, false, false,
		errorHandler, nil)
}

func FCEHOutReceiverArgsCE(output interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, args, false, false,
		errorHandler, nil)
}

func FCEHOutReceiverCE(output interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, nil, false, false,
		errorHandler, nil)
}

func FCEHOutAsyncCE(output interface{}, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, nil, true, false,
		errorHandler, nil)
}

func FCEHOutArgsAsyncCE(output interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, args, true, false,
		errorHandler, nil)
}

func FCEHOutReceiverArgsAsyncCE(output interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, args, true, false,
		errorHandler, nil)
}

func FCEHOutReceiverAsyncCE(output interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, nil, true, false,
		errorHandler, nil)
}

func FCEHOutDeferCE(output interface{}, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, nil, false, true,
		errorHandler, nil)
}

func FCEHOutArgsDeferCE(output interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, args, false, true,
		errorHandler, nil)
}

func FCEHOutReceiverArgsDeferCE(output interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, args, false, true,
		errorHandler, nil)
}

func FCEHOutReceiverDeferCE(output interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, nil, false, true,
		errorHandler, nil)
}

func FCEHNewOutCE(output interface{}, newOutput interface{}, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, nil, false, false,
		errorHandler, nil)
}

func FCEHNewOutArgsCE(output interface{}, newOutput interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, args, false, false,
		errorHandler, nil)
}

func FCEHNewOutReceiverArgsCE(output interface{}, newOutput interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, args, false, false,
		errorHandler, nil)
}

func FCEHNewOutReceiverCE(output interface{}, newOutput interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, nil, false, false,
		errorHandler, nil)
}

func FCEHNewOutAsyncCE(output interface{}, newOutput interface{}, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, nil, true, false,
		errorHandler, nil)
}

func FCEHNewOutArgsAsyncCE(output interface{}, newOutput interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, args, true, false,
		errorHandler, nil)
}

func FCEHNewOutReceiverArgsAsyncCE(output interface{}, newOutput interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, args, true, false,
		errorHandler, nil)
}

func FCEHNewOutReceiverAsyncCE(output interface{}, newOutput interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, nil, true, false,
		errorHandler, nil)
}

func FCEHNewOutDeferCE(output interface{}, newOutput interface{}, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, nil, false, true,
		errorHandler, nil)
}

func FCEHNewOutArgsDeferCE(output interface{}, newOutput interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, args, false, true,
		errorHandler, nil)
}

func FCEHNewOutReceiverArgsDeferCE(output interface{}, newOutput interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, args, false, true,
		errorHandler, nil)
}

func FCEHNewOutReceiverDeferCE(output interface{}, newOutput interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, nil, false, true,
		errorHandler, nil)
}

func FCCHCE(funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, nil, false, false,
		nil, cleanningHandler)
}

func FCCHArgsCE(funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, args, false, false,
		nil, cleanningHandler)
}

func FCCHReceiverArgsCE(receiver string, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, args, false, false,
		nil, cleanningHandler)
}

func FCCHReceiverCE(receiver string, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, nil, false, false,
		nil, cleanningHandler)
}

func FCCHAsyncCE(funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, nil, true, false,
		nil, cleanningHandler)
}

func FCCHArgsAsyncCE(funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, args, true, false,
		nil, cleanningHandler)
}

func FCCHReceiverArgsAsyncCE(receiver string, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, args, true, false,
		nil, cleanningHandler)
}

func FCCHReceiverAsyncCE(receiver string, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, nil, true, false,
		nil, cleanningHandler)
}

func FCCHDeferCE(funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, nil, false, true,
		nil, cleanningHandler)
}

func FCCHArgsDeferCE(funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, args, false, true,
		nil, cleanningHandler)
}

func FCCHReceiverArgsDeferCE(receiver string, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, args, false, true,
		nil, cleanningHandler)
}

func FCCHReceiverDeferCE(receiver string, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, nil, false, true,
		nil, cleanningHandler)
}

func FCCHOutCE(output interface{}, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, nil, false, false,
		nil, cleanningHandler)
}

func FCCHOutArgsCE(output interface{}, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, args, false, false,
		nil, cleanningHandler)
}

func FCCHOutReceiverArgsCE(output interface{}, receiver string, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, args, false, false,
		nil, cleanningHandler)
}

func FCCHOutReceiverCE(output interface{}, receiver string, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, nil, false, false,
		nil, cleanningHandler)
}

func FCCHOutAsyncCE(output interface{}, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, nil, true, false,
		nil, cleanningHandler)
}

func FCCHOutArgsAsyncCE(output interface{}, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, args, true, false,
		nil, cleanningHandler)
}

func FCCHOutReceiverArgsAsyncCE(output interface{}, receiver string, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, args, true, false,
		nil, cleanningHandler)
}

func FCCHOutReceiverAsyncCE(output interface{}, receiver string, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, nil, true, false,
		nil, cleanningHandler)
}

func FCCHOutDeferCE(output interface{}, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, nil, false, true,
		nil, cleanningHandler)
}

func FCCHOutArgsDeferCE(output interface{}, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, args, false, true,
		nil, cleanningHandler)
}

func FCCHOutReceiverArgsDeferCE(output interface{}, receiver string, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, args, false, true,
		nil, cleanningHandler)
}

func FCCHOutReceiverDeferCE(output interface{}, receiver string, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, nil, false, true,
		nil, cleanningHandler)
}

func FCCHNewOutCE(newOutput interface{}, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, nil, false, false,
		nil, cleanningHandler)
}

func FCCHNewOutArgsCE(newOutput interface{}, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, args, false, false,
		nil, cleanningHandler)
}

func FCCHNewOutReceiverArgsCE(newOutput interface{}, receiver string, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, args, false, false,
		nil, cleanningHandler)
}

func FCCHNewOutReceiverCE(newOutput interface{}, receiver string, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, nil, false, false,
		nil, cleanningHandler)
}

func FCCHNewOutAsyncCE(newOutput interface{}, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, nil, true, false,
		nil, cleanningHandler)
}

func FCCHNewOutArgsAsyncCE(newOutput interface{}, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, args, true, false,
		nil, cleanningHandler)
}

func FCCHNewOutReceiverArgsAsyncCE(newOutput interface{}, receiver string, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, args, true, false,
		nil, cleanningHandler)
}

func FCCHNewOutReceiverAsyncCE(newOutput interface{}, receiver string, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, nil, true, false,
		nil, cleanningHandler)
}

func FCCHNewOutDeferCE(newOutput interface{}, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, nil, false, true,
		nil, cleanningHandler)
}

func FCCHNewOutArgsDeferCE(newOutput interface{}, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, args, false, true,
		nil, cleanningHandler)
}

func FCCHNewOutReceiverArgsDeferCE(newOutput interface{}, receiver string, funcName string, args interface{}, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, args, false, true,
		nil, cleanningHandler)
}

func FCCHNewOutReceiverDeferCE(newOutput interface{}, receiver string, funcName string, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, nil, false, true,
		nil, cleanningHandler)
}

func FCEHCHCE(output interface{}, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, nil, false, false,
		errorHandler, cleanningHandler)
}

func FCEHCHArgsCE(output interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, args, false, false,
		errorHandler, cleanningHandler)
}

func FCEHCHReceiverArgsCE(output interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, args, false, false,
		errorHandler, cleanningHandler)
}

func FCEHCHReceiverCE(output interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, nil, false, false,
		errorHandler, cleanningHandler)
}

func FCEHCHAsyncCE(output interface{}, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, nil, true, false,
		errorHandler, cleanningHandler)
}

func FCEHCHArgsAsyncCE(output interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, args, true, false,
		errorHandler, cleanningHandler)
}

func FCEHCHReceiverArgsAsyncCE(output interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, args, true, false,
		errorHandler, cleanningHandler)
}

func FCEHCHReceiverAsyncCE(output interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, nil, true, false,
		errorHandler, cleanningHandler)
}

func FCEHCHDeferCE(output interface{}, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, nil, false, true,
		errorHandler, cleanningHandler)
}

func FCEHCHArgsDeferCE(output interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, "", funcName, args, false, true,
		errorHandler, cleanningHandler)
}

func FCEHCHReceiverArgsDeferCE(output interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, args, false, true,
		errorHandler, cleanningHandler)
}

func FCEHCHReceiverDeferCE(output interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, nil, receiver, funcName, nil, false, true,
		errorHandler, cleanningHandler)
}

func FCEHCHOutCE(output interface{}, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, nil, false, false,
		errorHandler, cleanningHandler)
}

func FCEHCHOutArgsCE(output interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, args, false, false,
		errorHandler, cleanningHandler)
}

func FCEHCHOutReceiverArgsCE(output interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, args, false, false,
		errorHandler, cleanningHandler)
}

func FCEHCHOutReceiverCE(output interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, nil, false, false,
		errorHandler, cleanningHandler)
}

func FCEHCHOutAsyncCE(output interface{}, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, nil, true, false,
		errorHandler, cleanningHandler)
}

func FCEHCHOutArgsAsyncCE(output interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, args, true, false,
		errorHandler, cleanningHandler)
}

func FCEHCHOutReceiverArgsAsyncCE(output interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, args, true, false,
		errorHandler, cleanningHandler)
}

func FCEHCHOutReceiverAsyncCE(output interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, nil, true, false,
		errorHandler, cleanningHandler)
}

func FCEHCHOutDeferCE(output interface{}, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, nil, false, true,
		errorHandler, cleanningHandler)
}

func FCEHCHOutArgsDeferCE(output interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, "", funcName, args, false, true,
		errorHandler, cleanningHandler)
}

func FCEHCHOutReceiverArgsDeferCE(output interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, args, false, true,
		errorHandler, cleanningHandler)
}

func FCEHCHOutReceiverDeferCE(output interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(nil, output, receiver, funcName, nil, false, true,
		errorHandler, cleanningHandler)
}

func FCEHCHNewOutCE(output interface{}, newOutput interface{}, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, nil, false, false,
		errorHandler, cleanningHandler)
}

func FCEHCHNewOutArgsCE(output interface{}, newOutput interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, args, false, false,
		errorHandler, cleanningHandler)
}

func FCEHCHNewOutReceiverArgsCE(output interface{}, newOutput interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, args, false, false,
		errorHandler, cleanningHandler)
}

func FCEHCHNewOutReceiverCE(output interface{}, newOutput interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, nil, false, false,
		errorHandler, cleanningHandler)
}

func FCEHCHNewOutAsyncCE(output interface{}, newOutput interface{}, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, nil, true, false,
		errorHandler, cleanningHandler)
}

func FCEHCHNewOutArgsAsyncCE(output interface{}, newOutput interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, args, true, false,
		errorHandler, cleanningHandler)
}

func FCEHCHNewOutReceiverArgsAsyncCE(output interface{}, newOutput interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, args, true, false,
		errorHandler, cleanningHandler)
}

func FCEHCHNewOutReceiverAsyncCE(output interface{}, newOutput interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, nil, true, false,
		errorHandler, cleanningHandler)
}

func FCEHCHNewOutDeferCE(output interface{}, newOutput interface{}, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, nil, false, true,
		errorHandler, cleanningHandler)
}

func FCEHCHNewOutArgsDeferCE(output interface{}, newOutput interface{}, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, "", funcName, args, false, true,
		errorHandler, cleanningHandler)
}

func FCEHCHNewOutReceiverArgsDeferCE(output interface{}, newOutput interface{}, receiver string, funcName string, args interface{}, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, args, false, true,
		errorHandler, cleanningHandler)
}

func FCEHCHNewOutReceiverDeferCE(output interface{}, newOutput interface{}, receiver string, funcName string, errorHandler *golang.ErrorHandler, cleanningHandler *golang.CleanningHandler) *golang.CodeElement {
	return golang.FunctionCallCE(newOutput, nil, receiver, funcName, nil, false, true,
		errorHandler, cleanningHandler)
}
