package models

import (
	"fmt"
	"strings"
)

type Comparision int

const (
	Equal Comparision = iota
	NotEqual
	Less
	LessOrEqual
	Greater
	GreaterOrEqual
)

func (c Comparision) String() string {
	return [...]string{"=", "!=", "<", "<=", ">", ">="}[c]
}

func (c Comparision) SQL() string {
	return [...]string{"=", "<>", "<", "<=", ">", ">="}[c]
}

type LogicCondition int

const (
	And LogicCondition = iota
	Or
	Not
)

func (c LogicCondition) String() string {
	return [...]string{"AND", "OR", "NOT"}[c]
}

func (c LogicCondition) SQL() string {
	return [...]string{"AND", "OR", "NOT"}[c]
}

type Condition struct {
	Left  string
	Right string
	Op    Comparision
}

func (c Condition) SQL() string {
	return fmt.Sprintf("%s %s %s", c.Left, c.Op.SQL(), c.Right)
}

func (c Condition) String() string {
	return fmt.Sprintf("%s %s %s", c.Left, c.Op.String(), c.Right)
}

type Conditions []Condition

func (c Conditions) SQL() string {
	var sqlConditions []string
	for _, cond := range c {
		sqlConditions = append(sqlConditions, cond.SQL())
	}
	return strings.Join(sqlConditions, " ")
}

func (c Conditions) String() string {
	var sqlConditions []string
	for _, cond := range c {
		sqlConditions = append(sqlConditions, cond.String())
	}
	return strings.Join(sqlConditions, " ")
}

type Filter struct {
	Conditions Conditions
	Logic      LogicCondition
}

type DataOp struct {
}

type ReadOp struct {
	DataOp
	ModelID int64
	Conds   Conditions
}

type AddOp struct {
	DataOp
}

type AddOrReplace struct {
	DataOp
}

type UpdateOp struct {
	DataOp
}

type DeleteOp struct {
	DataOp
}
