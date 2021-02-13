package dbservice

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	//Библиотека PostgresQL
	_ "github.com/lib/pq"
)

//MDB Структура, содержит ссылку на интерфейс к ДБ
type MDB struct {
	Pdb **sql.DB
}

const dbConnectString = "host=localhost port=5432 user=postgres password=Mm000000 dbname=bakery sslmode=disable"
const dbConnectStringInit = "host=localhost port=5432 user=postgres password=Mm000000 dbname=postgres sslmode=disable"

//InitDatabase Проверка наличия БД, создание и обновление до последней версии
func (mdb MDB) InitDatabase() {

	initdb, err := sql.Open("postgres", dbConnectStringInit)

	if err != nil {
		fmt.Println("Database opening error:", err)
		panic("Database error")
	}
	defer initdb.Close()

	rows, err := initdb.Query("SELECT datname FROM pg_database WHERE datistemplate = false AND datname = 'bakery';")

	if err != nil {
		fmt.Println("Error searching database:", err)
		panic("Error searching database")
	}
	if rows.Next() {
		fmt.Println("Database bakery found")
	} else {
		createDb(initdb)
	}
	initdb.Close()
	db, err := sql.Open("postgres", dbConnectString)
	if err != nil {
		fmt.Println("Database opening error:", err)
		panic("Database error")
	}
	updateDb(db)
	mdb.Pdb = &db

}
func createDb(db *sql.DB) {
	_, err := db.Exec("CREATE DATABASE bakery WITH OWNER postgres;")
	if err != nil {
		fmt.Println("Error creating database:", err)
		panic("Error creating database")
	}
	fmt.Println("Database created successfully")

}

