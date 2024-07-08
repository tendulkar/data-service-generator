package golang

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
	"stellarsky.ai/platform/codegen/data-service-generator/base"
)

type CodeBlock interface {
	ToCode() string
}

// CodeElement struct handles all operations and control structures
type CodeElement struct {
	Add                interface{}                 `yaml:"add,omitempty"`
	Subtract           interface{}                 `yaml:"sub,omitempty"`
	Multiply           interface{}                 `yaml:"mul,omitempty"`
	Divide             interface{}                 `yaml:"div,omitempty"`
	Modulo             interface{}                 `yaml:"mod,omitempty"`
	And                interface{}                 `yaml:"and,omitempty"`
	Or                 interface{}                 `yaml:"or,omitempty"`
	Not                interface{}                 `yaml:"not,omitempty"`
	Equal              interface{}                 `yaml:"eq,omitempty"`
	NotEqual           interface{}                 `yaml:"ne,omitempty"`
	LessThan           interface{}                 `yaml:"lt,omitempty"`
	LessThanOrEqual    interface{}                 `yaml:"le,omitempty"`
	GreaterThan        interface{}                 `yaml:"gt,omitempty"`
	GreaterThanOrEqual interface{}                 `yaml:"ge,omitempty"`
	LeftShift          interface{}                 `yaml:"lshift,omitempty"`
	RightShift         interface{}                 `yaml:"rshift,omitempty"`
	BitwiseAnd         interface{}                 `yaml:"bit_and,omitempty"`
	BitwiseOr          interface{}                 `yaml:"bit_or,omitempty"`
	BitwiseXor         interface{}                 `yaml:"bit_xor,omitempty"`
	BitwiseNot         interface{}                 `yaml:"bit_not,omitempty"`
	PostIncrement      interface{}                 `yaml:"post_inc,omitempty"`
	PostDecrement      interface{}                 `yaml:"post_dec,omitempty"`
	PreIncrement       interface{}                 `yaml:"pre_inc,omitempty"`
	PreDecrement       interface{}                 `yaml:"pre_dec,omitempty"`
	Assign             *Assignment                 `yaml:"assign,omitempty"`
	NewAssign          *NewAssignment              `yaml:"new_assign,omitempty"`
	If                 *IfElement                  `yaml:"if,omitempty"`
	Cases              CaseElement                 `yaml:"cases,omitempty"`
	MatchCases         MatchCases                  `yaml:"match_cases,omitempty"`
	RepeatCond         *RepeatByCondition          `yaml:"repeat_cond,omitempty"`
	RepeatInitCond     *RepeatInitConditionElement `yaml:"repeat_init_cond,omitempty"`
	RepeatLoop         *RepeatLoopElement          `yaml:"repeat,omitempty"`
	Iterate            *IterateElement             `yaml:"iterate,omitempty"`
	Return             []string                    `yaml:"return,omitempty"`
	StructCreation     *StructCreation             `yaml:"create,omitempty"`
	GoRoutine          *GoRoutine                  `yaml:"async,omitempty"`
	FunctionCall       *FunctionCall               `yaml:"call,omitempty"`
	MemberFunctionCall *MemberFunctionCall         `yaml:"obj_call,omitempty"`
	Steps              []*CodeElement              `yaml:"steps,omitempty"`
	Imports            []string                    `yaml:"imports,omitempty"`
	Literal            interface{}                 `yaml:"lit,omitempty"`
}

type BinaryOp struct {
	Left      interface{} `yaml:"left"`
	Right     interface{} `yaml:"right"`
	Output    string      `yaml:"out,omitempty"`
	NewOutput string      `yaml:"nout,omitempty"`
}

type UnaryOp struct {
	Input     interface{} `yaml:"in"`
	Output    interface{} `yaml:"out,omitempty"`
	NewOutput string      `yaml:"nout,omitempty"`
}

type Add struct {
	BinaryOp `yaml:",inline"`
}

type Sub struct {
	BinaryOp `yaml:",inline"`
}

type Mul struct {
	BinaryOp `yaml:",inline"`
}

