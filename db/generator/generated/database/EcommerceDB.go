package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

var EcommerceDb *ecommerceDb

type ecommerceDb struct {
	User    *User_DB
	Product *Product_DB
	Order   *Order_DB
}

func InitEcommerceDb() error {
	db, err := SetupDBConnection()
	if err != nil {
		return err
	}
	stmtMapUser, err := UserPrepareStmts(db)
	if err != nil {
		return err
	}
	user := &User_DB{
		db:            db,
		preparedCache: stmtMapUser,
	}
	stmtMapProduct, err := ProductPrepareStmts(db)
	if err != nil {
		return err
	}
	product := &Product_DB{
		db:            db,
		preparedCache: stmtMapProduct,
	}
	stmtMapOrder, err := OrderPrepareStmts(db)
	if err != nil {
		return err
	}
	order := &Order_DB{
		db:            db,
		preparedCache: stmtMapOrder,
	}
	EcommerceDb = &ecommerceDb{
		User:    user,
		Product: product,
		Order:   order,
	}
	return nil
}
func SetupDBConnection() (*sql.DB, error) {
	driverName, dsn := "postgres", "user=user password=password dbname=ecommerce port=5432 host=localhost"
	idleConnTimeout, connMaxLifetime := (time.Second * 10), (time.Minute * 30)
	idleConns, maxOpenConns := 5, 10
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(idleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxIdleTime(idleConnTimeout)
	db.SetConnMaxLifetime(connMaxLifetime)
	return db, nil
}
