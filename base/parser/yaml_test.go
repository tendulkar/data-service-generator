package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string `yaml:"name"`
	Age  int    `yaml:"age"`
}

func TestReadYamlTo_ValidFile(t *testing.T) {
	// Arrange
	filename := "valid_test.yaml"
	expected := &TestStruct{Name: "test", Age: 30}

	// Act
	result, err := ReadYamlTo[TestStruct](filename)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadYamlTo_InvalidFile(t *testing.T) {
	// Arrange
	filename := "invalid_test.yaml"

	// Act
	result, err := ReadYamlTo[TestStruct](filename)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, result)
}

func TestReadYamlTo_FileNotFound(t *testing.T) {
	// Arrange
	filename := "nonexistent_test.yaml"

	// Act
	result, err := ReadYamlTo[TestStruct](filename)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestReadYamlTo_InvalidYamlFormat(t *testing.T) {
	// Arrange
	filename := "invalid_temp_test.yaml"

	// Create a file with invalid YAML content
	err := os.WriteFile(filename, []byte("name: test\nage: thirty"), 0644)
	assert.NoError(t, err)
	defer os.Remove(filename)

	// Act
	result, err := ReadYamlTo[TestStruct](filename)

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, result)
}

func TestReadYamlTo_EmptyFile(t *testing.T) {
	// Arrange
	filename := "empty_temp_test.yaml"

	// Create an empty file
	err := os.WriteFile(filename, []byte(""), 0644)
	assert.NoError(t, err)
	defer os.Remove(filename)

	// Act
	result, err := ReadYamlTo[TestStruct](filename)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestStringToYaml(t *testing.T) {
	// Arrange
	yamlString := `
name: John Doe
age: 30
address:
  street: 123 Main St
  city: Anytown
  state: CA
`
	type expectedStruct struct {
		Name    string
		Age     int
		Address struct {
			Street string
			City   string
			State  string
		}
	}
	expected := &expectedStruct{
		Name: "John Doe",
		Age:  30,
		Address: struct {
			Street string
			City   string
			State  string
		}{
			Street: "123 Main St",
			City:   "Anytown",
			State:  "CA",
		},
	}

	// Act
	result, err := StringToYaml[expectedStruct](yamlString)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
