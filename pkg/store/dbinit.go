package store

import (
	"database/sql"
	"fmt"
)

const dbConnectString = "host=localhost port=5432 user=postgres password=qazplm dbname=bakery sslmode=disable"
const dbConnectStringInit = "host=localhost port=5432 user=postgres password=qazplm dbname=postgres sslmode=disable"

//InitDatabase Проверка наличия БД, создание и обновление до последней версии
func (mdb *MDB) InitDatabase() {

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
	(*mdb).Pdb = &db

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
	var dbVersion int
	dbVersion = -1
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
			integer_Value integer,
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
	if dbVersion < 1 {
		updateDbToVerion1(db)
	}
	if dbVersion < 2 {
		updateDbToVerion2(db)
	}
}
func getDbVersion(db *sql.DB) int {
	rows, err := db.Query(`SELECT integer_Value FROM public.settings WHERE name = 'dbVersion'`)
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

func setDbVersion(version float32, db *sql.DB) {
	insertString := ""
	if getDbVersion(db) == -1 {

		insertString = `INSERT INTO public.settings(name, integer_Value) VALUES ($1,$2);`
	} else {
		insertString = `UPDATE public.settings SET integer_Value =$2 WHERE name = $1;`
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
		short_name varchar(200) not null,
		constraint id_units primary key (id));
	CREATE TABLE public.materials(
		id serial,
		name varchar(200) not null,
		recipe_unit_id integer not null,
		price_unit_id integer not null,
		coefficient numeric(15,3) not null,
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
		output numeric(15,3) not null,
		unit_id integer not null,
		constraint id_recipes primary key (id));
	CREATE TABLE public.recipes_content(
		id integer not null,
		material_id integer not null,
		string_order integer not null,
		qty numeric(15,3) not null);
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
		unit_id integer not null,
		price numeric(15,2) not null,
		plan_qty numeric(15,3) not null,
		plan_cost numeric(15,2) not null,
		fact_qty numeric(15,3) not null,
		fact_cost numeric(15,2) not null,
		materials_cost numeric(15,2) not null,
		constraint id_orders primary key (id));
	CREATE TABLE public.order_details(
		id integer not null,
		material_id integer not null,
		qty numeric(15,3),
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

func updateDbToVerion1(db *sql.DB) {
	initTableString := `
	ALTER TABLE public.order_details ADD COLUMN by_recipe boolean;
		 `
	_, err := db.Exec(initTableString)
	if err != nil {
		fmt.Println("Error updating tables:", err)
		panic("Error updating tables")
	}
	setDbVersion(1, db)

}

func updateDbToVerion2(db *sql.DB) {
	initTableString := `
	ALTER TABLE public.order_details ALTER COLUMN price SET DATA TYPE numeric(15,4);
		 `
	_, err := db.Exec(initTableString)
	if err != nil {
		fmt.Println("Error updating tables:", err)
		panic("Error updating tables")
	}
	setDbVersion(2, db)

}
