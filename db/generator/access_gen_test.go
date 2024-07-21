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

func TestGenerate_Success(t *testing.T) {
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
	goSrc, err := Generate(cfg)
	t.Log(goSrc.SourceCode())

	assert.NoError(t, err)
	assert.NotNil(t, goSrc)
}

func TestGenerate_ErrorCase(t *testing.T) {
	// Create a sample ModelConfig that triggers an error scenario
	config := defs.ModelConfig{
		// Populate the necessary fields for testing the error case
	}

	goSrc, err := Generate(config)

	assert.Error(t, err)
	assert.Nil(t, goSrc)
	// You can add more specific assertions based on the expected error behavior
}
