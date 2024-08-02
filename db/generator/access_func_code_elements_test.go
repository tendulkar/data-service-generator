package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"stellarsky.ai/platform/codegen/data-service-generator/db/generator/defs"
)

func TestFindCodeFunction(t *testing.T) {
	modelName := "User"
	name := "FindUser"
	attributes := []string{"id", "name"}

	expectedFnCode := `func FindUser(ctx context.Context, db *User_DB, requestParams FindUserParams) (results []User, err error) {
	stmt := db.preparedCache["FindUser"]
	values, err := FindUserParseParams(requestParams)
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
		scanErr := rows.Scan(&item.id, &item.name)
		if scanErr != nil {
			return nil, scanErr
		}
		results = append(results, item)
	}
	return results, nil
}`

	expectedImports := map[string]bool{}
	fn := FindCodeFunction(modelName, "User_DB", name, attributes)
	fnCode, fnImports := fn.FunctionCode()
	assert.Equal(t, expectedFnCode, fnCode)
	assert.Equal(t, expectedImports, fnImports)
	assert.Equal(t, name, fn.Name)
	assert.Equal(t, 3, len(fn.Parameters)) // Assuming ctx, db, requestParams
	assert.Equal(t, 6, len(fn.Body))       // Assuming 6 code elements in the body
	assert.Equal(t, 2, len(fn.Returns))
	assert.Equal(t, 0, len(fn.Dependencies))
}

func TestUpdateCodeFunction(t *testing.T) {
	name := "UpdateUser"

	expectedFnCode := `func UpdateUser(ctx context.Context, db *User_DB, requestParams UpdateUserParams) (int64, error) {
	stmt := db.preparedCache["UpdateUser"]
	values, err := UpdateUserParseParams(requestParams)
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
}`

	expectedImports := map[string]bool{}
	fn := UpdateCodeFunction("UpdateUser", "User_DB")
	fnCode, fnImports := fn.FunctionCode()
	assert.Equal(t, expectedFnCode, fnCode)
	assert.Equal(t, expectedImports, fnImports)
	assert.Equal(t, name, fn.Name)
	assert.Equal(t, 3, len(fn.Parameters)) // Assuming ctx, db, requestParams
	assert.Equal(t, 5, len(fn.Body))       // Assuming 6 code elements in the body
	assert.Equal(t, 2, len(fn.Returns))
	assert.Equal(t, 0, len(fn.Dependencies))
}

func TestAddCodeFunction(t *testing.T) {
	name := "AddUser"

	expectedFnCode := `func AddUser(ctx context.Context, db *User_DB, requestParams AddUserParams) (int64, error) {
	stmt := db.preparedCache["AddUser"]
	values, err := AddUserParseParams(requestParams)
	if err != nil {
		return int64(0), err
	}
	var id int64
	queryErr := stmt.QueryRow(values...).Scan(&id)
	if queryErr != nil {
		return int64(0), queryErr
	}
	return id, nil
}`

	expectedImports := map[string]bool{}
	fn := AddCodeFunction(name, "User_DB")
	fnCode, fnImports := fn.FunctionCode()
	// t.Log(fnCode)
	assert.Equal(t, expectedFnCode, fnCode)
	assert.Equal(t, expectedImports, fnImports)
	assert.Equal(t, name, fn.Name)
	assert.Equal(t, 3, len(fn.Parameters)) // Assuming ctx, db, requestParams
	assert.Equal(t, 5, len(fn.Body))       // Assuming 5 code elements in the body
	assert.Equal(t, 2, len(fn.Returns))
	assert.Equal(t, 0, len(fn.Dependencies))
}

