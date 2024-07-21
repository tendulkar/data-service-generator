package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"stellarsky.ai/platform/codegen/data-service-generator/config"
	"stellarsky.ai/platform/codegen/data-service-generator/constructs/golang"
	"stellarsky.ai/platform/codegen/data-service-generator/db/generator/defs"
)

func TestGenerateFindConfigs(t *testing.T) {

	findConfigs := []defs.AccessConfig{
		{
			Name:       "FindConfig1",
			Attributes: []string{"attr1", "attr2"},
			Filter: []defs.Filter{
				{Attribute: "attr1", Operator: "=", ParamName: "p1"},
				{Attribute: "attr2", Operator: ">", ParamName: "p2"},
				{Operator: "AND", Conditions: []defs.Filter{
					{Attribute: "attr3", Operator: "=", ParamName: "p3"},
					{Attribute: "attr4", Operator: ">", ParamName: "p4"},
					{Operator: "OR", Conditions: []defs.Filter{
						{Attribute: "attr5", Operator: "=", ParamName: "p5"},
						{Attribute: "attr6", Operator: ">", ParamName: "p6"},
						{Attribute: "attr7", Operator: "BETWEEN", ParamName: "p7"},
						{Attribute: "attr8", Operator: "IN", ParamName: "p8"},
					}},
				}},
			},
		},
	}

	queries, functions, structs, err := GenerateFindConfigs("Product", findConfigs)
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.NotNil(t, functions)
	assert.NotNil(t, structs)

	sourceFile := &golang.GoSourceFile{
		Package:   "database",
		Structs:   structs,
		Functions: functions,
	}
	code, deps, _ := sourceFile.SourceCode()
	t.Log("TestGenerateFindConfigs Code:", code, "Dependencies:", deps, "Imports:", sourceFile.Imports, "queries:", queries)
	// Add more assertions as needed to validate the output
}

func TestGenerateV2_Success(t *testing.T) {
	// Create a sample ModelConfig for testing
	cfg := defs.ModelConfig{
		Model: defs.Model{
			Name:       "Product",
			Attributes: []int64{2000001, 2000002, 2000003, 2000004},
		},
		Access: defs.Access{
			Find: []defs.AccessConfig{
				{
					Name:       "GetProductByID",
					Attributes: []string{"id", "name", "price", "quantity"},
					Filter: []defs.Filter{{
						Attribute: "ID",
						Operator:  "=",
						ParamName: "id",
					},
					},
				},
			},
			Update: []defs.AccessConfig{
				{
					Name: "UpdateProduct",
					Filter: []defs.Filter{{
						Attribute: "id",
						Operator:  "=",
						ParamName: "id",
					},
					},
					Set: []defs.Update{
						{
							Attribute: "name",
							ParamName: "name",
						},
					},
				},
			},
			Add: []defs.AccessConfig{
				{
					Name: "AddProduct",
					Values: []defs.Update{
						{
							Attribute: "Sku",
							ParamName: "sku",
						},
						{
							Attribute: "Price",
							ParamName: "price",
						},
					},
				},
			},
			AddOrReplace: []defs.AccessConfig{
				{
					Name: "AddOrReplaceProduct",
					Values: []defs.Update{
						{
							Attribute: "Sku",
							ParamName: "sku",
						},
						{
							Attribute: "Price",
							ParamName: "price",
						},
					},
				},
			},
			Delete: []defs.AccessConfig{
				{
					Name: "DeleteProduct",
					Filter: []defs.Filter{{
						Attribute: "id",
						Operator:  "=",
						ParamName: "id",
					}},
				},
			},
		},
	}

	config.LoadConfig()
	goSrc, err := GenerateV2(cfg)
	t.Log(goSrc.SourceCode())

	assert.NoError(t, err)
	assert.NotNil(t, goSrc)
}

func TestGenerateV2_ErrorCase(t *testing.T) {
	// Create a sample ModelConfig that triggers an error scenario
	config := defs.ModelConfig{
		// Populate the necessary fields for testing the error case
	}

	goSrc, err := GenerateV2(config)

	assert.Error(t, err)
	assert.Nil(t, goSrc)
	// You can add more specific assertions based on the expected error behavior
}

