package golang

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"stellarsky.ai/platform/codegen/data-service-generator/base"
)

// CodeElement struct handles all operations and control structures
type CodeElement struct {
	Add                []string                    `yaml:"add,omitempty"`
	Subtract           []string                    `yaml:"sub,omitempty"`
	Multiply           []string                    `yaml:"mul,omitempty"`
	Divide             []string                    `yaml:"div,omitempty"`
	Modulo             []string                    `yaml:"mod,omitempty"`
	And                []string                    `yaml:"and,omitempty"`
	Or                 []string                    `yaml:"or,omitempty"`
	Not                string                      `yaml:"not,omitempty"`
	Equal              []string                    `yaml:"eq,omitempty"`
	NotEqual           []string                    `yaml:"ne,omitempty"`
	LessThan           []string                    `yaml:"lt,omitempty"`
	LessThanOrEqual    []string                    `yaml:"le,omitempty"`
	GreaterThan        []string                    `yaml:"gt,omitempty"`
	GreaterThanOrEqual []string                    `yaml:"ge,omitempty"`
	LeftShift          []string                    `yaml:"lshift,omitempty"`
	RightShift         []string                    `yaml:"rshift,omitempty"`
	BitwiseAnd         []string                    `yaml:"bit_and,omitempty"`
	BitwiseOr          []string                    `yaml:"bit_or,omitempty"`
	BitwiseXor         []string                    `yaml:"bit_xor,omitempty"`
	BitwiseNot         string                      `yaml:"bit_not,omitempty"`
	Assign             *Assignment                 `yaml:"assign,omitempty"`
	NewAssign          *NewAssignment              `yaml:"new_assign,omitempty"`
	If                 *IfElement                  `yaml:"if,omitempty"`
	RepeatCond         *RepeatConditionElement     `yaml:"repeat_cond,omitempty"`
	RepeatInitCond     *RepeatInitConditionElement `yaml:"repeat_init_cond,omitempty"`
	RepeatLoop         *RepeatLoopElement          `yaml:"repeat,omitempty"`
	Iterate            *IterateElement             `yaml:"iterate,omitempty"`
	Return             []string                    `yaml:"return,omitempty"`
	StructCreation     *StructCreation             `yaml:"create,omitempty"`
	GoRoutine          *GoRoutine                  `yaml:"async,omitempty"`
	FunctionCall       *FunctionCall               `yaml:"call,omitempty"`
	MemberFunctionCall *MemberFunctionCall         `yaml:"obj_call,omitempty"`
	Steps              []*CodeElement              `yaml:"steps,omitempty"`
}

// Supporting structs
type Assignment struct {
	Variables     []string       `yaml:"variables"`
	Values        []string       `yaml:"values"`
	ValueElements []*CodeElement `yaml:"value_elements"`
}

type NewAssignment struct {
	Variables     []string       `yaml:"variables"`
	Values        []string       `yaml:"values"`
	ValueElements []*CodeElement `yaml:"value_elements"`
}

type IfElement struct {
	Condition *CodeElement   `yaml:"condition"`
	Then      []*CodeElement `yaml:"then"`
	Else      []*CodeElement `yaml:"else,omitempty"`
}

type RepeatConditionElement struct {
	Condition *CodeElement `yaml:"condition"`
	Body      []*CodeElement
}

type RepeatInitConditionElement struct {
	Init      []*CodeElement
	Condition *CodeElement   `yaml:"condition"`
	Body      []*CodeElement `yaml:"body"`
}

type RepeatLoopElement struct {
	Init      []*CodeElement `yaml:"init"`
	Condition *CodeElement   `yaml:"condition"`
	Step      []*CodeElement `yaml:"step"`
	Body      []*CodeElement
}

type IterateElement struct {
	Variables []string       `yaml:"variables"`
	RangeOn   *CodeElement   `yaml:"range_on"`
	Body      []*CodeElement `yaml:"body"`
}

type StructCreation struct {
	StructType string         `yaml:"struct_type"`
	Values     []*CodeElement `yaml:"values"`
}

