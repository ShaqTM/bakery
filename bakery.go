package main

import (
	"bakery/pkg/api"
	"bakery/pkg/store"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	var db *sql.DB
	mdb := store.MDB{Pdb: &db}
	//	fmt.Println(mdb.Pdb)
	(&mdb).InitDatabase()
	//	fmt.Println(mdb.Pdb)
	router := mux.NewRouter()
	api.AddMaterialsRoutes(&router, mdb)
	api.AddUnitsRoutes(&router, mdb)
	api.AddOrdersRoutes(&router, mdb)
	api.AddRecipesRoutes(&router, mdb)

	http.ListenAndServe(":5000", router)
	for {

	}
}