// func TestApplyTransformation(t *testing.T) {
// 	tests := []struct {
// 		attribute      string
// 		transformation string
// 		expected       string
// 	}{
// 		{"name", "UPPER", "UPPER(name)"},
// 		{"price", "", "price"},
// 		{"created_at", "DATE", "DATE(created_at)"},
// 	}

// 	for _, test := range tests {
// 		result := ApplyTransformation(test.attribute, test.transformation)
// 		if result != test.expected {
// 			t.Errorf("Expected '%s', but got '%s'", test.expected, result)
// 		}
// 	}
// }

// func TestArgs(t *testing.T) {
// 	tests := []struct {
// 		conditions []Filter
// 		expected   string
// 	}{
// 		{
// 			[]Filter{
// 				{Attribute: "sku", Operator: "=", Value: "12345"},
// 				{Attribute: "price", Operator: ">", Value: 50},
// 			},
// 			"sku string, price string",
// 		},
// 		{
// 			[]Filter{
// 				{Operator: "AND", Conditions: []Filter{
// 					{Attribute: "sku", Operator: "=", Value: "12345"},
// 					{Attribute: "price", Operator: ">", Value: 50},
// 				}},
// 			},
// 			"sku string, price string",
// 		},
// 	}

// 	for _, test := range tests {
// 		result := Args(test.conditions)
// 		if result != test.expected {
// 			t.Errorf("Expected '%s', but got '%s'", test.expected, result)
// 		}
// 	}
// }

// func TestWhereClause(t *testing.T) {
// 	tests := []struct {
// 		conditions []Filter
// 		expected   string
// 	}{
// 		{
// 			[]Filter{
// 				{Attribute: "sku", Operator: "=", Value: "12345"},
// 				{Attribute: "price", Operator: ">", Value: 50},
// 			},
// 			"sku = $1 AND price > $2",
// 		},
// 		{
// 			[]Filter{
// 				{Operator: "AND", Conditions: []Filter{
// 					{Attribute: "sku", Operator: "=", Value: "12345"},
// 					{Attribute: "price", Operator: ">", Value: 50},
// 				}},
// 			},
// 			"(sku = $1 AND price > $2)",
// 		},
// 		{
// 			[]Filter{
// 				{Operator: "NOT", Conditions: []Filter{
// 					{Attribute: "terminated", Operator: "=", Value: true},
// 				}},
// 			},
// 			"NOT(terminated = $1)",
// 		},
// 	}

// 	for _, test := range tests {
// 		result := WhereClause(test.conditions)
// 		if result != test.expected {
// 			t.Errorf("Expected '%s', but got '%s'", test.expected, result)
// 		}
// 	}
// }

// func TestAttributeNames(t *testing.T) {
// 	tests := []struct {
// 		updates  []Update
// 		expected string
// 	}{
// 		{
// 			[]Update{
// 				{Attribute: "sku", Value: "67890"},
// 				{Attribute: "name", Value: "New Product"},
// 			},
// 			"sku, name",
// 		},
// 	}

// 	for _, test := range tests {
// 		result := AttributeNames(test.updates)
// 		if result != test.expected {
// 			t.Errorf("Expected '%s', but got '%s'", test.expected, result)
// 		}
// 	}
// }

// func TestAttributeValues(t *testing.T) {
// 	tests := []struct {
// 		updates  []Update
// 		expected string
// 	}{
// 		{
// 			[]Update{
// 				{Attribute: "sku", Value: "67890"},
// 				{Attribute: "name", Value: "New Product"},
// 			},
// 			"$1, $2",
// 		},
// 	}

// 	for _, test := range tests {
// 		result := AttributeValues(test.updates)
// 		if result != test.expected {
// 			t.Errorf("Expected '%s', but got '%s'", test.expected, result)
// 		}
// 	}
// }

// func TestSetClause(t *testing.T) {
// 	tests := []struct {
// 		updates  []Update
// 		expected string
// 	}{
// 		{
// 			[]Update{
// 				{Attribute: "price", Value: 29.99},
// 				{Attribute: "quantity", Value: 50},
// 			},
// 			"price = $1, quantity = $2",
// 		},
// 	}

