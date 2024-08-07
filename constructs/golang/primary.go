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
	// All binary or unary operations (Add, Subtract, Multiply ... PreDecrement) can take either array or corresponding binary op
	// For example Add can take array like ["x", "y"] or Add{Left: "x", Right: "y"}
	// For example Subtract can take array like ["x", "y"] or Subtract{Left: "x", Right: "y"}
	// For example GreaterThan can take array like ["x", "y"] or GreaterThan{Left: "x", Right: "y"}
	//
	// All structs can take output or new output, (only one)
	// For example Add{Output: "c", Left: "a", Right: "b"} will generate c = (a + b)
	// For example Add{NewOutput: "c", Left: "a", Right: "b"} will generate c := (a + b)
	Add                interface{} `yaml:"add,omitempty"`
	Subtract           interface{} `yaml:"sub,omitempty"`
	Multiply           interface{} `yaml:"mul,omitempty"`
	Divide             interface{} `yaml:"div,omitempty"`
	Modulo             interface{} `yaml:"mod,omitempty"`
	And                interface{} `yaml:"and,omitempty"`
	Or                 interface{} `yaml:"or,omitempty"`
	Not                interface{} `yaml:"not,omitempty"`
	Equal              interface{} `yaml:"eq,omitempty"`
	NotEqual           interface{} `yaml:"ne,omitempty"`
	LessThan           interface{} `yaml:"lt,omitempty"`
	LessThanOrEqual    interface{} `yaml:"le,omitempty"`
	GreaterThan        interface{} `yaml:"gt,omitempty"`
	GreaterThanOrEqual interface{} `yaml:"ge,omitempty"`
	LeftShift          interface{} `yaml:"lshift,omitempty"`
	RightShift         interface{} `yaml:"rshift,omitempty"`
	BitwiseAnd         interface{} `yaml:"bit_and,omitempty"`
	BitwiseOr          interface{} `yaml:"bit_or,omitempty"`
	BitwiseXor         interface{} `yaml:"bit_xor,omitempty"`
	BitwiseNot         interface{} `yaml:"bit_not,omitempty"`
	PostIncrement      interface{} `yaml:"post_inc,omitempty"`
	PostDecrement      interface{} `yaml:"post_dec,omitempty"`
	PreIncrement       interface{} `yaml:"pre_inc,omitempty"`
	PreDecrement       interface{} `yaml:"pre_dec,omitempty"`

	// Variable can take name to create new variable
	// Example `var int x`
	//
	// Assign can take name and value (don't create new variable)
	// Example Assignment{Left: "x", Right: "y"} will generate `x = y`
	// NewAssign can take name and value (create new variable)
	// Example NewAssignment{Left: "x", Right: "y"} will generate `x := y`
	//
	// Assign and new assign can take complex values as well on right side or left side
	// Example NewAssignment{Left: "x", Right: Add{Left: "y", Right: "z"}} will generate `x := y + z`
	// Example Assignment{Left: "x", Right: Add{Left: "y", Right: "z"}} will generate `x = y + z`
	// Example Assigment{Left: Literal{Value: "x", Attribute: "a"}, Right: Add{Left: "y", Right: "z"}} will generate `x["a"] = y + z`
	// Example Assigment{Left: Literal{Value: "x", Attribute: 1}, Right: Add{Left: "y", Right: "z"}} will generate `x[1] = y + z`
	Variable  *Variable      `yaml:"var,omitempty"`
	Constant  *Constant      `yaml:"const,omitempty"`
	Assign    *Assignment    `yaml:"assign,omitempty"`
	NewAssign *NewAssignment `yaml:"new_assign,omitempty"`

	// Control structures
	If         *IfElement  `yaml:"if,omitempty"`
	Cases      CaseElement `yaml:"cases,omitempty"`
	MatchCases MatchCases  `yaml:"match_cases,omitempty"`

	// Loop structures
	// RepeatCond can take condition to repeat (equivalant of `for condition {body}`)
	// RepeatInitCond can take condition to repeat (equivalant of `for init; condition;  {body}`)
	// RepeatLoop is full for loop, equivalant of `for init; condition; step {body}`
	// RepeatN can take i, start and end. Equivalant of `for i = s; i < n; i++ {body}`
	// Iterate will iterate on array, equivalant of `for i, v := range array {body}`
	RepeatCond     *RepeatByCondition          `yaml:"repeat_cond,omitempty"`
	RepeatInitCond *RepeatInitConditionElement `yaml:"repeat_init_cond,omitempty"`
	RepeatLoop     *RepeatLoopElement          `yaml:"repeat,omitempty"`
	RepeatN        *RepeatNElement             `yaml:"repeat_n,omitempty"`
	Iterate        *IterateElement             `yaml:"iterate,omitempty"`

	// Generate return statement, example `return x, y, z`
	Return interface{} `yaml:"return,omitempty"`

	// Create struct with fields, for example `&Point{x: 1, y: 2}`
	StructCreation *MakeStruct `yaml:"create,omitempty"`

	// Function call, Examples are
	// FunctionCall{Function: "Signal"} will generate `Signal()`
	// FunctionCall{Function: "Sort", Receiver: "x"} will generate `x.Sort()`
	// FunctionCall{Function: "Sort", Receiver: "x", Async: true} will generate `go x.Sort()` (go routine)
	// FunctionCall{Function: "Close", Receiver: "x", Async: true} will generate `defer x.Close()` (defer function call)
	// FunctionCall{Function: "Add", Args: []interface{}{"x", "y"}} will generate `Add(x, y)`
	// FunctionCall{Function: "Add", Args: []interface{}{"x", "y"}, Output: "z"} will generate `z = Add(x, y)`
	// FunctionCall{Function: "Add", Args: []interface{}{"x", "y"}, NewOutput: "z"} will generate `z := Add(x, y)`
	// FunctionCall{Function: "Add", Args: []interface{}{Literal{Value: "x", Index: "a"}, Literal{Value: "y", Index: "b"}}} will generate `Add(x["a"], y["b"])`
	// FunctionCall{Function: "Add", Args: []interface{}{Literal{Value: "x", Index: 1}, Literal{Value: "y", Index: 2}}} will generate `Add(x[1], y[2])`
	// FunctionCall{Function: "Open", Receiver: "io", Args: []string{"file.txt"}, NewOutput: []string{"f", "err"},
	// 				ErrorHandler: &ErrorHandler{ErrorReturns: []string{"nil", "err"}}}
	//		will generate `f, err := io.Open("file.txt"); if err != nil {return nil, err}`
	// FunctionCall{Function: "Open", Receiver: "io", Args: []string{"file.txt"}, NewOutput: []string{"f", "err"},
	// 					ErrorHandler: &ErrorHandler{ErrorReturns: []string{"nil", "err"}},
	// 					CleanningHandler: &CleanningHandler{Receiver: "f", Function: "Close"}}
	// 		will generate `f, err := io.Open("file.txt"); if err != nil {return nil, err}; defer f.Close()`
	FunctionCall     *FunctionCall     `yaml:"call,omitempty"`
	ErrorHandler     *ErrorHandler     `yaml:"on_error,omitempty"`
	CleanningHandler *CleanningHandler `yaml:"clean,omitempty"`

	// Steps often help us to generate more complex code, like body of if, for, function
	Steps []*CodeElement `yaml:"steps,omitempty"`

	// All imports used in the code (could be receiver, attribute)
	// For example "fmt", "database/sql", "github.com/sirupsen/logrus"...
	Imports []string `yaml:"imports,omitempty"`

	// All dependencies used in the code
	// For example "github.com/sirupsen/logrus", version "1.2.3"...
	// No need add builtin imports as dependencies ("database/sql", "fmt", "time"...)
	Dependencies []Dependency `yaml:"dependencies,omitempty"`

	// Literal allows us write code without quotes, and express array index, map index or attribute of struct
	// For example: Literal{Value: "x", Index: "a"} generates x["a"]
	//              Literal{Value: "x", Index: 1} generates x[1]
	//              Literal{Value: "x", Attribute: "y"} generates x.y
	//              Literal{Value: "x", Attribute: Literal{Value: "y", Index: 1}} generates x.y[1]
	// 				Literal{Value: "name"} generates "name" (string with quotes)
	// 				Literal{Value: true} generates true (boolen value, no quotes)
	// 				Literal{Value: 1} generates 1 (int value, no quotes)
	// 				Literal{Value: 1.1} generates 1.1 (float value, no quotes)
	Literal interface{} `yaml:"lit,omitempty"`

	// Other constructs to simplify code generation
	MapLookup *MapLookup     `yaml:"lookup,omitempty"`
	IfError   []*CodeElement `yaml:"if_error,omitempty"`
}

