package database

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

type User struct {
	Email           string `db:"email"`
	Name            string `db:"name"`
	ShippingAddress string `db:"shipping_address"`
	BillingAddress  string `db:"billing_address"`
}

type User_DB struct {
	db            *sql.DB
	preparedCache map[string]*sql.Stmt
}

type GetUserByEmailParams struct {
	Email interface{} `json:"email"`
}

type GetUserByEmailRequest struct {
	Params GetUserByEmailParams `json:"params"`
}

type GetUserByNameParams struct {
	Name interface{} `json:"name"`
}

type GetUserByNameRequest struct {
	Params GetUserByNameParams `json:"params"`
}

type GetUserByIDParams struct {
	Id interface{} `json:"id"`
}

type GetUserByIDRequest struct {
	Params GetUserByIDParams `json:"params"`
}

type UpdateUserParams struct {
	Name interface{} `json:"name"`
	Id   interface{} `json:"id"`
}

type UpdateUserRequest struct {
	Params UpdateUserParams `json:"params"`
}

type AddUserParams struct {
	Name            interface{} `json:"name"`
	Email           interface{} `json:"email"`
	ShippingAddress interface{} `json:"shipping_address"`
	BillingAddress  interface{} `json:"billing_address"`
}

type AddUserRequest struct {
	Params AddUserParams `json:"params"`
}

type AddOrReplaceUserParams struct {
	Name            interface{} `json:"name"`
	Email           interface{} `json:"email"`
	ShippingAddress interface{} `json:"shipping_address"`
	BillingAddress  interface{} `json:"billing_address"`
}

type AddOrReplaceUserRequest struct {
	Params AddOrReplaceUserParams `json:"params"`
}

type DeleteUserParams struct {
	Id interface{} `json:"id"`
}

type DeleteUserRequest struct {
	Params DeleteUserParams `json:"params"`
}