func TestAddOrReplaceCodeFunction(t *testing.T) {
	name := "AddOrReplaceUser"

	expectedFnCode := `func AddOrReplaceUser(ctx context.Context, db *User_DB, requestParams AddOrReplaceUserParams) (int64, bool, error) {
	stmt := db.preparedCache["AddOrReplaceUser"]
	values, err := AddOrReplaceUserParseParams(requestParams)
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
}`

	expectedImports := map[string]bool{}

	fn := AddOrReplaceCodeFunction(name, "User_DB")
	fnCode, fnImports := fn.FunctionCode()
	assert.Equal(t, expectedFnCode, fnCode)
	assert.Equal(t, expectedImports, fnImports)
	assert.Equal(t, name, fn.Name)
	assert.Equal(t, 3, len(fn.Parameters)) // Assuming ctx, db, requestParams
	assert.Equal(t, 6, len(fn.Body))       // Assuming 4 code elements in the body
	assert.Equal(t, 3, len(fn.Returns))
	assert.Equal(t, 0, len(fn.Dependencies))
}

func TestDeleteCodeFunction(t *testing.T) {
	name := "DeleteUser"

	expectedFnCode := `func DeleteUser(ctx context.Context, db *User_DB, requestParams DeleteUserParams) (int64, error) {
	stmt := db.preparedCache["DeleteUser"]
	values, err := DeleteUserParseParams(requestParams)
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
}`

	expectedImports := map[string]bool{}
	fn := DeleteCodeFunction(name, "User_DB")
	fnCode, fnImports := fn.FunctionCode()
	assert.Equal(t, expectedFnCode, fnCode)
	assert.Equal(t, expectedImports, fnImports)
	assert.Equal(t, name, fn.Name)
	assert.Equal(t, 3, len(fn.Parameters)) // Assuming ctx, db, param
	assert.Equal(t, 5, len(fn.Body))       // Assuming 5 code elements in the body
	assert.Equal(t, 2, len(fn.Returns))
	assert.Equal(t, 0, len(fn.Dependencies))
}

func TestReadParamsFunction(t *testing.T) {
	paramRefs := []defs.ParameterRef{
		{Name: "age", Index: -1},
		{Name: "salary_range", Index: 0},
		{Name: "salary_range", Index: 1},
		{Name: "department", Index: -1},
	}
	confName := "config"
	valuesName := "values"
	paramsName := "params"

	expectedFnCode := `func configReadParams(params configParams) ([]interface{}, error) {
	var values []interface{}
	values = append(values, params.age)
	values = append(values, params.salary_range.([]interface{})[0])
	values = append(values, params.salary_range.([]interface{})[1])
	values = append(values, params.department)
	return values, nil
}`
	// Call the function under test
	resultFunction := ReadParamsFunction(paramRefs, confName, valuesName, paramsName)
	resultCode, _ := resultFunction.FunctionCode()
	t.Log(resultFunction.FunctionCode())
	assert.Equal(t, expectedFnCode, resultCode)
	// Assert the expected output
	// assert.Equal(t, resultFunction, expectedFunction)
}

func TestPrepareStmtFunction(t *testing.T) {
	queries := []NamedQuery{
		{"query1", "SELECT * FROM table1 WHERE id = $1 AND name = $2"},
		{"query2", "INSERT INTO table2 (id, name, age) VALUES ($1, $2, $3)"},
		{"query3", "UPDATE table3 SET name = $1, age = $2 WHERE id = $3"},
		{"query4", "DELETE FROM table4 WHERE id = $1"},
		{"query5", "SELECT * FROM table5"},
	}

	expectedCode := `func UserPrepareStmt(db *sql.DB, queries map[string]string) (map[string]*sql.Stmt, error) {
	preparedCache := make(map[string]*sql.Stmt)
	var err error
	preparedCache["query1"], err = db.Prepare("SELECT * FROM table1 WHERE id = $1 AND name = $2")
	if err != nil {
		return nil, err
	}
	preparedCache["query2"], err = db.Prepare("INSERT INTO table2 (id, name, age) VALUES ($1, $2, $3)")
	if err != nil {
		return nil, err
	}
	preparedCache["query3"], err = db.Prepare("UPDATE table3 SET name = $1, age = $2 WHERE id = $3")
	if err != nil {
		return nil, err
	}
	preparedCache["query4"], err = db.Prepare("DELETE FROM table4 WHERE id = $1")
	if err != nil {
		return nil, err
	}
	preparedCache["query5"], err = db.Prepare("SELECT * FROM table5")
	if err != nil {
		return nil, err
	}
	return preparedCache, nil
}`

	actual := PrepareStmtFunction("User", queries)
	actualCode, _ := actual.FunctionCode()
	t.Log(actualCode)
	assert.Equal(t, expectedCode, actualCode)
}