type CodeElements []*CodeElement

func (c CodeElements) ToCode() string {
	var buf bytes.Buffer
	for i, v := range c {
		buf.WriteString(v.ToCode())
		if i < len(c)-1 {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

func (c CodeElements) Imports() []string {
	var imports []string
	for _, v := range c {
		imports = append(imports, v.Imports...)
	}
	return imports
}

func (c CodeElements) Dependencies() []Dependency {
	var deps []Dependency
	for _, v := range c {
		deps = append(deps, v.Dependencies...)
	}
	return deps
}

type BinaryOp struct {
	Left      interface{} `yaml:"left"`
	Right     interface{} `yaml:"right"`
	Output    interface{} `yaml:"out,omitempty"`
	NewOutput interface{} `yaml:"nout,omitempty"`
}

type UnaryOp struct {
	Input     interface{} `yaml:"in"`
	Output    interface{} `yaml:"out,omitempty"`
	NewOutput interface{} `yaml:"nout,omitempty"`
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

type Variable struct {
	// Names can take single string or array of strings
	Names       interface{} `yaml:"name"`
	Type        string      `yaml:"type"`
	ModuleName  string      `yaml:"module,omitempty"`
	IsReference bool        `yaml:"reference,omitempty"`
	// Value can be anything, string or array of strings or code elements
	Values    interface{} `yaml:"val,omitempty"`
	Variables interface{} `yaml:"var,omitempty"`
}

type Constant struct {
	// Names can take single string or array of strings
	Name     string      `yaml:"name"`
	Type     string      `yaml:"type"`
	Value    interface{} `yaml:"val,omitempty"`
	Variable interface{} `yaml:"var,omitempty"`
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
	// Simpilfied construct for if and else
	Condition    interface{}    `yaml:"cond"`
	Then         []*CodeElement `yaml:"then"`
	Break        interface{}    `yaml:"break,omitempty"`
	Continue     interface{}    `yaml:"continue,omitempty"`
	BreakElse    interface{}    `yaml:"break_else,omitempty"`
	ContinueElse interface{}    `yaml:"continue_else,omitempty"`
	Else         []*CodeElement `yaml:"else,omitempty"`
}

type OneCaseElement struct {
	Condition interface{}    `yaml:"cond"`
	Body      []*CodeElement `yaml:"body"`
	Break     interface{}    `yaml:"break,omitempty"`
	Continue  interface{}    `yaml:"continue,omitempty"`
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

type RepeatNElement struct {
	Iterator string         `yaml:"iter"`
	Start    string         `yaml:"start"`
	Limit    string         `yaml:"limit"`
	Body     []*CodeElement `yaml:"body"`
}

type IterateElement struct {
	Variables []string       `yaml:"variables"`
	RangeOn   *CodeElement   `yaml:"range_on"`
	Body      []*CodeElement `yaml:"body"`
}

type KeyValue struct {
	Key      string      `yaml:"key"`
	Value    interface{} `yaml:"value"`
	Variable interface{} `yaml:"var,omitempty"`
}

type KeyValues []*KeyValue

type MakeStruct struct {
	Output      interface{} `yaml:"out,omitempty"`
	NewOutput   interface{} `yaml:"nout,omitempty"`
	ModuleName  string      `yaml:"module_name"`
	StructType  string      `yaml:"struct_type"`
	KeyValues   KeyValues   `yaml:"values"`
	NoReference bool        `yaml:"no_ref,omitempty"`
}

type FunctionCall struct {
	Output           interface{} `yaml:"out,omitempty"`
	NewOutput        interface{} `yaml:"nout,omitempty"`
	Receiver         string      `yaml:"obj,omitempty"`
	Function         string      `yaml:"func"`
	Args             interface{} `yaml:"args,omitempty"`
	Defer            bool        `yaml:"defer,omitempty"`
	Async            bool        `yaml:"async,omitempty"`
	*ErrorHandler    `yaml:",inline"`
	CleanningHandler *CleanningHandler `yaml:"clean,omitempty"`
}

type Literal struct {
	Value     interface{} `yaml:"val"`
	Type      string      `yaml:"type,omitempty"`
	Indexes   interface{} `yaml:"indexes,omitempty"`
	Attribute interface{} `yaml:"attr,omitempty"`
}

type MapLookup struct {
	Output           interface{}      `yaml:"out,omitempty"`
	NewOutput        interface{}      `yaml:"nout,omitempty"`
	Receiver         string           `yaml:"obj,omitempty"`
	Name             string           `yaml:"name,omitempty"`
	Key              interface{}      `yaml:"key,omitempty"`
	CleanningHandler CleanningHandler `yaml:"clean,omitempty"`
}

type ErrorHandler struct {
	Error                string       `yaml:"err,omitempty"`
	ErrorReturns         interface{}  `yaml:"err_returns,omitempty"`
	ErrorFunctionReturns []*Parameter `yaml:"err_func_returns,omitempty"`
	ErrorSteps           CodeElements `yaml:"err_steps,omitempty"`
}

type CleanningHandler struct {
	Receiver string       `yaml:"obj,omitempty"`
	Function string       `yaml:"func,omitempty"`
	Args     interface{}  `yaml:"args,omitempty"`
	Steps    CodeElements `yaml:"steps,omitempty"`
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
	// base.LOG.Info("resolveArrayInterface", "array", arr)
	values := make([]string, len(arr))
	for i, v := range arr {
		values[i] = resolveStringOrCodeElement(v, 1, sep)
	}

	return strings.Join(values, sep)
}

// write a literal using yaml and resolve that don't need quotes
// we don't want to write golag code, either quotes, or brackets for arrays, we just want the literal
// This will be platform compatible
func resolveLiteral(v interface{}) string {
	// base.LOG.Info("resolveLiteral", "value", v)
	switch rv := v.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", rv)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", rv)
	case Literal:
		return rv.ToCode()
	case *Literal:
		return rv.ToCode()
	case map[string]interface{}:
		st := convertMapToStruct[*Literal](rv)
		return st.ToCode()
	case []interface{}:
		return resolveArrayInterface(rv, ", ")

	default:
		return fmt.Sprintf("%v", rv)
	}
}

func binaryOpToCode(op string, b *BinaryOp) string {
	lvalue := resolveStringOrCodeElement(b.Left, 0, ", ")
	rvalue := resolveStringOrCodeElement(b.Right, 0, ", ")

	if b.Output == nil && b.NewOutput == nil {
		return fmt.Sprintf("(%s %s %s)", lvalue, op, rvalue)
	} else if b.Output != nil {
		return fmt.Sprintf("%s = (%s %s %s)", resolveStringOrCodeElement(b.Output, 0, ", "), lvalue, op, rvalue)
	} else {
		return fmt.Sprintf("%s := (%s %s %s)", resolveStringOrCodeElement(b.NewOutput, 0, ", "), lvalue, op, rvalue)
	}
}

func unaryOpToCode(op string, u *UnaryOp) string {
	inpValue := resolveStringOrCodeElement(u.Input, 0, ", ")

	if u.Output == nil && u.NewOutput == nil {
		return fmt.Sprintf("%s %s", op, inpValue)
	} else if u.Output == nil {
		return fmt.Sprintf("%s = %s%s", u.Output, op, inpValue)
	} else {
		return fmt.Sprintf("%s := %s%s", u.NewOutput, op, inpValue)
	}
}

func unaryPostOpToCode(op string, u *UnaryOp) string {
	inpValue := resolveStringOrCodeElement(u.Input, 0, ", ")
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

func resolveStringOrCodeElement(v interface{}, indentSize int, sep string) string {
	// base.LOG.Info("resolveStringOrCodeElement", "value", v, "sep", sep)
	switch rv := v.(type) {
	case []interface{}:
		return resolveArrayInterface(rv, sep)
	case []string:
		return strings.Join(rv, sep)
	case string:
		return rv
	case *Literal:
		return rv.ToCode()
	case []*Literal:
		codeParts := make([]string, len(rv))
		for i, v := range rv {
			codeParts[i] = resolveLiteral(v)
		}
		return strings.Join(codeParts, sep)
	case []*CodeElement:
		return bodyCodeGen(rv, indentSize, sep)
	case *CodeElement:
		return rv.ToCode()
	case CodeBlock:
		return rv.ToCode()
	case map[string]interface{}:
		st := convertMapToStruct[*CodeElement](rv)
		return st.ToCode()
	default:
		if rv == nil {
			return ""
		}
		return fmt.Sprintf("%v", rv)
	}
}

func buidBreakStatement(v interface{}) string {
	if v == nil {
		return ""
	}
	switch rv := v.(type) {
	case bool:
		if rv {
			return "break"
		}
		return ""
	case string:
		// break with label
		return fmt.Sprintf("break %s", rv)
	default:
		return fmt.Sprintf("%v", rv)
	}
}

func buildContinueStatement(v interface{}) string {
	if v == nil {
		return ""
	}
	switch rv := v.(type) {
	case bool:
		if rv {
			return "continue"
		}
		return ""
	case string:
		// continue with label
		return fmt.Sprintf("continue %s", rv)
	default:
		return fmt.Sprintf("%v", rv)
	}
}

func bodyWithBreakAndContinue(body []*CodeElement, b interface{}, c interface{}) string {

	bodyCode := resolveStringOrCodeElement(body, 1, "\n")
	if b == nil && c == nil {
		return bodyCode
	}
	if b != nil {
		breakCode := buidBreakStatement(b)
		return fmt.Sprintf("%s\n%s%s", bodyCode, Indent, breakCode)
	}
	continueCode := buildContinueStatement(c)
	return fmt.Sprintf("%s\n%s%s", bodyCode, Indent, continueCode)
}

func (vc *Variable) ToCode() string {
	typeName := vc.Type
	if vc.ModuleName != "" {
		typeName = fmt.Sprintf("%s.%s", vc.ModuleName, typeName)
	}
	if vc.IsReference {
		typeName = fmt.Sprintf("*%s", typeName)
	}

	valueName := resolveStringOrCodeElement(vc.Values, 0, ", ")
	if valueName != "" {
		valueName = fmt.Sprintf(" = %s", valueName)
	}

	names := resolveStringOrCodeElement(vc.Names, 0, ", ")
	return fmt.Sprintf("var %s %s%s", names, typeName, valueName)
}

func (c *Constant) ToCode() string {
	typeNameWithSpace := ""
	if c.Type != "" {
		typeNameWithSpace = fmt.Sprintf(" %s", c.Type)
	}

	valueName := resolveStringOrCodeElement(c.Value, 0, ", ")
	varName := resolveStringOrCodeElement(c.Name, 0, "")

	return fmt.Sprintf("const %s%s = %s", varName, typeNameWithSpace, valueName)
}

// Implementation of ToCode for each struct
func (a *Assignment) ToCode() string {
	leftSide := resolveStringOrCodeElement(a.Left, 0, ", ")
	rightSide := resolveStringOrCodeElement(a.Right, 0, ", ")
	return fmt.Sprintf("%s = %s", leftSide, rightSide)
}

func (na *NewAssignment) ToCode() string {
	leftSide := resolveStringOrCodeElement(na.Left, 0, ", ")
	rightSide := resolveStringOrCodeElement(na.Right, 0, ", ")
	return fmt.Sprintf("%s := %s", leftSide, rightSide)
}

func (ie *IfElement) ToCode() string {
	condCode := resolveStringOrCodeElement(ie.Condition, 0, " && ")
	thenCode := bodyWithBreakAndContinue(ie.Then, ie.Break, ie.Continue)
	if len(ie.Else) > 0 {
		elseCode := bodyWithBreakAndContinue(ie.Else, ie.BreakElse, ie.ContinueElse)
		return fmt.Sprintf("if %s {\n%s\n} else {\n%s\n}", condCode, thenCode, elseCode)
	}
	return fmt.Sprintf("if %s {\n%s\n}", condCode, thenCode)
}

func (ce CaseElement) ToCode() string {
	code := ""
	for i, caseElem := range ce {
		condCode := resolveStringOrCodeElement(caseElem.Condition, 0, " && ")
		bodyCode := bodyWithBreakAndContinue(caseElem.Body, caseElem.Break, caseElem.Continue)
		base.LOG.Info("Case element", "i", i, "caseElem", caseElem, "condCode", condCode, "bodyCode", bodyCode)
		if i == 0 {
			code = fmt.Sprintf("if %s {\n%s\n}", condCode, bodyCode)
		} else if condCode == "" {
			code = fmt.Sprintf("%s else {\n%s\n}", code, bodyCode)
		} else {
			code = fmt.Sprintf("%s else if %s {\n%s\n}", code, condCode, bodyCode)
		}
	}
	return code
}

func (mc MatchCases) ToCode() string {
	code := ""
	for _, caseElem := range mc.MatchCases {
		condCode := resolveStringOrCodeElement(caseElem.MatchWith, 0, ", ")
		bodyCode := resolveStringOrCodeElement(caseElem.Body, 0, "\n")
		if condCode == "" {
			code = fmt.Sprintf("\ndefault:\n%s%s", Indent, bodyCode)
		} else {
			code = fmt.Sprintf("case %s:\n%s%s", condCode, Indent, bodyCode)
		}
	}

	return code
}

func bodyCodeGen(body []*CodeElement, indentSize int, sep string) string {
	var bodyStrings []string
	for _, b := range body {
		bCode := b.ToCode()
		if bCode == "" {
			continue
		}
		indentedBody := IndentCode(bCode, indentSize)
		bodyStrings = append(bodyStrings, indentedBody)
	}
	return strings.Join(bodyStrings, sep)
}

// ToCode for RepeatByCondition: Repeats based on a condition (like a 'while' loop in other languages)
func (r *RepeatByCondition) ToCode() string {
	bodyCode := bodyCodeGen(r.Body, 1, "\n")
	return fmt.Sprintf("for %s {\n%s\n}", r.Condition.ToCode(), bodyCode)
}

// ToCode for RepeatInitConditionElement: Similar to 'for' loop with an init and condition but no increment step
func (ric *RepeatInitConditionElement) ToCode() string {
	var initStrings []string
	for _, i := range ric.Init {
		initStrings = append(initStrings, i.ToCode())
	}

	bodyCode := bodyCodeGen(ric.Body, 1, "\n")
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

	bodyCode := bodyCodeGen(rl.Body, 1, "\n")
	return fmt.Sprintf("for %s; %s; %s {\n%s\n}", strings.Join(initStrings, ", "), rl.Condition.ToCode(), strings.Join(stepStrings, ", "), bodyCode)
}

func (rn *RepeatNElement) ToCode() string {
	bodyCode := bodyCodeGen(rn.Body, 1, "\n")
	start := "0"
	if rn.Start != "" {
		start = rn.Start
	}
	return fmt.Sprintf("for %s := %s; %s < %s; %s++ {\n%s\n}", rn.Iterator, start, rn.Iterator, rn.Limit, rn.Iterator, bodyCode)
}

// ToCode for IterateElement: 'for' loop for iterating over slices, arrays, or maps
func (it *IterateElement) ToCode() string {
	bodyCode := bodyCodeGen(it.Body, 1, "\n")
	return fmt.Sprintf("for %s := range %s {\n%s\n}", strings.Join(it.Variables, ", "), it.RangeOn.ToCode(), bodyCode)
}

func (kv *KeyValue) ToCode() string {
	value := ""
	if kv.Value != nil {
		value = resolveLiteral(kv.Value)
	} else if kv.Variable != nil {
		value = resolveStringOrCodeElement(kv.Variable, 0, ", ")
	}
	return fmt.Sprintf("%s: %s", kv.Key, value)
}

func (kvs KeyValues) ToCode() string {
	kvCodeParts := make([]string, len(kvs))
	for i, kv := range kvs {
		kvCode := kv.ToCode()
		kvCodeParts[i] = fmt.Sprintf("%s,\n", kvCode)
	}
	return strings.Join(kvCodeParts, "")
}

func (sc *MakeStruct) ToCode() string {
	leftSide := resolveOutputs(sc.Output, sc.NewOutput)
	fieldsCode := IndentCode(sc.KeyValues.ToCode(), 1)
	typeName := sc.StructType
	if sc.ModuleName != "" {
		typeName = fmt.Sprintf("%s.%s", sc.ModuleName, sc.StructType)
	}
	if !sc.NoReference {
		typeName = fmt.Sprintf("&%s", typeName)
	}
	return fmt.Sprintf("%s%s{\n%s}", leftSide, typeName, fieldsCode)
}

func resolveOutputs(output, newOutput interface{}) string {
	if output == nil && newOutput == nil {
		return ""
	}

	// output or newOutput could be literal like array/map index, so resolveStringOrCdoeElement will do that
	if output != nil {
		return fmt.Sprintf("%s = ", resolveStringOrCodeElement(output, 0, ", "))
	}

	return fmt.Sprintf("%s := ", resolveStringOrCodeElement(newOutput, 0, ", "))
}

func findZeroValues(params []*Parameter, errName string) []string {
	var zeroValues []string
	for _, p := range params {
		if p.Type.Name == "error" {
			continue
		}
		zeroValues = append(zeroValues, p.Type.ZeroValue())
	}

	zeroValues = append(zeroValues, errName)
	return zeroValues
}

func (fc *FunctionCall) ToCode() string {
	leftSide := resolveOutputs(fc.Output, fc.NewOutput)
	argsCode := resolveStringOrCodeElement(fc.Args, 0, ", ")
	fnName := fc.Function
	// base.LOG.Info("FunctionCall ToCode", "fc", *fc, "leftSide", leftSide, "params", argsCode)
	if fc.Receiver != "" {
		fnName = fmt.Sprintf("%s.%s", fc.Receiver, fc.Function)
	}

	// if defered or async, just call the function, no assignment nor error handling
	if fc.Defer {
		return fmt.Sprintf("defer %s(%s)", fnName, argsCode)
	} else if fc.Async {
		return fmt.Sprintf("go %s(%s)", fnName, argsCode)
	}

	fnPart := fmt.Sprintf("%s%s(%s)", leftSide, fnName, argsCode)
	fullCode := fnPart
	if fc.ErrorHandler != nil {
		errPart := fc.ErrorHandler.ToCode()
		fullCode = fmt.Sprintf("%s\n%s", fnPart, errPart)
	}
	if fc.CleanningHandler != nil {
		cleanPart := fc.CleanningHandler.ToCode()
		fullCode = fmt.Sprintf("%s\n%s", fullCode, cleanPart)
	}

	return fullCode
}

func (l *Literal) ToCode() string {
	if l.Indexes != nil {
		indicesCode := []string{}
		if reflect.TypeOf(l.Indexes).Kind() == reflect.Slice {
			for _, index := range l.Indexes.([]*Literal) {
				indicesCode = append(indicesCode, fmt.Sprintf("[%s]", resolveLiteral(index)))
			}
		} else {
			indicesCode = append(indicesCode, fmt.Sprintf("[%s]", resolveLiteral(l.Indexes)))
		}

		return fmt.Sprintf("%s%s", resolveStringOrCodeElement(l.Value, 0, ", "), strings.Join(indicesCode, ""))
	}

	if l.Attribute != nil {
		attributes := []string{}
		if reflect.TypeOf(l.Attribute).Kind() == reflect.Slice {
			for _, attr := range l.Attribute.([]*Literal) {
				attributes = append(attributes, resolveLiteral(attr))
			}
		} else {
			if reflect.TypeOf(l.Attribute).Kind() == reflect.String {
				attributes = append(attributes, l.Attribute.(string))
			} else {
				attributes = append(attributes, resolveLiteral(l.Attribute))
			}
		}
		return fmt.Sprintf("%s.%s", resolveStringOrCodeElement(l.Value, 0, ", "), strings.Join(attributes, "."))
	}
	if l.Type == "" {
		return resolveLiteral(l.Value)
	} else {
		return resolveTypeLiteral(l.Value, l.Type)
	}
}

func (ml *MapLookup) ToCode() string {
	name := ml.Name
	if ml.Receiver != "" {
		name = fmt.Sprintf("%s.%s", ml.Receiver, ml.Name)
	}

	key := resolveLiteral(ml.Key)
	leftSide := resolveOutputs(ml.Output, ml.NewOutput)
	return fmt.Sprintf("%s%s[%s]", leftSide, name, key)
}

func (eh *ErrorHandler) ToCode() string {
	errName := eh.Error
	if errName == "" {
		errName = "err"
	}

	errPart := ""
	if eh.ErrorReturns != nil {
		errPart = IndentCode(ReturnToCode(eh.ErrorReturns), 1)
	} else if eh.ErrorSteps != nil {
		errPart = bodyCodeGen(eh.ErrorSteps, 1, "\n")
	} else if eh.ErrorFunctionReturns != nil {
		errPart = IndentCode(ReturnToCode(findZeroValues(eh.ErrorFunctionReturns, errName)), 1)
	} else {
		return ""
	}

	errCode := fmt.Sprintf("if %s != nil {\n%s\n}", errName, errPart)
	return errCode
}

func (ch *CleanningHandler) ToCode() string {

	argsCode := resolveStringOrCodeElement(ch.Args, 0, ", ")
	if ch.Steps == nil && ch.Function == "" {
		return ""
	} else if ch.Steps == nil {
		fnName := ch.Function
		if ch.Receiver != "" {
			fnName = fmt.Sprintf("%s.%s", ch.Receiver, ch.Function)
		}
		return fmt.Sprintf("defer %s(%s)", fnName, argsCode)
	} else {
		bodyCode := bodyCodeGen(ch.Steps, 1, "\n")
		return fmt.Sprintf("defer func(%s) {\n%s\n} ()", bodyCode, argsCode)
	}
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
		binOp := BinaryOp{
			Left:  v[0],
			Right: v[1],
		}
		return binaryOpToCode(op, &binOp)
	case []interface{}:
		binOp := BinaryOp{
			Left:  v[0],
			Right: v[1],
		}
		return binaryOpToCode(op, &binOp)
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
	return fmt.Sprintf("return %v", resolveStringOrCodeElement(a, 0, ", "))
}

func IfErrorToCode(ce CodeElements) string {
	codePart := ce.ToCode()
	indentedCode := IndentCode(codePart, 1)
	return fmt.Sprintf("if err != nil {\n%s\n}", indentedCode)
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
		"ife":       IfErrorToCode,
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
{{if .Variable}}{{.Variable.ToCode}}{{end}}
{{if .Constant}}{{.Constant.ToCode}}{{end}}
{{if .Assign}}{{assign .Assign}}{{end}}
{{if .NewAssign}}{{newassign .NewAssign}}{{end}}
{{if .If}}{{.If.ToCode}}{{end}}
{{if .Cases}}{{.Cases.ToCode}}{{end}}
{{if .MatchCases}}{{.MatchCases.ToCode}}{{end}}
{{if .RepeatCond}}{{.RepeatCond.ToCode}}{{end}}
{{if .RepeatInitCond}}{{.RepeatInitCond.ToCode}}{{end}}
{{if .RepeatLoop}}{{.RepeatLoop.ToCode}}{{end}}
{{if .RepeatN}}{{.RepeatN.ToCode}}{{end}}
{{if .Iterate}}{{.Iterate.ToCode}}{{end}}
{{if .StructCreation}}{{.StructCreation.ToCode}}{{end}}
{{if .FunctionCall}}{{.FunctionCall.ToCode}}{{end}}
{{if .ErrorHandler}}{{.ErrorHandler.ToCode}}{{end}}
{{if .CleanningHandler}}{{.CleanningHandler.ToCode}}{{end}}
{{if .Return}}{{return .Return}}{{end}}
{{if .MapLookup}}{{.MapLookup.ToCode}}{{end}}
{{if .IfError}}{{ife .IfError}}{{end}}
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
	return strings.Trim(buf.String(), "\n")
}

func CodeElementToCode(ce *CodeElement) string {
	return ce.ToCode()
}