type GoRoutine struct {
	FunctionCall       *CodeElement `yaml:"call"`
	MemberFunctionCall *CodeElement `yaml:"obj_call"`
}

type FunctionCall struct {
	Function string         `yaml:"function"`
	Params   []*CodeElement `yaml:"params"`
}

type MemberFunctionCall struct {
	Receiver string         `yaml:"receiver"`
	Function string         `yaml:"function"`
	Params   []*CodeElement `yaml:"params"`
}

// Implementation of ToCode for each struct
func (a *Assignment) ToCode() string {
	var values []string
	if len(a.Values) > 0 {
		values = append(values, a.Values...)
	} else {
		for _, v := range a.ValueElements {
			values = append(values, v.ToCode())
		}
	}
	return fmt.Sprintf("%s = %s", strings.Join(a.Variables, ", "), strings.Join(values, ", "))
}

func (na *NewAssignment) ToCode() string {
	var values []string
	if len(na.Values) > 0 {
		values = append(values, na.Values...)
	} else {
		for _, v := range na.ValueElements {
			values = append(values, v.ToCode())
		}
	}
	return fmt.Sprintf("var %s = %s", strings.Join(na.Variables, ", "), strings.Join(values, ", "))
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

// ToCode for RepeatConditionElement: Repeats based on a condition (like a 'while' loop in other languages)
func (r *RepeatConditionElement) ToCode() string {
	var bodyStrings []string
	for _, b := range r.Body {
		bodyStrings = append(bodyStrings, b.ToCode())
	}
	return fmt.Sprintf("for %s {\n%s\n}", r.Condition.ToCode(), strings.Join(bodyStrings, "\n"))
}

// ToCode for RepeatInitConditionElement: Similar to 'for' loop with an init and condition but no increment step
func (ric *RepeatInitConditionElement) ToCode() string {
	var initStrings, bodyStrings []string
	for _, i := range ric.Init {
		initStrings = append(initStrings, i.ToCode())
	}
	for _, b := range ric.Body {
		bodyStrings = append(bodyStrings, b.ToCode())
	}
	return fmt.Sprintf("for %s; %s; {\n%s\n}", strings.Join(initStrings, ", "), ric.Condition.ToCode(), strings.Join(bodyStrings, "\n"))
}

// ToCode for RepeatLoopElement: Full 'for' loop with init, condition, and step
func (rl *RepeatLoopElement) ToCode() string {
	var initStrings, stepStrings, bodyStrings []string
	for _, i := range rl.Init {
		initStrings = append(initStrings, i.ToCode())
	}
	for _, s := range rl.Step {
		stepStrings = append(stepStrings, s.ToCode())
	}
	for _, b := range rl.Body {
		bodyStrings = append(bodyStrings, b.ToCode())
	}
	return fmt.Sprintf("for %s; %s; %s {\n%s\n}", strings.Join(initStrings, ", "), rl.Condition.ToCode(), strings.Join(stepStrings, ", "), strings.Join(bodyStrings, "\n"))
}

// ToCode for IterateElement: 'for' loop for iterating over slices, arrays, or maps
func (it *IterateElement) ToCode() string {
	var bodyStrings []string
	for _, b := range it.Body {
		bodyStrings = append(bodyStrings, b.ToCode())
	}
	return fmt.Sprintf("for %s := range %s {\n%s\n}", strings.Join(it.Variables, ", "), it.RangeOn.ToCode(), strings.Join(bodyStrings, "\n"))
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

func (fc *FunctionCall) ToCode() string {
	var params []string
	for _, p := range fc.Params {
		params = append(params, p.ToCode())
	}
	return fmt.Sprintf("%s(%s)", fc.Function, strings.Join(params, ", "))
}

func (mfc *MemberFunctionCall) ToCode() string {
	var params []string
	for _, p := range mfc.Params {
		params = append(params, p.ToCode())
	}
	return fmt.Sprintf("%s.%s(%s)", mfc.Receiver, mfc.Function, strings.Join(params, ", "))
}

// Define templates for each construct
var tpl *template.Template

func init() {
	tpl = template.Must(template.New("code").Parse(`
{{define "Arithmetic"}}
{{if .Add}}{{index .Add 0}} + {{index .Add 1}}{{end}}
{{if .Subtract}}{{index .Subtract 0}} - {{index .Subtract 1}}{{end}}
{{if .Multiply}}{{index .Multiply 0}} * {{index .Multiply 1}}{{end}}
{{if .Divide}}{{index .Divide 0}} / {{index .Divide 1}}{{end}}
{{if .Modulo}}{{index .Modulo 0}} % {{index .Modulo 1}}{{end}}
{{end}}
{{define "Logical"}}
{{if .And}}{{index .And 0}} && {{index .And 1}}{{end}}
{{if .Or}}{{index .Or 0}} || {{index .Or 1}}{{end}}
{{if .Not}}!{{index .Not 0}}{{end}}
{{end}}
{{define "Compare"}}
{{if .Equal}}{{index .Equal 0}} == {{index .Equal 1}}{{end}}
{{if .NotEqual}}{{index .NotEqual 0}} != {{index .NotEqual 1}}{{end}}
{{if .GreaterThan}}{{index .GreaterThan 0}} > {{index .GreaterThan 1}}{{end}}
{{if .GreaterThanOrEqual}}{{index .GreaterThanOrEqual 0}} >= {{index .GreaterThanOrEqual 1}}{{end}}
{{if .LessThan}}{{index .LessThan 0}} < {{index .LessThan 1}}{{end}}
{{if .LessThanOrEqual}}{{index .LessThanOrEqual 0}} <= {{index .LessThanOrEqual 1}}{{end}}
{{end}}
{{define "Bitwise"}}
{{if .BitwiseAnd}}{{index .BitwiseAnd 0}} & {{index .BitwiseAnd 1}}{{end}}
{{if .BitwiseOr}}{{index .BitwiseOr 0}} | {{index .BitwiseOr 1}}{{end}}
{{if .BitwiseXor}}{{index .BitwiseXor 0}} ^ {{index .BitwiseXor 1}}{{end}}
{{if .BitwiseNot}}^{{index .BitwiseNot 0}}{{end}}
{{end}}
{{define "Assign"}}
{{range $index, $variable := .Variables}}{{if $index}}, {{end}}{{$variable}}{{end}} = {{range $index, $value := .Values}}{{if $index}}, {{end}}{{template "code" $value}}{{end}}
{{end}}
{{define "NewAssign"}}
{{range $index, $variable := .Variables}}{{if $index}}, {{end}}{{$variable}}{{end}} := {{range $index, $value := .Values}}{{if $index}}, {{end}}{{template "code" $value}}{{end}}
{{end}}

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
{{if .Assign}}{{.Assign.ToCode}}{{end}}
{{if .NewAssign}}{{.NewAssign.ToCode}}{{end}}
{{if .If}}{{.If.ToCode}}{{end}}
{{if .RepeatCond}}{{.RepeatCond.ToCode}}{{end}}
{{if .RepeatInitCond}}{{.RepeatInitCond.ToCode}}{{end}}
{{if .RepeatLoop}}{{.RepeatLoop.ToCode}}{{end}}
{{if .Iterate}}{{.Iterate.ToCode}}{{end}}
{{if .StructCreation}}{{.StructCreation.ToCode}}{{end}}
{{if .GoRoutine}}{{.GoRoutine.ToCode}}{{end}}
{{if .FunctionCall}}{{.FunctionCall.ToCode}}{{end}}
{{if .MemberFunctionCall}}{{.MemberFunctionCall.ToCode}}{{end}}
{{if .Return}}return {{range $index, $elem := .Return}} {{if $index}}, {{end}} {{$elem}} {{end}}{{end}}
{{if .Steps}}{{range $index, $step := .Steps}}{{template "code" $step}}{{end}}{{end}}
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
