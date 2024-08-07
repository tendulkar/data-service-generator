package golang

func AddCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		Add: &Add{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func AddOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		Add: &Add{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}

func AddNewOutCE(newOutput, a, b interface{}) *CodeElement {
	return &CodeElement{
		Add: &Add{BinaryOp: BinaryOp{Left: a, Right: b, NewOutput: newOutput}},
	}
}

func SubCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		Subtract: &Sub{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func SubOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		Subtract: &Sub{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}

func SubNewOutCE(newOutput, a, b interface{}) *CodeElement {
	return &CodeElement{
		Subtract: &Sub{BinaryOp: BinaryOp{Left: a, Right: b, NewOutput: newOutput}},
	}
}

func MulCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		Multiply: &Mul{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func MulOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		Multiply: &Mul{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}

func MulNewOutCE(newOutput, a, b interface{}) *CodeElement {
	return &CodeElement{
		Multiply: &Mul{BinaryOp: BinaryOp{Left: a, Right: b, NewOutput: newOutput}},
	}
}

func DivCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		Divide: &Div{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func DivOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		Divide: &Div{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}

func DivNewOutCE(newOutput, a, b interface{}) *CodeElement {
	return &CodeElement{
		Divide: &Div{BinaryOp: BinaryOp{Left: a, Right: b, NewOutput: newOutput}},
	}
}

func ModCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		Modulo: &Mod{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func ModOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		Modulo: &Mod{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}

func ModNewOutCE(newOutput, a, b interface{}) *CodeElement {
	return &CodeElement{
		Modulo: &Mod{BinaryOp: BinaryOp{Left: a, Right: b, NewOutput: newOutput}},
	}
}

func AndCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		And: &And{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func AndOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		And: &And{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}

func AndNewOutCE(newOutput, a, b interface{}) *CodeElement {
	return &CodeElement{
		And: &And{BinaryOp: BinaryOp{Left: a, Right: b, NewOutput: newOutput}},
	}
}

func OrCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		Or: &Or{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func OrOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		Or: &Or{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}

func OrNewOutCE(newOutput, a, b interface{}) *CodeElement {
	return &CodeElement{
		Or: &Or{BinaryOp: BinaryOp{Left: a, Right: b, NewOutput: newOutput}},
	}
}
func NotCE(a interface{}) *CodeElement {
	return &CodeElement{
		Not: &Not{UnaryOp: UnaryOp{Input: a}},
	}
}

func NotOutCE(output, a interface{}) *CodeElement {
	return &CodeElement{
		Not: &Not{UnaryOp: UnaryOp{Input: a, Output: output}},
	}
}

func NotNewOutCE(newOutput, a interface{}) *CodeElement {
	return &CodeElement{
		Not: &Not{UnaryOp: UnaryOp{Input: a, NewOutput: newOutput}},
	}
}

func EqCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		Equal: &Equal{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func EqOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		Equal: &Equal{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}

func EqNewOutCE(newOutput, a, b interface{}) *CodeElement {
	return &CodeElement{
		Equal: &Equal{BinaryOp: BinaryOp{Left: a, Right: b, NewOutput: newOutput}},
	}
}

func NeCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		NotEqual: &NotEqual{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func NeOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		NotEqual: &NotEqual{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}

func NeNewOutCE(newOutput, a, b interface{}) *CodeElement {
	return &CodeElement{
		NotEqual: &NotEqual{BinaryOp: BinaryOp{Left: a, Right: b, NewOutput: newOutput}},
	}
}

func GtCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		GreaterThan: &GreaterThan{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func GtOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		GreaterThan: &GreaterThan{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}

func GtNewOutCE(newOutput, a, b interface{}) *CodeElement {
	return &CodeElement{
		GreaterThan: &GreaterThan{BinaryOp: BinaryOp{Left: a, Right: b, NewOutput: newOutput}},
	}
}

func GeCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		GreaterThanOrEqual: &GreaterThanOrEqual{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func GeOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		GreaterThanOrEqual: &GreaterThanOrEqual{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}

func GeNewOutCE(newOutput, a, b interface{}) *CodeElement {
	return &CodeElement{
		GreaterThanOrEqual: &GreaterThanOrEqual{BinaryOp: BinaryOp{Left: a, Right: b, NewOutput: newOutput}},
	}
}

func LtCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		LessThan: &LessThan{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func LtOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		LessThan: &LessThan{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}

func LtNewOutCE(newOutput, a, b interface{}) *CodeElement {
	return &CodeElement{
		LessThan: &LessThan{BinaryOp: BinaryOp{Left: a, Right: b, NewOutput: newOutput}},
	}
}

func LeCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		LessThanOrEqual: &LessThanOrEqual{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func LeOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		LessThanOrEqual: &LessThanOrEqual{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}
func LsCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		LeftShift: &LeftShift{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func LsOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		LeftShift: &LeftShift{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}

func LsNewOutCE(newOutput, a, b interface{}) *CodeElement {
	return &CodeElement{
		LeftShift: &LeftShift{BinaryOp: BinaryOp{Left: a, Right: b, NewOutput: newOutput}},
	}
}

func RsCE(a, b interface{}) *CodeElement {
	return &CodeElement{
		RightShift: &RightShift{BinaryOp: BinaryOp{Left: a, Right: b}},
	}
}

func RsOutCE(output, a, b interface{}) *CodeElement {
	return &CodeElement{
		RightShift: &RightShift{BinaryOp: BinaryOp{Left: a, Right: b, Output: output}},
	}
}

func RsNewOutCE(newOutput, a, b interface{}) *CodeElement {
	return &CodeElement{
		RightShift: &RightShift{BinaryOp: BinaryOp{Left: a, Right: b, NewOutput: newOutput}},
	}
}
