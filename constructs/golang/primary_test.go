package golang

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCodeElement_ToCode(t *testing.T) {
	// Test case: Arithmetic operation
	ce := &CodeElement{
		Add: []string{"a", "b"},
	}
	expected := "a + b"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: Logical operation
	ce = &CodeElement{
		And: []string{"true", "false"},
	}
	expected = "true && false"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: Compare operation
	ce = &CodeElement{
		Equal: []string{"a", "b"},
	}
	expected = "a == b"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: Bitwise operation
	ce = &CodeElement{
		BitwiseAnd: []string{"a", "b"},
	}
	expected = "a & b"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: Assignment
	ce = &CodeElement{
		Assign: &Assignment{
			Left:  "a",
			Right: "b",
		},
	}
	expected = "a = b"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: New Assignment
	ce = &CodeElement{
		NewAssign: &NewAssignment{
			Left:  []string{"a", "b"},
			Right: []string{"1", "2"},
		},
	}
	expected = "a, b := 1, 2"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: If statement
	ce = &CodeElement{
		If: &IfElement{
			Condition: &CodeElement{GreaterThan: []string{"a", "b"}},
			Then: []*CodeElement{{Assign: &Assignment{
				Left:  "a",
				Right: &CodeElement{Add: []string{"a", "b"}},
			}}},
		},
	}
	expected = "if a > b {\na = a + b\n}"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: For statement
	ce = &CodeElement{
		RepeatCond: &RepeatByCondition{
			Condition: &CodeElement{GreaterThan: []string{"a", "b"}},
			Body: []*CodeElement{{Assign: &Assignment{
				Left:  "a",
				Right: &CodeElement{Add: []string{"a", "b"}},
			}}},
		},
	}
	expected = "for a > b {\na = a + b\n}"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}
}

func TestAddToCodeYaml(t *testing.T) {
	// Test case: Arithmetic operation
	addYaml := `left: a
right: b
out: c`

	add := &Add{}
	yaml.Unmarshal([]byte(addYaml), add)
	t.Log(add)
	expected := "c = a + b"
	if result := add.ToCode(); result != expected {
		t.Errorf("AddToCode() = %v, want %v", result, expected)
	}

	expectedFib := `a, b, c := 0, 1, 0
for i := 0; i < 10; i++ {
	c = a + b
	a, b = b, c
}
return c`

	fibYaml := `
steps:
- new_assign:
    left: [a, b, c]
    right: [0, 1, 0]
- repeat:
    init: 
        - {new_assign: {left: i, right: 0}}
    cond:  {lt: [i, 10]}
    step: 
        - post_inc: i
    body:
        - add:
            out: c
            left: a
            right: b
        - assign:
            left: [a, b]
            right: [b, c]
- return: [c]`

	fibCE := &CodeElement{}
	err := yaml.Unmarshal([]byte(fibYaml), fibCE)
	if err != nil {
		t.Error(err)
	}
	t.Log(fibCE.Steps[0].Assign)
	t.Log(fibCE.Steps[1].RepeatLoop.Init, fibCE.Steps[1].RepeatLoop.Condition,
		fibCE.Steps[1].RepeatLoop.Step, fibCE.Steps[1].RepeatLoop.Body)
	t.Log(fibCE.Steps[2].Return)
	if result := fibCE.ToCode(); result != expectedFib {
		t.Errorf("FibCode() = %v, want %v", result, expectedFib)
	}

}

func TestAddToCode2(t *testing.T) {
	// Test case: Logical operation
	addYaml := `left: a
right: b
out: c`

	add := &Add{}
	yaml.Unmarshal([]byte(addYaml), add)
	t.Log(add)
	expected := "c = a + b"
	if result := AddToCode(add); result != expected {
		t.Errorf("AddToCode() = %v, want %v", result, expected)
	}

	addStrings := []string{"a", "b"}
	expected = "a + b"
	if result := AddToCode(addStrings); result != expected {
		t.Errorf("AddToCode() = %v, want %v", result, expected)
	}

}
