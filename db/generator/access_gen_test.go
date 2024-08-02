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

	queries, functions, structs, err := GenerateFindConfigs("Product", "Product_DB", findConfigs)
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
	goSrc, _, err := Generate(cfg)
	t.Log(goSrc.SourceCode())

	assert.NoError(t, err)
	assert.NotNil(t, goSrc)
}

func TestGenerate_ErrorCase(t *testing.T) {
	// Create a sample ModelConfig that triggers an error scenario
	config := defs.ModelConfig{
		// Populate the necessary fields for testing the error case
	}

	goSrc, _, err := Generate(config)

	assert.Nil(t, err)
	assert.NotNil(t, goSrc)
	// You can add more specific assertions based on the expected error behavior
}

func TestGenerateFamily(t *testing.T) {

	modelNameMaps := modelNameMappings{
		{ModelName: "Product", ModelStructName: "Product", ModelDBStructName: "Product_DB"},
		{ModelName: "User", ModelStructName: "User", ModelDBStructName: "User_DB"},
		{ModelName: "Order", ModelStructName: "Order", ModelDBStructName: "Order_DB"},
		{ModelName: "OrderItem", ModelStructName: "OrderItem", ModelDBStructName: "OrderItem_DB"},
		{ModelName: "UserCart", ModelStructName: "UserCart", ModelDBStructName: "UserCart_DB"},
	}

	familyName := "EcommerceDB"

	stDef, fnDef, vars, err := GenerateFamily(familyName, modelNameMaps)
	assert.NoError(t, err)
	assert.NotNil(t, stDef)
	assert.NotNil(t, fnDef)
	assert.NotNil(t, vars)

	expectedSrcCode := `package database

var EcommerceDb *ecommerceDb

type ecommerceDb struct {
	Product   Product
	User      User
	Order     Order
	OrderItem OrderItem
	UserCart  UserCart
}

func InitEcommerceDb() error {
	db, err = SetupDBConnection()
	if err != nil {
		return err
	}
	stmtMapProduct, err = ProductPrepareStmt()
	if err != nil {
		return err
	}
	product := &Product_DB{
		db:            db,
		preparedCache: stmtMapProduct,
	}
	stmtMapUser, err = UserPrepareStmt()
	if err != nil {
		return err
	}
	user := &User_DB{
		db:            db,
		preparedCache: stmtMapUser,
	}
	stmtMapOrder, err = OrderPrepareStmt()
	if err != nil {
		return err
	}
	order := &Order_DB{
		db:            db,
		preparedCache: stmtMapOrder,
	}
	stmtMapOrderItem, err = OrderItemPrepareStmt()
	if err != nil {
		return err
	}
	orderItem := &OrderItem_DB{
		db:            db,
		preparedCache: stmtMapOrderItem,
	}
	stmtMapUserCart, err = UserCartPrepareStmt()
	if err != nil {
		return err
	}
	userCart := &UserCart_DB{
		db:            db,
		preparedCache: stmtMapUserCart,
	}
	EcommerceDb = &ecommerceDb{
		Product:   product,
		User:      user,
		Order:     order,
		OrderItem: orderItem,
		UserCart:  userCart,
	}
	return nil
}
`
	srcFile := golang.GoSourceFile{
		Package:   "database",
		Structs:   stDef,
		Functions: fnDef,
		Variables: vars,
	}
	srcCode, _, err := srcFile.SourceCode()
	assert.NoError(t, err)
	t.Log(srcFile.SourceCode())
	assert.Equal(t, expectedSrcCode, srcCode)
}

func TestGenerateDB(t *testing.T) {
	config.LoadConfig()

	dataConfig := &defs.DataConfig{
		FamilyName:     "EcommerceDB",
		DatabaseConfig: &defs.DatabaseConfig{}, // Just to avoid error returned by GenerateDB function
		Models: []defs.ModelConfig{
			{
				Model: defs.Model{
					Name:       "User",
					Attributes: []int64{2000007, 2000008, 2000009, 2000010},
				},
				Access: defs.Access{
					Find: []defs.AccessConfig{
						{
							Name:       "GetUserByID",
							Attributes: []string{"name", "email", "shopping_address", "billing_address"},
							Filter: []defs.Filter{{
								Attribute: "email",
								Operator:  "=",
								ParamName: "email",
							}},
						},
					},
				},
			},
			{
				Model: defs.Model{
					Name:       "Product",
					Attributes: []int64{2000001, 2000002, 2000003, 2000004},
				},
				Access: defs.Access{
					Find: []defs.AccessConfig{
						{
							Name:       "GetProductByID",
							Attributes: []string{"id", "sku", "price"},
							Filter: []defs.Filter{{
								Attribute: "sku",
								Operator:  "=",
								ParamName: "sku",
							}},
						},
					},
				},
			},
			{
				Model: defs.Model{
					Name:       "Order",
					Attributes: []int64{2000012, 2000013, 2000014, 2000015, 2000016},
				},
				Access: defs.Access{
					Find: []defs.AccessConfig{
						{
							Name:       "GetOrderByID",
							Attributes: []string{"order_date", "order_status", "payment_method", "total_amount"},
							Filter: []defs.Filter{{
								Attribute: "order_date",
								Operator:  "=",
								ParamName: "order_date",
							}},
						},
					},
				},
			},
		},
	}

	unitModules, err := GenerateDB(dataConfig)
	assert.Nil(t, err)
	assert.NotNil(t, unitModules)
	assert.Equal(t, 4, len(unitModules))
	t.Log(unitModules)

	for _, unitModule := range unitModules {
		t.Log(*unitModule)
		t.Log(unitModule.GenerateCode("database"))
	}

	// You can add more specific assertions based on the expected behavior of the GenerateDB function
}