// 	for _, test := range tests {
// 		result := SetClause(test.updates)
// 		if result != test.expected {
// 			t.Errorf("Expected '%s', but got '%s'", test.expected, result)
// 		}
// 	}
// }

// func TestScanArgs(t *testing.T) {
// 	tests := []struct {
// 		attributes []string
// 		expected   string
// 	}{
// 		{
// 			[]string{"sku", "name"},
// 			"&item.sku, &item.name",
// 		},
// 	}

// 	for _, test := range tests {
// 		result := ScanArgs(test.attributes)
// 		if result != test.expected {
// 			t.Errorf("Expected '%s', but got '%s'", test.expected, result)
// 		}
// 	}
// }

// func TestGenerate(t *testing.T) {
// 	config := ModelConfig{
// 		Model: Model{
// 			ID:        3000001,
// 			Namespace: "e-commerce",
// 			Family:    "inventory",
// 			Name:      "Product",
// 			Attributes: []struct {
// 				ID   int    `yaml:"id"`
// 				Name string `yaml:"name"`
// 				Type string `yaml:"type"`
// 			}{
// 				{ID: 2000001, Name: "sku", Type: "string"},
// 				{ID: 2000002, Name: "name", Type: "string"},
// 				{ID: 2000003, Name: "price", Type: "float"},
// 				{ID: 2000009, Name: "quantity", Type: "int"},
// 			},
// 			UniqueConstraints: []struct {
// 				ConstraintName string `yaml:"constraint_name"`
// 				Attributes     []int  `yaml:"attributes"`
// 			}{
// 				{ConstraintName: "Unique SKU", Attributes: []int{2000001}},
// 			},
// 		},
// 		Access: Access{
// 			Find: []AccessConfig{
// 				{
// 					Name: "FindProductBySku",
// 					Filter: []Filter{
// 						{Attribute: "sku", Operator: "=", Value: "12345"},
// 					},
// 					Attributes: []string{"sku", "name", "price", "quantity"},
// 				},
// 			},
// 			Update: []AccessConfig{
// 				{
// 					Name: "UpdateProductPriceAndQuantityBySku",
// 					Filter: []Filter{
// 						{Attribute: "sku", Operator: "=", Value: "12345"},
// 					},
// 					Set: []Update{
// 						{Attribute: "price", Value: 29.99},
// 						{Attribute: "quantity", Value: 50},
// 					},
// 					Autoincrement:    []string{"version"},
// 					CaptureTimestamp: []string{"last_updated"},
// 				},
// 			},

// 			Add: []AccessConfig{
// 				{
// 					Name: "AddNewProduct",
// 					Values: []Update{
// 						{Attribute: "sku", Value: "67890"},
// 						{Attribute: "name", Value: "New Product"},
// 						{Attribute: "price", Value: 19.99},
// 						{Attribute: "quantity", Value: 100},
// 					},
// 				},
// 			},
// 			AddOrReplace: []AccessConfig{
// 				{
// 					Name: "AddOrReplaceProduct",
// 					Values: []Update{
// 						{Attribute: "sku", Value: "67890"},
// 						{Attribute: "name", Value: "New Product"},
// 						{Attribute: "price", Value: 19.99},
// 						{Attribute: "quantity", Value: 100},
// 					},
// 				},
// 			},
// 			Delete: []AccessConfig{
// 				{
// 					Name: "DeleteProductBySku",
// 					Filter: []Filter{
// 						{Attribute: "sku", Operator: "=", Value: "67890"},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	err := Generate(config)
// 	if err != nil {
// 		t.Fatalf("Error generating code: %v", err)
// 	}

// 	// Check if the file was created
// 	fileName := "Product_gen.go"
// 	if _, err := os.Stat(fileName); os.IsNotExist(err) {
// 		t.Fatalf("Expected file %s to be created, but it does not exist", fileName)
// 	}

// 	// Clean up the generated file
// 	// err = os.Remove(fileName)
// 	// if err != nil {
// 	// 	t.Fatalf("Error cleaning up generated file: %v", err)
// 	// }
// }