func updateDb(db *sql.DB) {

	dbVersion := -1
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'settings'")
	if err != nil {
		fmt.Println("Error searching table:", err)
		panic("Error searching database")
	}
	if rows.Next() {
		fmt.Println("Table settings found")
		dbVersion = getDbVersion(db)
	} else {
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS public.settings(
			name varchar(50) not null,
			string_value varchar(200),
			numeric_value numeric(15,3),
			time_value timestamp,
			constraint name_data primary key (name));`)
		if err != nil {
			fmt.Println("Table create error", err)
			panic("Table create error")
		}
	}
	if dbVersion < 0 {
		updateDbToVerion0(db)
	}

}
func getDbVersion(db *sql.DB) int {
	rows, err := db.Query(`SELECT numeric_value FROM public.settings WHERE name = 'dbVersion'`)
	if err != nil {
		fmt.Println("Error inserting dbVersion data: ", err)
		panic("Error inserting dbVersion data")
	}
	var version int
	version = -1
	if rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			fmt.Println("Error getting db version: ", err)
			panic("Error getting db version")
		}
	}
	return version
}

func setDbVersion(version uint8, db *sql.DB) {
	insertString := ""
	if getDbVersion(db) == -1 {

		insertString = `INSERT INTO public.settings(name, numeric_value) VALUES ($1,$2);`
	} else {
		insertString = `UPDATE public.settings SET numeric_value =$2 WHERE name = $1;`
	}
	_, err := db.Exec(insertString, "dbVersion", version)

	if err != nil {
		fmt.Println("Error inserting dbVersion data: ", err)
		panic("Error inserting dbVersion data")
	}
}

func updateDbToVerion0(db *sql.DB) {
	initTableString := `
	CREATE TABLE public.units(
		id serial,
		name varchar(200) not null,
		constraint id_units primary key (id));
	CREATE TABLE public.materials(
		id serial,
		name varchar(200) not null,
		recipe_unit_id integer not null,
		price_unit_id integer not null,
		coefficient numeric(15,2) not null,
		constraint id_materials primary key (id));
	CREATE TABLE public.material_prices(
		id serial,
		material_id	integer not null,
		date timestamp default current_timestamp,
		price numeric(15,2),
		constraint id_material_prices primary key (id));
	CREATE TABLE public.recipes(
		id serial,
		name varchar(200) not null,
		output numeric(15,2) not null,
		constraint id_recipes primary key (id));
	CREATE TABLE public.recipes_content(
		id integer not null,
		material_id integer not null,
		string_order integer not null,
		qty numeric(15,2) not null);
	CREATE TABLE public.recipe_prices(
		id serial,
		recipe_id integer not null,
		date timestamp default current_timestamp,
		price numeric(15,2),
		constraint id_recipe_prices primary key (id));
	CREATE TABLE public.orders(
		id serial,
		customer varchar(200) not null,
		recipe_id integer not null,
		date timestamp default current_timestamp,
		release_date timestamp,
		price numeric(15,2) not null,
		plan_qty numeric(15,2) not null,
		plan_cost numeric(15,2) not null,
		fact_qty numeric(15,2) not null,
		fact_cost numeric(15,2) not null,
		materials_cost numeric(15,2) not null,
		constraint id_orders primary key (id));
	CREATE TABLE public.order_details(
		id integer not null,
		material_id integer not null,
		qty numeric(15,2),
		string_order integer not null,
		price numeric(15,2) not null,
		cost numeric(15,2));
		 `
	_, err := db.Exec(initTableString)
	if err != nil {
		fmt.Println("Error creating tables:", err)
		panic("Error creating tables")
	}
	setDbVersion(0, db)

}

//UpdateData Запись данных в таблицу
func (mdb MDB) UpdateData(tableName string, data map[string]interface{}) int {
	db := *mdb.Pdb

	rows, err := db.Query(`SELECT * FROM ` + tableName + ` WHERE false`)
	if err != nil {
		fmt.Println("Error reading data: ", err)
		return -1
	}
	cols, _ := rows.Columns()

	vals := []interface{}{}
	namesString := ""
	placeholdersString := ""
	newRecord := false
	valueIndex := 1
	for _, col := range cols {
		if col == "id" && data[col] == -1 {
			newRecord = true
			continue
		}
		if data[col] == nil {
			continue
		}
		namesString += (col + ",")
		placeholdersString += ("$" + strconv.Itoa(valueIndex) + ",")
		vals = append(vals, data[col])
		valueIndex++
	}
	namesString = strings.TrimSuffix(namesString, ",")
	placeholdersString = strings.TrimSuffix(placeholdersString, ",")

	if newRecord {
		lastInsertID := -1
		err := db.QueryRow(`insert into public.`+tableName+`(`+namesString+`) VALUES(`+placeholdersString+`) RETURNING id;`, vals...).Scan(&lastInsertID)
		if err != nil {
			fmt.Println("Error inserting data: ", err)
			return -1
		}
		return lastInsertID
	}
	_, err = db.Exec(`UPDATE public.`+tableName+` SET (`+namesString+`) = (`+placeholdersString+`) WHERE id = `+strconv.Itoa(data["id"].(int))+`;`, vals...)
	if err != nil {
		fmt.Println("Error updating data: ", err)
		return -1
	}
	return data["id"].(int)

}

//UpdateTableData Обновляет данные в табличной части
func (mdb MDB) UpdateTableData(tableName string, data []map[string]interface{}, id int) bool {
	db := *mdb.Pdb
	rows, err := db.Query(`SELECT * FROM ` + tableName + ` WHERE false`)
	if err != nil {
		fmt.Println("Error reading data: ", err)
		return false
	}
	cols, _ := rows.Columns()

	_, err = db.Exec(`DELETE FROM public.` + tableName + ` WHERE id = ` + strconv.Itoa(id) + `;`)
	if err != nil {
		fmt.Println("Error deleting old rows: ", err)
		return false
	}
	if len(data) == 0 {
		return true
	}
	namesString := ""
	placeholdersString := ""
	valueIndex := 1
	for _, col := range cols {
		if data[0][col] == nil {
			continue
		}
		namesString += (col + ",")
		placeholdersString += ("$" + strconv.Itoa(valueIndex) + ",")
		valueIndex++
	}
	namesString = strings.TrimSuffix(namesString, ",")
	placeholdersString = strings.TrimSuffix(placeholdersString, ",")
	stmt, err := db.Prepare(`insert into public.` + tableName + `(` + namesString + `) VALUES(` + placeholdersString + `);`)
	if err != nil {
		fmt.Println("Error preparing for inserting data: ", err)
		return false
	}
	for _, item := range data {
		vals := []interface{}{}
		for _, col := range cols {
			if item[col] == nil {
				continue
			}
			vals = append(vals, item[col])
		}
		_, err := stmt.Exec(vals...)
		if err != nil {
			fmt.Println("Error inserting data: ", err)
			return false
		}
	}
	return true
}

//DeleteData уделение данных в таблице по id
func (mdb MDB) DeleteData(tableName string, id int) bool {
	db := *mdb.Pdb
	_, err := db.Exec(`DELETE FROM public.` + tableName + ` WHERE id = ` + strconv.Itoa(id) + `;`)
	if err != nil {
		fmt.Println("Error deleting old rows: ", err)
		return false
	}
	return true
}

//ReadRow Выполняет запрос к SQL и возвращает одну строку в виде map
func (mdb MDB) ReadRow(queryText string) (bool, map[string]interface{}) {
	data := make(map[string]interface{})
	db := *mdb.Pdb
	rows, err := db.Query(queryText)
	if err != nil {
		fmt.Println("Error reading data: ", err)
		return false, data
	}
	cols, _ := rows.Columns()
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	// Scan the result into the column pointers...
	if err := rows.Scan(columnPointers...); err != nil {
		fmt.Println("Error reading data: ", err)
		return false, data
	}

	for i, colName := range cols {
		val := columnPointers[i].(*interface{})
		data[colName] = *val
	}
	return true, data
}

//ReadRows Выполняет запрос к SQL и возвращает строки в виде []map
func (mdb MDB) ReadRows(queryText string) (bool, []map[string]interface{}) {

	dataArray := make([]map[string]interface{}, 0)
	db := *mdb.Pdb
	rows, err := db.Query(queryText) // Note: Ignoring errors for brevity
	if err != nil {
		fmt.Println("Error reading data: ", err)
		return false, dataArray
	}
	cols, _ := rows.Columns()
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	for rows.Next() {
		if err := rows.Scan(columnPointers...); err != nil {
			fmt.Println("Error reading data: ", err)
			return false, dataArray
		}
		data := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			data[colName] = *val
		}
		dataArray = append(dataArray, data)
	}
	return true, dataArray
}

//ReadOrderList Читает список заказов
func (mdb MDB) ReadOrderList() (bool, []map[string]interface{}) {
	return mdb.ReadRows(OrdersQuery)
}

//ReadOrder Читает заказ по id
func (mdb MDB) ReadOrder(id int) (bool, map[string]interface{}) {
	res, data := mdb.ReadRow(GetOrderQuerry(id))
	if !res {
		return res, data
	}
	res, content := mdb.ReadOrderContent(id)
	if !res {
		return res, data
	}
	data["content"] = content
	return res, data
}

//ReadOrderContent Читает табличную часть заказа
func (mdb MDB) ReadOrderContent(id int) (bool, []map[string]interface{}) {
	return mdb.ReadRows(GetOrderContentQuerry(id))
}

//ReadRecipesList Читает список рецептов
func (mdb MDB) ReadRecipesList() (bool, []map[string]interface{}) {
	return mdb.ReadRows(GetRowsQuerry("recipes"))
}

//ReadRecipe Читает рецепт по id
func (mdb MDB) ReadRecipe(id int) (bool, map[string]interface{}) {
	res, data := mdb.ReadRow(GetRowQuerry("recipes", id))
	if !res {
		return res, data
	}
	res, content := mdb.ReadRecipeContent(id)
	if !res {
		return res, data
	}
	data["content"] = content
	return res, data
}

//ReadRecipeContent Читает табличную часть рецепта по id
func (mdb MDB) ReadRecipeContent(id int) (bool, []map[string]interface{}) {
	return mdb.ReadRows(GetRecipeContentQuery(id))
}

//ReadRecipeContentWithPrice Читает табличную часть рецепта с ценами по id
func (mdb MDB) ReadRecipeContentWithPrice(id int) (bool, []map[string]interface{}) {
	return mdb.ReadRows(GetRecipeContentWithPricesQuery(id))
}

//ReadUnits читает список единиц измерения
func (mdb MDB) ReadUnits() (bool, []map[string]interface{}) {
	return mdb.ReadRows(GetRowsQuerry("units"))
}

//ReadMaterials читает список материалов
func (mdb MDB) ReadMaterials(prices bool) (bool, []map[string]interface{}) {
	if prices {
		return mdb.ReadRows(GetMaterialsWithPricesQuery())
	}
	return mdb.ReadRows(GetMaterialsQuery())
}

//ReadMaterial читает материал по id
func (mdb MDB) ReadMaterial(prices bool, id int) (bool, map[string]interface{}) {
	if prices {
		return mdb.ReadRow(GetMaterialWithPricesQuery(id))
	}
	return mdb.ReadRow(GetMaterialQuery(id))
}