type Div struct {
	BinaryOp `yaml:",inline"`
}

type Mod struct {
	BinaryOp `yaml:",inline"`
}

type And struct {
	BinaryOp `yaml:",inline"`
}

type Or struct {
	BinaryOp `yaml:",inline"`
}

type Not struct {
	UnaryOp `yaml:",inline"`
}

type Equal struct {
	BinaryOp `yaml:",inline"`
}

type NotEqual struct {
	BinaryOp `yaml:",inline"`
}

type LessThan struct {
	BinaryOp `yaml:",inline"`
}

type LessThanOrEqual struct {
	BinaryOp `yaml:",inline"`
}

type GreaterThan struct {
	BinaryOp `yaml:",inline"`
}

type GreaterThanOrEqual struct {
	BinaryOp `yaml:",inline"`
}

type LeftShift struct {
	BinaryOp `yaml:",inline"`
}

type RightShift struct {
	BinaryOp `yaml:",inline"`
}

type BitwiseAnd struct {
	BinaryOp `yaml:",inline"`
}

type BitwiseOr struct {
	BinaryOp `yaml:",inline"`
}

type BitwiseXor struct {
	BinaryOp `yaml:",inline"`
}

type BitwiseNot struct {
	UnaryOp `yaml:",inline"`
}

type PostIncrement struct {
	UnaryOp `yaml:",inline"`
}

type PostDecrement struct {
	UnaryOp `yaml:",inline"`
}

type PreIncrement struct {
	UnaryOp `yaml:",inline"`
}

type PreDecrement struct {
	UnaryOp `yaml:",inline"`
}

// Supporting structs
type Assignment struct {
	Left  interface{} `yaml:"left"`
	Right interface{} `yaml:"right"`
}

type NewAssignment struct {
	Left  interface{} `yaml:"left"`
	Right interface{} `yaml:"right"`
}

type IfElement struct {
	Condition *CodeElement   `yaml:"cond"`
	Then      []*CodeElement `yaml:"then"`
	Else      []*CodeElement `yaml:"else,omitempty"`
}

type OneCaseElement struct {
	Condition interface{}    `yaml:"cond"`
	Body      []*CodeElement `yaml:"body"`
}

type CaseElement []*OneCaseElement

type MatchCase struct {
	MatchWith interface{}    `yaml:"match_with,omitempty"`
	Body      []*CodeElement `yaml:"body"`
}

type MatchCases struct {
	MatchOn    interface{}  `yaml:"match"`
	MatchCases []*MatchCase `yaml:"cases"`
}

// type MatchCases []*MatchCase

type RepeatByCondition struct {
	Condition *CodeElement `yaml:"cond"`
	Body      []*CodeElement
}

type RepeatInitConditionElement struct {
	Init      []*CodeElement
	Condition *CodeElement   `yaml:"cond"`
	Body      []*CodeElement `yaml:"body"`
}

type RepeatLoopElement struct {
	Init      []*CodeElement `yaml:"init"`
	Condition *CodeElement   `yaml:"cond"`
	Step      []*CodeElement `yaml:"step"`
	Body      []*CodeElement `yaml:"body"`
}

type IterateElement struct {
	Variables []string       `yaml:"variables"`
	RangeOn   *CodeElement   `yaml:"range_on"`
	Body      []*CodeElement `yaml:"body"`
}

type StructCreation struct {
	Output     interface{}    `yaml:"out,omitempty"`
	NewOutput  interface{}    `yaml:"nout,omitempty"`
	StructType string         `yaml:"struct_type"`
	Values     []*CodeElement `yaml:"values"`
}

type GoRoutine struct {
	FunctionCall       *CodeElement `yaml:"call"`
	MemberFunctionCall *CodeElement `yaml:"obj_call"`
}

type FunctionCall struct {
	Output    interface{}    `yaml:"out,omitempty"`
	NewOutput interface{}    `yaml:"nout,omitempty"`
	Function  string         `yaml:"function"`
	Params    []*CodeElement `yaml:"params"`
}

