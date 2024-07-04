package golang

import (
	"testing"
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
			Variables: []string{"a"},
			Values:    []string{"b"},
		},
	}
	expected = "a = b"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: New Assignment
	ce = &CodeElement{
		NewAssign: &NewAssignment{
			Variables: []string{"a", "b"},
			Values:    []string{"1", "2"},
		},
	}
	expected = "var a, b = 1, 2"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: If statement
	// ce = &CodeElement{
	// 	If: &IfElement{
	// 		Condition: &CodeElement{GreaterThan: []string{"a", "b"}},
	// 		Then: []*CodeElement{Assign: &Assignment{
	// 			Variables: []string{"a"},
	// 			ValueElements: []*CodeElement{&CodeElement{
	// 				Add: []string{"a", "b"},
	// 			},
	// 			}}},
	// 	},
	// }
	// expected = "if a > b {\na = a + b\n}"
	// if result := ce.ToCode(); result != expected {
	// 	t.Errorf("ToCode() = %v, want %v", result, expected)
	// }

}
