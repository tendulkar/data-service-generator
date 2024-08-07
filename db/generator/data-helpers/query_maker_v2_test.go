package datahelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"stellarsky.ai/platform/codegen/data-service-generator/config"
	"stellarsky.ai/platform/codegen/data-service-generator/db/generator/defs"
)

func TestBuildDeletePreparedStmt(t *testing.T) {
	t.Run("DeleteWithoutFilters", func(t *testing.T) {
		psb := &PreparedStmtBuilder{
			modelName: "test_table",
			dialect:   &PostgresDialect{},
		}
		query, params := psb.BuildDeletePreparedStmt()
		assert.Equal(t, "DELETE FROM `test_table`", query)
		assert.Empty(t, params)
	})

	t.Run("DeleteWithSingleFilter", func(t *testing.T) {
		psb := &PreparedStmtBuilder{
			modelName: "users",
			dialect:   &PostgresDialect{},
			accessConfig: defs.AccessConfig{
				Filter: []defs.Filter{
					{Attribute: "id", Operator: "=", ParamName: "user_id"},
				},
			},
		}
		query, params := psb.BuildDeletePreparedStmt()
		assert.Equal(t, "DELETE FROM `users` WHERE (`id` = $1)", query)
		assert.Equal(t, []defs.ParameterRef{{Name: "user_id", Index: -1}}, params)
	})

	t.Run("DeleteWithMultipleFilters", func(t *testing.T) {
		psb := &PreparedStmtBuilder{
			modelName: "orders",
			dialect:   &PostgresDialect{},
			accessConfig: defs.AccessConfig{
				Filter: []defs.Filter{
					{Attribute: "status", Operator: "=", ParamName: "order_status"},
					{Attribute: "created_at", Operator: "<", ParamName: "date"},
				},
			},
		}
		query, params := psb.BuildDeletePreparedStmt()
		assert.Equal(t, "DELETE FROM `orders` WHERE (`status` = $1 AND `created_at` < $2)", query)
		assert.Equal(t, []defs.ParameterRef{
			{Name: "order_status", Index: -1},
			{Name: "date", Index: -1},
		}, params)
	})

	t.Run("DeleteWithComplexFilters", func(t *testing.T) {
		psb := &PreparedStmtBuilder{
			modelName: "products",
			dialect:   &PostgresDialect{},
			accessConfig: defs.AccessConfig{
				Filter: []defs.Filter{
					{
						Operator: "OR",
						Conditions: []defs.Filter{
							{Attribute: "category", Operator: "=", ParamName: "category"},
							{Attribute: "price", Operator: ">", ParamName: "min_price"},
						},
					},
					{Attribute: "is_active", Operator: "=", ParamName: "active"},
				},
			},
		}
		query, params := psb.BuildDeletePreparedStmt()
		assert.Equal(t, "DELETE FROM `products` WHERE ((`category` = $1 OR `price` > $2) AND `is_active` = $3)", query)
		assert.Equal(t, []defs.ParameterRef{
			{Name: "category", Index: -1},
			{Name: "min_price", Index: -1},
			{Name: "active", Index: -1},
		}, params)
	})
}
func TestBuildFindPreparedStmt(t *testing.T) {
	t.Run("FindAllWithoutFilters", func(t *testing.T) {
		psb := &PreparedStmtBuilder{
			modelName: "users",
			dialect:   &PostgresDialect{},
		}
		query, params := psb.BuildFindPreparedStmt()
		assert.Equal(t, "SELECT * FROM `users`", query)
		assert.Empty(t, params)
	})

	t.Run("FindSpecificAttributesWithoutFilters", func(t *testing.T) {
		psb := &PreparedStmtBuilder{
			modelName: "products",
			dialect:   &PostgresDialect{},
			accessConfig: defs.AccessConfig{
				Attributes: []string{"id", "name", "price"},
			},
		}
		query, params := psb.BuildFindPreparedStmt()
		assert.Equal(t, "SELECT `id`, `name`, `price` FROM `products`", query)
		assert.Empty(t, params)
	})

	t.Run("FindAllWithSingleFilter", func(t *testing.T) {
		psb := &PreparedStmtBuilder{
			modelName: "orders",
			dialect:   &PostgresDialect{},
			accessConfig: defs.AccessConfig{
				Filter: []defs.Filter{
					{Attribute: "status", Operator: "=", ParamName: "order_status"},
				},
			},
		}
		query, params := psb.BuildFindPreparedStmt()
		assert.Equal(t, "SELECT * FROM `orders` WHERE (`status` = $1)", query)
		assert.Equal(t, []defs.ParameterRef{{Name: "order_status", Index: -1}}, params)
	})

	t.Run("FindSpecificAttributesWithMultipleFilters", func(t *testing.T) {
		psb := &PreparedStmtBuilder{
			modelName: "employees",
			dialect:   &PostgresDialect{},
			accessConfig: defs.AccessConfig{
				Attributes: []string{"id", "name", "department"},
				Filter: []defs.Filter{
					{Attribute: "age", Operator: ">", ParamName: "min_age"},
					{Attribute: "salary", Operator: "<", ParamName: "max_salary"},
				},
			},
		}
		query, params := psb.BuildFindPreparedStmt()
		assert.Equal(t, "SELECT `id`, `name`, `department` FROM `employees` WHERE (`age` > $1 AND `salary` < $2)", query)
		assert.Equal(t, []defs.ParameterRef{
			{Name: "min_age", Index: -1},
			{Name: "max_salary", Index: -1},
		}, params)
	})

	t.Run("FindAllWithComplexFilters", func(t *testing.T) {
		psb := &PreparedStmtBuilder{
			modelName: "customers",
			dialect:   &PostgresDialect{},
			accessConfig: defs.AccessConfig{
				Filter: []defs.Filter{
					{
						Operator: "OR",
						Conditions: []defs.Filter{
							{Attribute: "country", Operator: "=", ParamName: "country"},
							{Attribute: "total_purchases", Operator: ">", ParamName: "min_purchases"},
						},
					},
					{Attribute: "is_active", Operator: "=", ParamName: "active"},
				},
			},
		}
		query, params := psb.BuildFindPreparedStmt()
		assert.Equal(t, "SELECT * FROM `customers` WHERE ((`country` = $1 OR `total_purchases` > $2) AND `is_active` = $3)", query)
		assert.Equal(t, []defs.ParameterRef{
			{Name: "country", Index: -1},
			{Name: "min_purchases", Index: -1},
			{Name: "active", Index: -1},
		}, params)
	})

	t.Run("FindWithComplexNestedFilters", func(t *testing.T) {
		psb := &PreparedStmtBuilder{
			modelName: "products",
			dialect:   &PostgresDialect{},
			accessConfig: defs.AccessConfig{
				Attributes: []string{"id", "name", "price", "category"},
				Filter: []defs.Filter{
					{
						Operator: "AND",
						Conditions: []defs.Filter{
							{Attribute: "price", Operator: "BETWEEN", ParamName: "price_range"},
							{
								Operator: "OR",
								Conditions: []defs.Filter{
									{Attribute: "category", Operator: "IN", ParamName: "categories"},
									{Attribute: "name", Operator: "LIKE", ParamName: "name_pattern"},
								},
							},
						},
					},
					{
						Operator: "NOT",
						Conditions: []defs.Filter{
							{Attribute: "is_discontinued", Operator: "=", ParamName: "discontinued"},
						},
					},
				},
			},
		}
		query, params := psb.BuildFindPreparedStmt()
		assert.Equal(t, "SELECT `id`, `name`, `price`, `category` FROM `products` WHERE (((`price` BETWEEN $1 AND $2) AND (`category` = ANY($3) OR `name` LIKE $4)) AND (NOT(`is_discontinued` = $5)))", query)
		assert.Equal(t, []defs.ParameterRef{
			{Name: "price_range", Index: 0},
			{Name: "price_range", Index: 1},
			{Name: "categories", Index: -1},
			{Name: "name_pattern", Index: -1},
			{Name: "discontinued", Index: -1},
		}, params)
	})

}

