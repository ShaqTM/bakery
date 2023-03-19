package store

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

const dbConnectString = "host=localhost port=5432 user=postgres password=qazplm dbname=bakery sslmode=disable"
const dbConnectStringInit = "host=localhost port=5432 user=postgres password=qazplm dbname=postgres sslmode=disable"

//InitDatabase Проверка наличия БД, создание и обновление до последней версии
func (mdb *MDB) InitDatabase() {

	initdb, err := sql.Open("postgres", dbConnectStringInit)

	if err != nil {
		mdb.Log.Fatalf("Database opening error:", err)
	}
	defer initdb.Close()

	rows, err := initdb.Query("SELECT datname FROM pg_database WHERE datistemplate = false AND datname = 'bakery';")

	if err != nil {
		mdb.Log.Fatalf("Error searching database:", err)
	}
	if rows.Next() {
		mdb.Log.Info("Database bakery found")
	} else {
		createDb(initdb, mdb.Log)
	}
	initdb.Close()
	db, err := sql.Open("postgres", dbConnectString)
	if err != nil {
		mdb.Log.Fatalf("Database opening error:", err)
	}
	updateDb(db, mdb.Log)
	(*mdb).Pdb = &db

}
func createDb(db *sql.DB, log *logrus.Logger) {
	_, err := db.Exec("CREATE DATABASE bakery WITH OWNER postgres;")
	if err != nil {
		log.Fatalf("Error creating database:", err)
	}
	fmt.Println("Database created successfully")

}

func updateDb(db *sql.DB, log *logrus.Logger) {
	var dbVersion int
	dbVersion = -1
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'settings'")
	if err != nil {
		log.Fatalf("Error searching table settings:", err)
	}
	if rows.Next() {
		log.Info("Table settings found")
		dbVersion = getDbVersion(db, log)
	} else {
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS public.settings(
			name varchar(50) not null,
			string_value varchar(200),
			numeric_value numeric(15,3),
			integer_Value integer,
			time_value timestamp,
			constraint name_data primary key (name));`)
		if err != nil {
			log.Fatalf("Table settings create error", err)
		}
	}
	if dbVersion < 0 {
		updateDbToVerion0(db, log)
	}
	if dbVersion < 1 {
		updateDbToVerion1(db, log)
	}
	if dbVersion < 2 {
		updateDbToVerion2(db, log)
	}
}
func getDbVersion(db *sql.DB, log *logrus.Logger) int {
	rows, err := db.Query(`SELECT integer_Value FROM public.settings WHERE name = 'dbVersion'`)
	if err != nil {
		log.Fatalf("Error qetting dbVersion data: ", err)
	}
	var version int
	version = -1
	if rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			log.Fatalf("Error scanning dbVersion data: ", err)
		}
	}
	return version
}

func setDbVersion(version float32, db *sql.DB, log *logrus.Logger) {
	insertString := ""
	if getDbVersion(db, log) == -1 {

		insertString = `INSERT INTO public.settings(name, integer_Value) VALUES ($1,$2);`
	} else {
		insertString = `UPDATE public.settings SET integer_Value =$2 WHERE name = $1;`
	}
	_, err := db.Exec(insertString, "dbVersion", version)

	if err != nil {
		log.Fatalf("Error inserting dbVersion data: ", err)
	}
}

func updateDbToVerion0(db *sql.DB, log *logrus.Logger) {
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
		log.Fatalf("Error creating tables:", err)
		panic("Error creating tables")
	}
	setDbVersion(0, db, log)

}

func updateDbToVerion1(db *sql.DB, log *logrus.Logger) {
	initTableString := `
	ALTER TABLE public.order_details ADD COLUMN by_recipe boolean;
		 `
	_, err := db.Exec(initTableString)
	if err != nil {
		log.Fatalf("Error updating tables to version 1:", err)
	}
	setDbVersion(1, db, log)

}

func updateDbToVerion2(db *sql.DB, log *logrus.Logger) {
	initTableString := `
	ALTER TABLE public.order_details ALTER COLUMN price SET DATA TYPE numeric(15,4);
		 `
	_, err := db.Exec(initTableString)
	if err != nil {
		log.Fatalf("Error updating tables to version 2:", err)
	}
	setDbVersion(2, db, log)

}