type MemberFunctionCall struct {
	Output    interface{}    `yaml:"out,omitempty"`
	NewOutput interface{}    `yaml:"nout,omitempty"`
	Receiver  string         `yaml:"obj"`
	Function  string         `yaml:"function"`
	Params    []*CodeElement `yaml:"params"`
}

type Literal struct {
	Value interface{} `yaml:"val"`
	Type  string      `yaml:"type"`
}

func resolveTypeLiteral(v interface{}, t string) string {
	switch t {
	case "string":
		return fmt.Sprintf("\"%s\"", v)
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
		return fmt.Sprintf("%v", v)
	case "array_int", "int_array", "array_int32", "int32_array":
		return fmt.Sprintf("[]int{%v}", strings.Join(v.([]string), ", "))
	case "array_uint", "uint_array", "array_uint32", "uint32_array":
		return fmt.Sprintf("[]uint{%v}", strings.Join(v.([]string), ", "))
	case "array_int64", "int64_array":
		return fmt.Sprintf("[]int64{%v}", strings.Join(v.([]string), ", "))
	case "array_uint64", "uint64_array":
		return fmt.Sprintf("[]uint64{%v}", strings.Join(v.([]string), ", "))
	case "array_string", "string_array":
		return fmt.Sprintf("[]string{%v}", strings.Join(v.([]string), ", "))
	case "array_float", "float_array", "array_float32", "float32_array":
		return fmt.Sprintf("[]float64{%v}", strings.Join(v.([]string), ", "))
	case "array_float64", "float64_array":
		return fmt.Sprintf("[]float64{%v}", strings.Join(v.([]string), ", "))
	case "array_bool", "bool_array":
		return fmt.Sprintf("[]bool{%v}", strings.Join(v.([]string), ", "))

	default:
		return fmt.Sprintf("%v", v)
	}
}

func resolveArrayInterface(arr []interface{}, sep string) string {
	values := make([]string, len(arr))
	for i, v := range arr {
		values[i] = resolveCodeElementOrString(v)
	}

	return strings.Join(values, sep)
}

// write a literal using yaml and resolve that don't need quotes
// we don't want to write golag code, either quotes, or brackets for arrays, we just want the literal
// This will be platform compatible
func resolveLiteral(v interface{}) string {
	switch rv := v.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", rv)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", rv)
	case Literal:
		if rv.Type == "" {
			return resolveLiteral(rv.Value)
		} else {
			return resolveTypeLiteral(rv.Value, rv.Type)
		}

	case []interface{}:
		return resolveArrayInterface(rv, ", ")

	default:
		return fmt.Sprintf("%v", rv)
	}
}

func resolveCodeElementOrString(v interface{}) string {
	switch rv := v.(type) {
	case *Literal:
		return resolveLiteral(rv.Value)
	case *CodeElement:
		return rv.ToCode()
	default:
		return fmt.Sprintf("%v", rv)
	}
}

func binaryOpToCode(op string, b *BinaryOp) string {
	lvalue := resolveCodeElementOrString(b.Left)
	rvalue := resolveCodeElementOrString(b.Right)

	if b.Output == "" && b.NewOutput == "" {
		return fmt.Sprintf("%s %s %s", lvalue, op, rvalue)
	} else if b.NewOutput == "" {
		return fmt.Sprintf("%s = %s %s %s", b.Output, lvalue, op, rvalue)
	} else {
		return fmt.Sprintf("%s := %s %s %s", b.NewOutput, lvalue, op, rvalue)
	}
}

func unaryOpToCode(op string, u *UnaryOp) string {
	inpValue := resolveCodeElementOrString(u.Input)

	if u.Output == "" && u.NewOutput == "" {
		return fmt.Sprintf("%s %s", op, inpValue)
	} else if u.NewOutput == "" {
		return fmt.Sprintf("%s = %s%s", u.Output, op, inpValue)
	} else {
		return fmt.Sprintf("%s := %s%s", u.NewOutput, op, inpValue)
	}
}

