package database

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type Order struct {
	OrderDate     *time.Time `db:"order_date"`
	OrderStatus   string     `db:"order_status"`
	PaymentMethod string     `db:"payment_method"`
	TotalAmount   float64    `db:"total_amount"`
	Rating        float64    `db:"rating"`
}

type Order_DB struct {
	db            *sql.DB
	preparedCache map[string]*sql.Stmt
}

type GetOrderByIDParams struct {
	OrderDate interface{} `json:"order_date"`
}

type GetOrderByIDRequest struct {
	Params GetOrderByIDParams `json:"params"`
}

func NewOrder_DB(db *sql.DB, preparedCache map[string]*sql.Stmt) *Order_DB {
	return &Order_DB{
		db:            db,
		preparedCache: preparedCache,
	}
}
func GetOrderByIDReadParams(params GetOrderByIDParams) ([]interface{}, error) {
	var values []interface{}
	values = append(values, params.OrderDate)
	return values, nil
}
func GetOrderByID(ctx context.Context, db *Order_DB, requestParams GetOrderByIDParams) ([]Order, error) {
	stmt := db.preparedCache["GetOrderByID"]
	values, err := GetOrderByIDReadParams(requestParams)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []Order
	for rows.Next() {
		var item Order
		scanErr := rows.Scan(&item.OrderDate, &item.OrderStatus, &item.PaymentMethod, &item.TotalAmount)
		if scanErr != nil {
			return nil, scanErr
		}
		results = append(results, item)
	}
	return results, nil
}
func OrderPrepareStmts(db *sql.DB) (map[string]*sql.Stmt, error) {
	preparedCache := make(map[string]*sql.Stmt)
	var err error
	preparedCache["GetOrderByID"], err = db.Prepare("SELECT order_date, order_status, payment_method, total_amount FROM order WHERE (1 = 1) AND (order_date = $1)")
	if err != nil {
		return nil, err
	}
	return preparedCache, nil
}