func NewUser_DB(db *sql.DB, preparedCache map[string]*sql.Stmt) *User_DB {
	return &User_DB{
		db:            db,
		preparedCache: preparedCache,
	}
}
func GetUserByEmailReadParams(params GetUserByEmailParams) ([]interface{}, error) {
	var values []interface{}
	values = append(values, params.Email)
	return values, nil
}
func GetUserByEmail(ctx context.Context, db *User_DB, requestParams GetUserByEmailParams) ([]User, error) {
	stmt := db.preparedCache["GetUserByEmail"]
	values, err := GetUserByEmailReadParams(requestParams)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []User
	for rows.Next() {
		var item User
		scanErr := rows.Scan(&item.Name, &item.Email, &item.ShippingAddress, &item.BillingAddress)
		if scanErr != nil {
			return nil, scanErr
		}
		results = append(results, item)
	}
	return results, nil
}
func GetUserByNameReadParams(params GetUserByNameParams) ([]interface{}, error) {
	var values []interface{}
	values = append(values, params.Name)
	return values, nil
}
func GetUserByName(ctx context.Context, db *User_DB, requestParams GetUserByNameParams) ([]User, error) {
	stmt := db.preparedCache["GetUserByName"]
	values, err := GetUserByNameReadParams(requestParams)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []User
	for rows.Next() {
		var item User
		scanErr := rows.Scan(&item.Name, &item.Email, &item.ShippingAddress, &item.BillingAddress)
		if scanErr != nil {
			return nil, scanErr
		}
		results = append(results, item)
	}
	return results, nil
}
func GetUserByIDReadParams(params GetUserByIDParams) ([]interface{}, error) {
	var values []interface{}
	values = append(values, params.Id)
	return values, nil
}
func GetUserByID(ctx context.Context, db *User_DB, requestParams GetUserByIDParams) ([]User, error) {
	stmt := db.preparedCache["GetUserByID"]
	values, err := GetUserByIDReadParams(requestParams)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []User
	for rows.Next() {
		var item User
		scanErr := rows.Scan(&item.Name, &item.Email, &item.ShippingAddress, &item.BillingAddress)
		if scanErr != nil {
			return nil, scanErr
		}
		results = append(results, item)
	}
	return results, nil
}
func UpdateUserReadParams(params UpdateUserParams) ([]interface{}, error) {
	var values []interface{}
	values = append(values, params.Name)
	values = append(values, params.Id)
	return values, nil
}
func UpdateUser(ctx context.Context, db *User_DB, requestParams UpdateUserParams) (int64, error) {
	stmt := db.preparedCache["UpdateUser"]
	values, err := UpdateUserReadParams(requestParams)
	if err != nil {
		return int64(0), err
	}
	result, err := stmt.Exec(values...)
	if err != nil {
		return int64(0), err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return int64(0), err
	}
	return rowsAffected, nil
}
func AddUserReadParams(params AddUserParams) ([]interface{}, error) {
	var values []interface{}
	values = append(values, params.Name)
	values = append(values, params.Email)
	values = append(values, params.ShippingAddress)
	values = append(values, params.BillingAddress)
	return values, nil
}
func AddUser(ctx context.Context, db *User_DB, requestParams AddUserParams) (int64, error) {
	stmt := db.preparedCache["AddUser"]
	values, err := AddUserReadParams(requestParams)
	if err != nil {
		return int64(0), err
	}
	var id int64
	queryErr := stmt.QueryRow(values...).Scan(&id)
	if queryErr != nil {
		return int64(0), queryErr
	}
	return id, nil
}
func AddOrReplaceUserReadParams(params AddOrReplaceUserParams) ([]interface{}, error) {
	var values []interface{}
	values = append(values, params.Name)
	values = append(values, params.Email)
	values = append(values, params.ShippingAddress)
	values = append(values, params.BillingAddress)
	values = append(values, params.Name)
	values = append(values, params.Email)
	values = append(values, params.ShippingAddress)
	values = append(values, params.BillingAddress)
	return values, nil
}
func AddOrReplaceUser(ctx context.Context, db *User_DB, requestParams AddOrReplaceUserParams) (int64, bool, error) {
	stmt := db.preparedCache["AddOrReplaceUser"]
	values, err := AddOrReplaceUserReadParams(requestParams)
	if err != nil {
		return int64(0), false, err
	}
	var id int64
	var inserted bool
	queryErr := stmt.QueryRow(values...).Scan(&id, &inserted)
	if queryErr != nil {
		return int64(0), false, queryErr
	}
	return id, inserted, nil
}
func DeleteUserReadParams(params DeleteUserParams) ([]interface{}, error) {
	var values []interface{}
	values = append(values, params.Id)
	return values, nil
}
func DeleteUser(ctx context.Context, db *User_DB, requestParams DeleteUserParams) (int64, error) {
	stmt := db.preparedCache["DeleteUser"]
	values, err := DeleteUserReadParams(requestParams)
	if err != nil {
		return int64(0), err
	}
	result, err := stmt.Exec(values...)
	if err != nil {
		return int64(0), err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return int64(0), err
	}
	return rowsAffected, nil
}
func UserPrepareStmts(db *sql.DB) (map[string]*sql.Stmt, error) {
	preparedCache := make(map[string]*sql.Stmt)
	var err error
	preparedCache["GetUserByEmail"], err = db.Prepare("SELECT name, email, shipping_address, billing_address FROM user WHERE (1 = 1) AND (email = $1)")
	if err != nil {
		return nil, err
	}
	preparedCache["GetUserByName"], err = db.Prepare("SELECT name, email, shipping_address, billing_address FROM user WHERE (1 = 1) AND (name = $1)")
	if err != nil {
		return nil, err
	}
	preparedCache["GetUserByID"], err = db.Prepare("SELECT name, email, shipping_address, billing_address FROM user WHERE (1 = 1) AND (id = ANY($1))")
	if err != nil {
		return nil, err
	}
	preparedCache["UpdateUser"], err = db.Prepare("UPDATE user SET name = $1 WHERE (1 = 1) AND (id = $2)")
	if err != nil {
		return nil, err
	}
	preparedCache["AddUser"], err = db.Prepare("INSERT INTO user (name, email, shipping_address, billing_address) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return nil, err
	}
	preparedCache["AddOrReplaceUser"], err = db.Prepare("INSERT INTO user (name, email, shipping_address, billing_address) VALUES ($1, $2, $3, $4) ON CONFLICT DO UPDATE SET name = $5, email = $6, shipping_address = $7, billing_address = $8 RETURNING id, (xmax = 0)")
	if err != nil {
		return nil, err
	}
	preparedCache["DeleteUser"], err = db.Prepare("DELETE FROM user WHERE (1 = 1) AND (id = $1)")
	if err != nil {
		return nil, err
	}
	return preparedCache, nil
}
