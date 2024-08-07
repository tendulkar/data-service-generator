package database

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

type Product struct {
	Sku         string  `db:"sku"`
	ProductName string  `db:"product_name"`
	Description string  `db:"description"`
	Price       float64 `db:"price"`
}

type Product_DB struct {
	db            *sql.DB
	preparedCache map[string]*sql.Stmt
}

type GetProductByIDParams struct {
	Sku interface{} `json:"sku"`
}

type GetProductByIDRequest struct {
	Params GetProductByIDParams `json:"params"`
}

func NewProduct_DB(db *sql.DB, preparedCache map[string]*sql.Stmt) *Product_DB {
	return &Product_DB{
		db:            db,
		preparedCache: preparedCache,
	}
}
func GetProductByIDReadParams(params GetProductByIDParams) ([]interface{}, error) {
	var values []interface{}
	values = append(values, params.Sku)
	return values, nil
}
func GetProductByID(ctx context.Context, db *Product_DB, requestParams GetProductByIDParams) ([]Product, error) {
	stmt := db.preparedCache["GetProductByID"]
	values, err := GetProductByIDReadParams(requestParams)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []Product
	for rows.Next() {
		var item Product
		scanErr := rows.Scan(&item.Id, &item.Sku, &item.Price)
		if scanErr != nil {
			return nil, scanErr
		}
		results = append(results, item)
	}
	return results, nil
}
func ProductPrepareStmts(db *sql.DB) (map[string]*sql.Stmt, error) {
	preparedCache := make(map[string]*sql.Stmt)
	var err error
	preparedCache["GetProductByID"], err = db.Prepare("SELECT id, sku, price FROM product WHERE (1 = 1) AND (sku = $1)")
	if err != nil {
		return nil, err
	}
	return preparedCache, nil
}
