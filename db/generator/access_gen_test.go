package generator

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
