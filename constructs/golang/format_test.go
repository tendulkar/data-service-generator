package golang

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"HelloWorld", "hello_world"},
		{"Hello World", "hello_world"},
		{"helloWorld", "hello_world"},
		{"hello_world", "hello_world"},
		{"HELLO_WORLD", "hello_world"},
		{"helloWorld123", "hello_world_123"},
		{"hello-world", "hello_world"},
		{"Hello__World", "hello__world"},
	}

	for _, name := range cases {
		if ToSnakeCase(name.input) != name.expected {
			t.Errorf("%s -> %s\n", name.input, ToSnakeCase(name.input))
		}
	}
}

func TestToSnakeCaseArray(t *testing.T) {
	cases := []struct {
		input    []string
		expected []string
	}{
		{[]string{"HelloWorld", "Hello World", "helloWorld"}, []string{"hello_world", "hello_world", "hello_world"}},
		{[]string{"hello_world", "HELLO_WORLD", "helloWorld123"}, []string{"hello_world", "hello_world", "hello_world_123"}},
		{[]string{"hello-world", "Hello__World"}, []string{"hello_world", "hello__world"}},
	}

	for _, testCase := range cases {
		result := ToSnakeCaseArray(testCase.input)
		for i, res := range result {
			if res != testCase.expected[i] {
				t.Errorf("ToSnakeCaseArray(%v) = %v, want %v", testCase.input, result, testCase.expected)
			}
		}
	}
}