func unaryPostOpToCode(op string, u *UnaryOp) string {
	inpValue := resolveCodeElementOrString(u.Input)
	if u.Output == "" && u.NewOutput == "" {
		return fmt.Sprintf("%s %s", inpValue, op)
	} else if u.NewOutput == "" {
		return fmt.Sprintf("%s = %s %s", u.Output, inpValue, op)
	} else {
		return fmt.Sprintf("%s := %s %s", u.NewOutput, inpValue, op)
	}
}

func (a *Add) ToCode() string {
	return binaryOpToCode("+", &a.BinaryOp)
}

func (s *Sub) ToCode() string {
	return binaryOpToCode("-", &s.BinaryOp)
}

func (m *Mul) ToCode() string {
	return binaryOpToCode("*", &m.BinaryOp)
}

func (d *Div) ToCode() string {
	return binaryOpToCode("/", &d.BinaryOp)
}

func (m *Mod) ToCode() string {
	return binaryOpToCode("%", &m.BinaryOp)
}

func (a *And) ToCode() string {
	return binaryOpToCode("&&", &a.BinaryOp)
}

func (o *Or) ToCode() string {
	return binaryOpToCode("||", &o.BinaryOp)
}

func (n *Not) ToCode() string {
	return unaryOpToCode("!", &n.UnaryOp)
}

func (e *Equal) ToCode() string {
	return binaryOpToCode("==", &e.BinaryOp)
}

func (ne *NotEqual) ToCode() string {
	return binaryOpToCode("!=", &ne.BinaryOp)
}

func (l *LessThan) ToCode() string {
	return binaryOpToCode("<", &l.BinaryOp)
}

func (le *LessThanOrEqual) ToCode() string {
	return binaryOpToCode("<=", &le.BinaryOp)
}

func (l *GreaterThan) ToCode() string {
	return binaryOpToCode(">", &l.BinaryOp)
}

func (ge *GreaterThanOrEqual) ToCode() string {
	return binaryOpToCode(">=", &ge.BinaryOp)
}

func (lshift *LeftShift) ToCode() string {
	return binaryOpToCode("<<", &lshift.BinaryOp)
}

func (rshift *RightShift) ToCode() string {
	return binaryOpToCode(">>", &rshift.BinaryOp)
}

func (b *BitwiseAnd) ToCode() string {
	return binaryOpToCode("&", &b.BinaryOp)
}

func (b *BitwiseOr) ToCode() string {
	return binaryOpToCode("|", &b.BinaryOp)
}

func (b *BitwiseNot) ToCode() string {
	return unaryOpToCode("^", &b.UnaryOp)
}

func (b *BitwiseXor) ToCode() string {
	return binaryOpToCode("^", &b.BinaryOp)
}

func (p *PostIncrement) ToCode() string {
	return unaryPostOpToCode("++", &p.UnaryOp)
}

func (p *PostDecrement) ToCode() string {
	return unaryPostOpToCode("--", &p.UnaryOp)
}

func (p *PreIncrement) ToCode() string {
	return unaryOpToCode("++", &p.UnaryOp)
}

func (p *PreDecrement) ToCode() string {
	return unaryOpToCode("--", &p.UnaryOp)
}

func resolveStringOrCodeElement(v interface{}) string {
	switch rv := v.(type) {
	case []interface{}:
		return resolveArrayInterface(rv, ", ")
	case []string:
		return strings.Join(rv, ", ")
	case string:
		return rv
	case *Literal:
		return resolveLiteral(rv.Value)
	case []*CodeElement:
		rvCodes := make([]string, len(rv))
		for i, v := range rv {
			rvCodes[i] = v.ToCode()
		}
		return strings.Join(rvCodes, ", ")
	case *CodeElement:
		return rv.ToCode()
	default:
		return fmt.Sprintf("%v", rv)
	}
}

// Implementation of ToCode for each struct
func (a *Assignment) ToCode() string {
	leftSide := resolveStringOrArray(a.Left)
	rightSide := resolveStringOrCodeElement(a.Right)
	return fmt.Sprintf("%s = %s", leftSide, rightSide)
}