func TestSetupDBConnectionFunction(t *testing.T) {
	// Test case 1: Successful setup of the database
	dataConf := defs.DataConfig{
		Models: []defs.ModelConfig{{Model: defs.Model{Name: "Product"}}},
		DatabaseConfig: &defs.DatabaseConfig{
			DriverName: "postgres",
			DBConfigId: "default",
			UserName:   "postgres",
			Password:   "postgres",
			Host:       "localhost",
			Port:       5432,
			DBName:     "postgres",
			ConnectionConfig: &defs.ConnectionConfig{
				MaxLifetimeMins: 30,
			},
			ConnectionPoolConfig: &defs.ConnectionPoolConfig{
				MaxIdleConns: 10,
			},
		},
	}
	expectedImports := []string{"database/sql", "github.com/lib/pq", "time"}
	expectedReturns := typeOnlyParamsCE("*sql.DB", "error")

	fn, err := SetupDBConnectionFunction(dataConf)
	fnCode, fnImports := fn.FunctionCode()

	expectedFnCode := `func SetupDBConnection() (*sql.DB, error) {
	driverName, dsn, idleConns, connMaxLifetime := "postgres", "user=postgres password=postgres dbname=postgres port=5432 host=localhost", 10, 30
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(idleConns)
	db.SetConnMaxLifetime((time.Minute * connMaxLifetime))
	return db, nil
}`

	t.Log(fnCode)

	assert.Nil(t, err)
	assert.Equal(t, expectedFnCode, fnCode)
	assert.Equal(t, "SetupDBConnection", fn.Name)
	assert.Equal(t, expectedImports, fn.Imports)
	assert.Equal(t, expectedReturns, fn.Returns)

	// Test case 2: Check if the correct imports are present
	assert.Equal(t, map[string]bool{
		"database/sql":      true,
		"github.com/lib/pq": true,
		"time":              true,
	}, fnImports)
}

func TestGenerateNewFamilyFunction(t *testing.T) {
	// Create a mock DataConfig and modelNameMappings

	modelNameMaps := modelNameMappings{
		{
			ModelName:         "User",
			ModelStructName:   "User",
			ModelDBStructName: "User_DB",
		},
		{
			ModelName:         "Product",
			ModelStructName:   "Product",
			ModelDBStructName: "Product_DB",
		},
		{
			ModelName:         "Order",
			ModelStructName:   "Order",
			ModelDBStructName: "Order_DB",
		},
	}

	// Call the function
	fn, err := GenerateInitFamilyFunction(modelNameMaps, "RetailDB", "retailDB")
	if err != nil {
		t.Fatalf("GenerateNewFamilyFunctions returned an error: %v", err)
	}
	t.Log(fn.FunctionCode())

	expectedCode := `func InitRetailDB() error {
	db, err = SetupDBConnection()
	if err != nil {
		return err
	}
	stmtMapUser, err = UserPrepareStmt()
	if err != nil {
		return err
	}
	user := &User_DB{
		db: db,
		preparedCache: stmtMapUser,
	}
	stmtMapProduct, err = ProductPrepareStmt()
	if err != nil {
		return err
	}
	product := &Product_DB{
		db: db,
		preparedCache: stmtMapProduct,
	}
	stmtMapOrder, err = OrderPrepareStmt()
	if err != nil {
		return err
	}
	order := &Order_DB{
		db: db,
		preparedCache: stmtMapOrder,
	}
	RetailDB = &retailDB{
		User: user,
		Product: product,
		Order: order,
	}
	return nil
}`

	resultCode, _ := fn.FunctionCode()
	assert.Equal(t, expectedCode, resultCode)
}
