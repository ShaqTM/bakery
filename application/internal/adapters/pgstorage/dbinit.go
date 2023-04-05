package pgstorage

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

// InitDatabase Проверка наличия БД, создание и обновление до последней версии
func (s *Storage) Start() {

	initdb, err := sql.Open("postgres", s.config.PG_connect_string_init)

	if err != nil {
		s.Log.Fatalf("Database opening error:", err)
	}
	defer initdb.Close()

	rows, err := initdb.Query("SELECT datname FROM pg_database WHERE datistemplate = false AND datname = 'bakery';")

	if err != nil {
		s.Log.Fatalf("Error searching database:", err)
	}
	if rows.Next() {
		s.Log.Info("Database bakery found")
	} else {
		createDb(initdb, s.Log)
	}

	s.Pdb, err = sql.Open("postgres", s.config.PG_connect_string)
	if err != nil {
		s.Log.Fatalf("Database opening error:", err)
	}
	s.updateDb()

}
func createDb(db *sql.DB, log *logrus.Logger) {
	_, err := db.Exec("CREATE DATABASE bakery WITH OWNER postgres;")
	if err != nil {
		log.Fatalf("Error creating database:", err)
	}
	fmt.Println("Database created successfully")

}

func (s *Storage) updateDb() {
	var dbVersion int
	dbVersion = -1
	rows, err := s.Pdb.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'settings'")
	if err != nil {
		s.Log.Fatalf("Error searching table settings:", err)
	}
	if rows.Next() {
		s.Log.Info("Table settings found")
		dbVersion = s.getDbVersion()
	} else {
		_, err = s.Pdb.Exec(`CREATE TABLE IF NOT EXISTS public.settings(
			name varchar(50) not null,
			string_value varchar(200),
			numeric_value numeric(15,3),
			integer_Value integer,
			time_value timestamp,
			constraint name_data primary key (name));`)
		if err != nil {
			s.Log.Fatalf("Table settings create error", err)
		}
	}
	if dbVersion < 0 {
		s.updateDbToVerion0()
	}
	if dbVersion < 1 {
		s.updateDbToVerion1()
	}
	if dbVersion < 2 {
		s.updateDbToVerion2()
	}
}
func (s *Storage) getDbVersion() int {
	rows, err := s.Pdb.Query(`SELECT integer_Value FROM public.settings WHERE name = 'dbVersion'`)
	if err != nil {
		s.Log.Fatalf("Error qetting dbVersion data: ", err)
	}
	var version int
	version = -1
	if rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			s.Log.Fatalf("Error scanning dbVersion data: ", err)
		}
	}
	return version
}

func (s *Storage) setDbVersion(version float32) {
	insertString := ""
	if s.getDbVersion() == -1 {
		insertString = `INSERT INTO public.settings(name, integer_Value) VALUES ($1,$2);`
	} else {
		insertString = `UPDATE public.settings SET integer_Value =$2 WHERE name = $1;`
	}
	_, err := s.Pdb.Exec(insertString, "dbVersion", version)

	if err != nil {
		s.Log.Fatalf("Error inserting dbVersion data: ", err)
	}
}

func (s *Storage) updateDbToVerion0() {
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
	_, err := s.Pdb.Exec(initTableString)
	if err != nil {
		s.Log.Fatalf("Error creating tables:", err)
		panic("Error creating tables")
	}
	s.setDbVersion(0)

}

func (s *Storage) updateDbToVerion1() {
	initTableString := `
	ALTER TABLE public.order_details ADD COLUMN by_recipe boolean;
		 `
	_, err := s.Pdb.Exec(initTableString)
	if err != nil {
		s.Log.Fatalf("Error updating tables to version 1:", err)
	}
	s.setDbVersion(1)

}

func (s *Storage) updateDbToVerion2() {
	initTableString := `
	ALTER TABLE public.order_details ALTER COLUMN price SET DATA TYPE numeric(15,4);
		 `
	_, err := s.Pdb.Exec(initTableString)
	if err != nil {
		s.Log.Fatalf("Error updating tables to version 2:", err)
	}
	s.setDbVersion(2)

}