func (na *NewAssignment) ToCode() string {
	leftSide := resolveStringOrArray(na.Left)
	rightSide := resolveStringOrCodeElement(na.Right)
	return fmt.Sprintf("%s := %s", leftSide, rightSide)
}

func (ie *IfElement) ToCode() string {
	var thenPart, elsePart string
	for _, then := range ie.Then {
		thenPart += then.ToCode() + "\n"
	}
	if len(ie.Else) > 0 {
		elsePart = "else {\n"
		for _, el := range ie.Else {
			elsePart += el.ToCode() + "\n"
		}
		elsePart += "}"
	}
	return fmt.Sprintf("if %s {\n%s}%s", ie.Condition.ToCode(), thenPart, elsePart)
}

func (ce CaseElement) ToCode() string {
	code := ""
	for i, caseElem := range ce {
		condCode := resolveStringOrCodeElement(caseElem.Condition)
		bodyCode := resolveStringOrCodeElement(caseElem.Body)
		if i == 0 {
			code = fmt.Sprintf("if %s {\n%s\n}", condCode, bodyCode)
		} else if condCode == "" {
			code = fmt.Sprintf("else {\n%s\n}", bodyCode)
		} else {
			code = fmt.Sprintf("else if %s {\n%s\n}", condCode, bodyCode)
		}
	}
	return code
}

func (mc MatchCases) ToCode() string {
	code := ""
	for _, caseElem := range mc.MatchCases {
		condCode := resolveStringOrCodeElement(caseElem.MatchWith)
		bodyCode := resolveStringOrCodeElement(caseElem.Body)
		if condCode == "" {
			code = fmt.Sprintf("\ndefault:\n%s%s", Indent, bodyCode)
		} else {
			code = fmt.Sprintf("case %s:\n%s%s", condCode, Indent, bodyCode)
		}
	}

	return code
}

func bodyCodeGen(body []*CodeElement) string {
	var bodyStrings []string
	for _, b := range body {
		bCode := b.ToCode()
		if bCode == "" {
			continue
		}
		bodyStrings = append(bodyStrings, fmt.Sprintf("%s%s", Indent, bCode))
	}
	return strings.Join(bodyStrings, "\n")
}

// ToCode for RepeatByCondition: Repeats based on a condition (like a 'while' loop in other languages)
func (r *RepeatByCondition) ToCode() string {
	bodyCode := bodyCodeGen(r.Body)
	return fmt.Sprintf("for %s {\n%s\n}", r.Condition.ToCode(), bodyCode)
}

// ToCode for RepeatInitConditionElement: Similar to 'for' loop with an init and condition but no increment step
func (ric *RepeatInitConditionElement) ToCode() string {
	var initStrings []string
	for _, i := range ric.Init {
		initStrings = append(initStrings, i.ToCode())
	}

	bodyCode := bodyCodeGen(ric.Body)
	return fmt.Sprintf("for %s; %s; {\n%s\n}", strings.Join(initStrings, ", "), ric.Condition.ToCode(), bodyCode)
}

// ToCode for RepeatLoopElement: Full 'for' loop with init, condition, and step
func (rl *RepeatLoopElement) ToCode() string {
	var initStrings, stepStrings []string
	for _, i := range rl.Init {
		initStrings = append(initStrings, i.ToCode())
	}
	for _, s := range rl.Step {
		stepStrings = append(stepStrings, s.ToCode())
	}

	bodyCode := bodyCodeGen(rl.Body)
	return fmt.Sprintf("for %s; %s; %s {\n%s\n}", strings.Join(initStrings, ", "), rl.Condition.ToCode(), strings.Join(stepStrings, ", "), bodyCode)
}

// ToCode for IterateElement: 'for' loop for iterating over slices, arrays, or maps
func (it *IterateElement) ToCode() string {
	bodyCode := bodyCodeGen(it.Body)
	return fmt.Sprintf("for %s := range %s {\n%s\n}", strings.Join(it.Variables, ", "), it.RangeOn.ToCode(), bodyCode)
}

