package main

import (
	"fmt"

	"stellarsky.ai/platform/codegen/data-service-generator/base"
	"stellarsky.ai/platform/codegen/data-service-generator/config"
	"stellarsky.ai/platform/codegen/data-service-generator/server"
)

// func main() {
// 	http.HandleFunc("/generate-sql", generateSQLHandler)
// 	fmt.Println("Server is running on http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

func main() {
	cfg := config.LoadConfig()
	base.LOG.Info("Config loaded", "config", cfg)

	attributes := []server.Attribute{
		{ID: 1, Name: "sku", TypeID: 1},
		{ID: 2, Name: "price", TypeID: 3},
	}
	model := server.Model{
		ID:         1,
		Namespace:  "public",
		Family:     "product",
		Name:       "product",
		Attributes: []int{1, 2},
		UniqueConstraints: []struct {
			ConstraintName string `json:"constraint_name"`
			Attributes     []int  `json:"attributes"`
		}{
			{ConstraintName: "sku_unique", Attributes: []int{1}},
		},
	}
	ddl, err := server.GenerateDDL(model, attributes)
	if err != nil {
		fmt.Println("Error generating DDL:", err)
		return
	}
	fmt.Println("DDL SQL:", ddl)
}
