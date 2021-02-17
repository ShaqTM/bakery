package api

import (
	"bakery/pkg/store"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//AddMaterialsRoutes Добавляет обработку роутов
func AddMaterialsRoutes(router **mux.Router, mdb store.MDB) {
	(*router).Handle("/api/writematerial", writeMaterial(mdb)).Methods("POST", "OPTIONS")
	(*router).Handle("/api/readmaterials", readMaterials(mdb)).Methods("GET", "OPTIONS")
	(*router).Handle("/api/readmaterial/", readMaterial(mdb)).Methods("GET", "OPTIONS")

	(*router).Handle("/api/writeprice", writePrice(mdb)).Methods("POST", "OPTIONS")
}

func writeMaterial(mdb store.MDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		if r.Method != "POST" {
		//			sendAnswer405(w, "bad method")
		//			return
		//		}
		if r.Method == "OPTIONS" {
			sendAnswer200(w, "")
			return
		}
		rc := r.Body
		b, err := ioutil.ReadAll(rc)
		if err != nil {
			fmt.Println("error reading querry:", err)
			sendAnswer400(w, "")
			return
		}
		dataMap := make(map[string]interface{})
		err = json.Unmarshal(b, &dataMap)
		if err != nil {
			fmt.Println(string(b))
			sendAnswer400(w, err.Error())
			return
		}
		id, err := mdb.UpdateData("materials", dataMap)
		if err == nil {
			sendAnswer202(w, strconv.Itoa(id))
		} else {
			sendAnswer400(w, err.Error())
		}

	})
}

func readMaterials(mdb store.MDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		if r.Method != "GET" {
		//			sendAnswer405(w, "bad method")
		//			return
		//		}
		if r.Method == "OPTIONS" {
			sendAnswer200(w, "")
			return
		}
		query := r.URL.Query()
		withPrice := query["price"] != nil && query["price"][0] == "true"
		units, err := mdb.ReadMaterials(withPrice)
		if err != nil {
			sendAnswer405(w, err.Error())
			return
		}
		jsonArray, err := json.Marshal(units)
		if err != nil {
			fmt.Println("error reading querry:", err)
			sendAnswer400(w, err.Error())
			return
		}
		jsonString := string(jsonArray)
		sendAnswer200(w, jsonString)
	})
}

func readMaterial(mdb store.MDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		if r.Method != "GET" {
		//			sendAnswer405(w, "bad method")
		//			return
		//		}
		if r.Method == "OPTIONS" {
			sendAnswer200(w, "")
			return
		}
		query := r.URL.Query()
		if query["id"] == nil {
			sendAnswer400(w, "bad parameters")
			return
		}
		withPrice := query["price"] != nil && query["price"][0] == "true"
		id, err := strconv.Atoi(query["id"][0])
		if err != nil {
			sendAnswer400(w, err.Error())
			return
		}
		unit, err := mdb.ReadMaterial(withPrice, id)
		if err != nil {
			sendAnswer405(w, err.Error())
			return
		}
		jsonText, err := json.Marshal(unit)
		if err != nil {
			fmt.Println("error reading querry:", err)
			sendAnswer400(w, err.Error())
			return
		}
		jsonString := string(jsonText)
		sendAnswer200(w, jsonString)
	})
}

func writePrice(mdb store.MDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		if r.Method != "POST" {
		//			sendAnswer405(w, "bad method")
		//			return
		//		}
		if r.Method == "OPTIONS" {
			sendAnswer200(w, "")
			return
		}
		rc := r.Body
		b, err := ioutil.ReadAll(rc)
		if err != nil {
			fmt.Println("error reading querry:", err)
			sendAnswer400(w, "")
			return
		}
		dataMap := make(map[string]interface{})
		err = json.Unmarshal(b, &dataMap)
		if err != nil {
			fmt.Println(string(b))
			sendAnswer400(w, err.Error())
			return
		}
		id, err := mdb.UpdateData("material_prices", dataMap)
		if err == nil {
			sendAnswer202(w, strconv.Itoa(id))
		} else {
			sendAnswer400(w, err.Error())
		}

	})
}