func (sc *StructCreation) ToCode() string {
	var params []string
	for _, v := range sc.Values {
		params = append(params, v.ToCode())
	}
	return fmt.Sprintf("new(%s){%s}", sc.StructType, strings.Join(params, ", "))
}

func (gr *GoRoutine) ToCode() string {
	return fmt.Sprintf("go %s()", gr.FunctionCall.ToCode())
}

func resolveStringOrArray(s interface{}) string {
	base.LOG.Info("resolveStringOrArray with []string", "s", s, "type", fmt.Sprintf("%T", s))

	switch v := s.(type) {
	case []string:
		return strings.Join(v, ", ")
	case []interface{}:
		return resolveArrayInterface(v, ", ")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func resolveOutputs(output, newOutput interface{}) string {
	if output == nil && newOutput == nil {
		return ""
	}

	if output != nil {
		return fmt.Sprintf("%s = ", resolveStringOrArray(output))
	}

	return fmt.Sprintf("%s := ", resolveStringOrArray(newOutput))
}

func (fc *FunctionCall) ToCode() string {
	leftSide := resolveOutputs(fc.Output, fc.NewOutput)
	var params []string
	for _, p := range fc.Params {
		params = append(params, p.ToCode())
	}
	return fmt.Sprintf("%s%s(%s)", leftSide, fc.Function, strings.Join(params, ", "))
}

func (mfc *MemberFunctionCall) ToCode() string {
	leftSide := resolveOutputs(mfc.Output, mfc.NewOutput)
	var params []string
	for _, p := range mfc.Params {
		params = append(params, p.ToCode())
	}
	return fmt.Sprintf("%s%s.%s(%s)", leftSide, mfc.Receiver, mfc.Function, strings.Join(params, ", "))
}

func convertMapToStruct[T CodeBlock](m map[string]interface{}) T {
	var t T
	ptr := reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
	b, err := yaml.Marshal(m)
	yaml.Unmarshal(b, ptr)

	base.LOG.Debug("convertMapToStruct", "m", m, "b", string(b), "err", err,
		"ptr", ptr, "ttype", fmt.Sprintf("%T", t), "mtype", fmt.Sprintf("%T", m))
	return ptr
}

func genericBinaryOpToCode[T CodeBlock](i interface{}, op string) string {
	base.LOG.Debug("genericBinaryOpToCode", "i", i, "op", op, "type", fmt.Sprintf("%T", i))
	switch v := i.(type) {
	case []string:
		return strings.Join(v, fmt.Sprintf(" %s ", op))
	case []interface{}:
		return resolveArrayInterface(v, fmt.Sprintf(" %s ", op))
	case map[string]interface{}:
		st := convertMapToStruct[T](v)
		return st.ToCode()
	case T:
		return v.ToCode()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func genericUnaryOpToCode[T CodeBlock](i interface{}, op string) string {
	switch v := i.(type) {
	case string:
		return fmt.Sprintf("%s%s", op, v)
	case T:
		return v.ToCode()
	case map[string]interface{}:
		st := convertMapToStruct[T](v)
		return st.ToCode()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func genericUnaryPostOpToCode[T CodeBlock](i interface{}, op string) string {
	switch v := i.(type) {
	case string:
		return fmt.Sprintf("%s%s", v, op)
	case T:
		return v.ToCode()
	case map[string]interface{}:
		return fmt.Sprintf("%s(%v)", op, v["input"])
	default:
		return fmt.Sprintf("%v", v)
	}
}

func AddToCode(a interface{}) string {
	return genericBinaryOpToCode[*Add](a, "+")
}

func SubtractToCode(a interface{}) string {
	return genericBinaryOpToCode[*Sub](a, "-")
}

func MultiplyToCode(a interface{}) string {
	return genericBinaryOpToCode[*Mul](a, "*")
}

func DivideToCode(a interface{}) string {
	return genericBinaryOpToCode[*Div](a, "/")
}

func ModuloToCode(a interface{}) string {
	return genericBinaryOpToCode[*Mod](a, "%")
}

func AndToCode(a interface{}) string {
	return genericBinaryOpToCode[*And](a, "&&")
}

func OrToCode(a interface{}) string {
	return genericBinaryOpToCode[*Or](a, "||")
}

func NotToCode(a interface{}) string {
	return genericUnaryOpToCode[*Not](a, "!")
}

func EqualToCode(a interface{}) string {
	return genericBinaryOpToCode[*Equal](a, "==")
}

func NotEqualToCode(a interface{}) string {
	return genericBinaryOpToCode[*Equal](a, "!=")
}

func GreaterThanCode(a interface{}) string {
	return genericBinaryOpToCode[*GreaterThan](a, ">")
}

func GreaterThanOrEqualToCode(a interface{}) string {
	return genericBinaryOpToCode[*GreaterThan](a, ">=")
}

func LessThanCode(a interface{}) string {
	return genericBinaryOpToCode[*LessThan](a, "<")
}

func LessThanOrEqualToCode(a interface{}) string {
	return genericBinaryOpToCode[*LessThan](a, "<=")
}

func LeftShiftToCode(a interface{}) string {
	return genericBinaryOpToCode[*LeftShift](a, "<<")
}

func RightShiftToCode(a interface{}) string {
	return genericBinaryOpToCode[*RightShift](a, ">>")
}

func BitwiseAndToCode(a interface{}) string {
	return genericBinaryOpToCode[*BitwiseAnd](a, "&")
}

func BitwiseOrToCode(a interface{}) string {
	return genericBinaryOpToCode[*BitwiseOr](a, "|")
}

func BitwiseXorToCode(a interface{}) string {
	return genericBinaryOpToCode[*BitwiseOr](a, "^")
}

func BitwiseNotToCode(a interface{}) string {
	return genericUnaryOpToCode[*BitwiseNot](a, "^")
}

func PostIncrementToCode(a interface{}) string {
	return genericUnaryPostOpToCode[*PostIncrement](a, "++")
}

func PostDecrementToCode(a interface{}) string {
	return genericUnaryPostOpToCode[*PostDecrement](a, "--")
}

func PreIncrementToCode(a interface{}) string {
	return genericUnaryOpToCode[*PreIncrement](a, "++")
}

func PreDecrementToCode(a interface{}) string {
	return genericUnaryOpToCode[*PreDecrement](a, "--")
}

func AssignToCode(a *Assignment) string {
	return a.ToCode()
}

func NewAssignToCode(a *NewAssignment) string {
	return a.ToCode()
}

func ReturnToCode(a interface{}) string {
	return fmt.Sprintf("return %v", resolveStringOrCodeElement(a))
}

// Define templates for each construct
var tpl *template.Template

func init() {
	tpl = template.Must(template.New("code").Funcs(template.FuncMap{
		"add":       AddToCode,
		"sub":       SubtractToCode,
		"mul":       MultiplyToCode,
		"div":       DivideToCode,
		"mod":       ModuloToCode,
		"and":       AndToCode,
		"or":        OrToCode,
		"not":       NotToCode,
		"eq":        EqualToCode,
		"ne":        NotEqualToCode,
		"gt":        GreaterThanCode,
		"ge":        GreaterThanOrEqualToCode,
		"lt":        LessThanCode,
		"le":        LessThanOrEqualToCode,
		"ls":        LeftShiftToCode,
		"rs":        RightShiftToCode,
		"band":      BitwiseAndToCode,
		"bor":       BitwiseOrToCode,
		"bxor":      BitwiseXorToCode,
		"bnot":      BitwiseNotToCode,
		"postinc":   PostIncrementToCode,
		"postdec":   PostDecrementToCode,
		"preinc":    PreIncrementToCode,
		"predec":    PreDecrementToCode,
		"assign":    AssignToCode,
		"newassign": NewAssignToCode,
		"to_code":   CodeElementToCode,
		"return":    ReturnToCode,
	}).Parse(`
{{define "Arithmetic"}}
{{if .Add}}{{add .Add}}{{end}}
{{if .Subtract}}{{sub .Subtract}}{{end}}
{{if .Multiply}}{{mul .Multiply}}{{end}}
{{if .Divide}}{{div .Divide}}{{end}}
{{if .Modulo}}{{mod .Modulo}}{{end}}
{{end}}
{{define "Logical"}}
{{if .And}}{{and .And}}{{end}}
{{if .Or}}{{or .Or}}{{end}}
{{if .Not}}{{not .Not}}{{end}}
{{end}}
{{define "Compare"}}
{{if .Equal}}{{eq .Equal}}{{end}}
{{if .NotEqual}}{{ne .NotEqual}}{{end}}
{{if .GreaterThan}}{{gt .GreaterThan}}{{end}}
{{if .GreaterThanOrEqual}}{{ge .GreaterThanOrEqual}}{{end}}
{{if .LessThan}}{{lt .LessThan}}{{end}}
{{if .LessThanOrEqual}}{{le .LessThanOrEqual}}{{end}}
{{end}}
{{define "Bitwise"}}
{{if .LeftShift}}{{ls .LeftShift}}{{end}}
{{if .RightShift}}{{rs .RightShift}}{{end}}
{{if .BitwiseAnd}}{{band .BitwiseAnd}}{{end}}
{{if .BitwiseOr}}{{bor .BitwiseOr}}{{end}}
{{if .BitwiseXor}}{{bxor .BitwiseXor}}{{end}}
{{if .BitwiseNot}}{{bnot .BitwiseNot}}{{end}}
{{end}}
{{define "Unary"}}
{{if .PostIncrement}}{{postinc .PostIncrement}}{{end}}
{{if .PostDecrement}}{{postdec .PostDecrement}}{{end}}
{{if .PreIncrement}}{{preinc .PreIncrement}}{{end}}
{{if .PreDecrement}}{{predec .PreDecrement}}{{end}}
{{end}}
{{define "Assign"}}{{assign .Assign}}{{end}}
{{define "NewAssign"}}{{newassign .NewAssign}}{{end}}

{{define "If"}}
if {{template "code" .Condition}} {
    {{range .Then}}{{template "code" .}}
    {{end}}
}{{if .Else}} else {
    {{range .Else}}{{template "code" .}}
    {{end}}
}{{end}}
{{end}}

{{define "code"}}
{{template "Arithmetic" .}}
{{template "Logical" .}}
{{template "Compare" .}}
{{template "Bitwise" .}}
{{template "Unary" .}}
{{if .Assign}}{{assign .Assign}}{{end}}
{{if .NewAssign}}{{newassign .NewAssign}}{{end}}
{{if .If}}{{.If.ToCode}}{{end}}
{{if .Cases}}{{.Cases.ToCode}}{{end}}
{{if .MatchCases}}{{.MatchCases.ToCode}}{{end}}
{{if .RepeatCond}}{{.RepeatCond.ToCode}}{{end}}
{{if .RepeatInitCond}}{{.RepeatInitCond.ToCode}}{{end}}
{{if .RepeatLoop}}{{.RepeatLoop.ToCode}}{{end}}
{{if .Iterate}}{{.Iterate.ToCode}}{{end}}
{{if .StructCreation}}{{.StructCreation.ToCode}}{{end}}
{{if .GoRoutine}}{{.GoRoutine.ToCode}}{{end}}
{{if .FunctionCall}}{{.FunctionCall.ToCode}}{{end}}
{{if .MemberFunctionCall}}{{.MemberFunctionCall.ToCode}}{{end}}
{{if .Return}}{{return .Return}}{{end}}
{{if .Steps}}{{range $index, $step := .Steps}}
{{to_code $step}}{{end}}{{end}}
{{end}}
`))
}

// ToCode generates Go code for a *CodeElement using templates
func (ce *CodeElement) ToCode() string {
	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, "code", ce); err != nil {
		base.LOG.Error("Template execution error: %v", err)
		panic(err)
	}
	return strings.Trim(buf.String(), "\n\t ")
}

func CodeElementToCode(ce *CodeElement) string {
	return ce.ToCode()
}