func TestBuildUpdatePreparedStmt(t *testing.T) {
	t.Run("UpdateWithSetAndFilters", func(t *testing.T) {
		psb := &PreparedStmtBuilder{
			modelName: "users",
			dialect:   &PostgresDialect{},
			accessConfig: defs.AccessConfig{
				Set: []defs.Update{
					{Attribute: "name", ParamName: "new_name"},
					{Attribute: "email", ParamName: "new_email"},
				},
				Filter: []defs.Filter{
					{Attribute: "id", Operator: "=", ParamName: "user_id"},
				},
				Autoincrement:    []string{"version"},
				CaptureTimestamp: []string{"updated_at"},
			},
		}
		query, params := psb.BuildUpdatePreparedStmt()
		assert.Equal(t, "UPDATE `users` SET `name` = $1, `email` = $2, `version` = `version` + 1, `updated_at` = NOW() WHERE (`id` = $3)", query)
		assert.Equal(t, []defs.ParameterRef{
			{Name: "new_name", Index: -1},
			{Name: "new_email", Index: -1},
			{Name: "user_id", Index: -1},
		}, params)
	})
}

func TestBuildAddPreparedStmt(t *testing.T) {
	t.Run("AddNewRecord", func(t *testing.T) {
		psb := &PreparedStmtBuilder{
			modelName: "products",
			dialect:   &PostgresDialect{},
			accessConfig: defs.AccessConfig{
				Values: []defs.Update{
					{Attribute: "name", ParamName: "product_name"},
					{Attribute: "price", ParamName: "product_price"},
				},
			},
		}
		query, params := psb.BuildAddPreparedStmt()
		assert.Equal(t, "INSERT INTO `products` (`id`, `name`, `price`) VALUES ($1, $2, $3) RETURNING id", query)
		assert.Equal(t, []defs.ParameterRef{
			{FuncName: "UUIDV7"},
			{Name: "product_name", Index: -1},
			{Name: "product_price", Index: -1},
		}, params)
	})
}

