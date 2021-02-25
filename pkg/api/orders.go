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

//AddOrdersRoutes Добавляет обработку роутов
func AddOrdersRoutes(router **mux.Router, mdb store.MDB) {
	(*router).Handle("/api/writeorder", writeOrder(mdb)).Methods("POST", "OPTIONS")
	(*router).Handle("/api/readorders", readOrders(mdb)).Methods("GET", "OPTIONS")
	(*router).Handle("/api/readorder/", readOrder(mdb)).Methods("GET", "OPTIONS")
}

func writeOrder(mdb store.MDB) http.Handler {
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
		id, err := mdb.UpdateData("orders", dataMap)
		if err != nil {
			sendAnswer400(w, err.Error())
			return
		}
		err = mdb.UpdateTableData("order_details", dataMap["content"].([]interface{}), id)
		if err != nil {
			sendAnswer400(w, err.Error())
			return
		}

		sendAnswer202(w, strconv.Itoa(id))
	})
}

func readOrders(mdb store.MDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		if r.Method != "GET" {
		//			sendAnswer405(w, "bad method")
		//			return
		//		}
		if r.Method == "OPTIONS" {
			sendAnswer200(w, "")
			return
		}
		units, err := mdb.ReadOrders()
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

func readOrder(mdb store.MDB) http.Handler {
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
		id, err := strconv.Atoi(query["id"][0])
		if err != nil {
			sendAnswer400(w, err.Error())
			return
		}
		order, err := mdb.ReadOrder(id)
		if err != nil {
			sendAnswer405(w, err.Error())
			return
		}
		jsonText, err := json.Marshal(order)
		if err != nil {
			fmt.Println("error reading querry:", err)
			sendAnswer400(w, err.Error())
			return
		}
		jsonString := string(jsonText)
		sendAnswer200(w, jsonString)
	})
}