func TestBuildAddOrReplacePreparedStmt(t *testing.T) {
	t.Run("AddOrReplaceRecord", func(t *testing.T) {
		psb := &PreparedStmtBuilder{
			modelName: "inventory",
			dialect:   &PostgresDialect{},
			accessConfig: defs.AccessConfig{
				Values: []defs.Update{
					{Attribute: "product_id", ParamName: "pid"},
					{Attribute: "quantity", ParamName: "qty"},
				},
			},
		}
		query, params := psb.BuildAddOrReplacePreparedStmt()
		assert.Equal(t,
			"INSERT INTO `inventory` (`id`, `product_id`, `quantity`) VALUES ($1, $2, $3) ON CONFLICT DO UPDATE SET `product_id` = $4, `quantity` = $5 RETURNING id, (xmax = 0) AS inserted",
			query)
		assert.Equal(t, []defs.ParameterRef{
			{FuncName: "UUIDV7"},
			{Name: "pid", Index: -1},
			{Name: "qty", Index: -1},
			{Name: "pid", Index: -1},
			{Name: "qty", Index: -1},
		}, params)
	})
}
func TestBuildCreateTable(t *testing.T) {
	config.LoadConfig()
	t.Run("BasicModelWithoutAttributes", func(t *testing.T) {
		sb := &SchemaBuilder{
			dialect: &PostgresDialect{},
		}
		model := &defs.ModelConfig{
			Model: defs.Model{
				Name: "test_table",
			},
		}
		result := sb.BuildCreateTable(model)
		expected := "CREATE TABLE `test_table` (\n" +
			"	`id` UUID PRIMARY KEY DEFAULT gen_random_uuid(),\n" +
			"	`version` INTEGER NOT NULL DEFAULT 1,\n" +
			"	`updated_at` TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP\n" +
			");\n\n\n"

		assert.Equal(t, expected, result)
	})

	t.Run("ModelWithAttributes", func(t *testing.T) {
		sb := &SchemaBuilder{
			dialect: &PostgresDialect{},
		}
		model := &defs.ModelConfig{
			Model: defs.Model{
				Name:       "products",
				Attributes: []int64{2000001, 2000002, 2000003},
			},
		}
		result := sb.BuildCreateTable(model)
		expected := "CREATE TABLE `products` (\n" +
			"	`id` UUID PRIMARY KEY DEFAULT gen_random_uuid(),\n" +
			"	`version` INTEGER NOT NULL DEFAULT 1,\n" +
			"	`updated_at` TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
			"	`sku` TEXT,\n" +
			"	`product_name` TEXT,\n" +
			"	`description` TEXT\n" +
			");\n\n\n"

		assert.Equal(t, expected, result)
	})

	t.Run("ModelWithAttributesAndIndexes", func(t *testing.T) {
		sb := &SchemaBuilder{
			dialect: &PostgresDialect{},
		}
		model := &defs.ModelConfig{
			Model: defs.Model{
				Name:       "products",
				Attributes: []int64{2000001, 2000002, 2000003},
				Indexes: []struct {
					IndexName  string  `yaml:"index_name"`
					Attributes []int64 `yaml:"attributes"`
				}{
					{
						IndexName:  "idx_product_name",
						Attributes: []int64{2000001},
					},
					{
						IndexName:  "idx_product_category",
						Attributes: []int64{2000002},
					},
				},
			},
		}
		result := sb.BuildCreateTable(model)
		expected := "CREATE TABLE `products` (\n" +
			"	`id` UUID PRIMARY KEY DEFAULT gen_random_uuid(),\n" +
			"	`version` INTEGER NOT NULL DEFAULT 1,\n" +
			"	`updated_at` TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
			"	`sku` TEXT,\n" +
			"	`product_name` TEXT,\n" +
			"	`description` TEXT\n" +
			");\n\n" +
			"CREATE INDEX ON `products` (`product_name`);\n" +
			"CREATE INDEX ON `products` (`sku`);\n"

		assert.Equal(t, expected, result)
	})

	t.Run("ModelWithAttributesAndFilters", func(t *testing.T) {
		sb := &SchemaBuilder{
			dialect: &PostgresDialect{},
		}
		model := &defs.ModelConfig{
			Model: defs.Model{
				Name:       "orders",
				Attributes: []int64{2000001, 2000002, 2000003},
			},
			Access: defs.Access{
				Find: []defs.AccessConfig{{
					Filter: []defs.Filter{
						{Attribute: "sku", Operator: "="},
						{Attribute: "product_name", Operator: "="},
					},
				},
				},
			},
		}
		result := sb.BuildCreateTable(model)
		expected := "CREATE TABLE `orders` (\n" +
			"	`id` UUID PRIMARY KEY DEFAULT gen_random_uuid(),\n" +
			"	`version` INTEGER NOT NULL DEFAULT 1,\n" +
			"	`updated_at` TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
			"	`sku` TEXT,\n" +
			"	`product_name` TEXT,\n" +
			"	`description` TEXT\n" +
			");\n\n" +
			"CREATE INDEX ON `orders` (`product_name`);\n" +
			"CREATE INDEX ON `orders` (`sku`);\n"
		assert.Equal(t, expected, result)
	})
	t.Run("ModelWithAttributesAndFiltersAndIndexes", func(t *testing.T) {
		sb := &SchemaBuilder{
			dialect: &PostgresDialect{},
		}
		model := &defs.ModelConfig{
			Model: defs.Model{
				Name:       "orders",
				Attributes: []int64{2000001, 2000002, 2000003},
				Indexes: []struct {
					IndexName  string  `yaml:"index_name"`
					Attributes []int64 `yaml:"attributes"`
				}{
					{
						IndexName:  "idx_sku_product_name",
						Attributes: []int64{2000001, 2000002},
					},
					{
						IndexName:  "idx_product_description",
						Attributes: []int64{2000002, 2000003},
					},
				},
			},
			Access: defs.Access{
				Find: []defs.AccessConfig{{
					Filter: []defs.Filter{
						{Attribute: "product_name", Operator: "="},
						{Attribute: "sku", Operator: "="},
					},
				},
				},
			},
		}
		result := sb.BuildCreateTable(model)
		expected := "CREATE TABLE `orders` (\n" +
			"	`id` UUID PRIMARY KEY DEFAULT gen_random_uuid(),\n" +
			"	`version` INTEGER NOT NULL DEFAULT 1,\n" +
			"	`updated_at` TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" +
			"	`sku` TEXT,\n" +
			"	`product_name` TEXT,\n" +
			"	`description` TEXT\n" +
			");\n\n" +
			"CREATE INDEX ON `orders` (`product_name`);\n" +
			"CREATE INDEX ON `orders` (`product_name`, `description`);\n" +
			"CREATE INDEX ON `orders` (`sku`);\n" +
			"CREATE INDEX ON `orders` (`sku`, `product_name`);\n"
		assert.Equal(t, expected, result)
	})
}
