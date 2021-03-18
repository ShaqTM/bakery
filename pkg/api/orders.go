package api

import (
	"bakery/pkg/store"
	"encoding/json"
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
			mdb.Log.Error("Error reading querry:", err)
			sendAnswer400(w, err.Error())
			return
		}
		dataMap := make(map[string]interface{})
		err = json.Unmarshal(b, &dataMap)
		if err != nil {
			mdb.Log.Error("Error unmarshal query: ", err)
			mdb.Log.Error(string(b))
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
		orders, err := mdb.ReadOrders()
		if err != nil {
			sendAnswer405(w, err.Error())
			return
		}
		jsonArray, err := json.Marshal(orders)
		if err != nil {
			mdb.Log.Error("Error marshal orders:", err)
			mdb.Log.Error(orders)
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
			mdb.Log.Error("Error reading id:", query["id"][0])
			mdb.Log.Error(err.Error())
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
			mdb.Log.Error("Error marshal order:", err)
			mdb.Log.Error(order)
			sendAnswer400(w, err.Error())
			return
		}
		jsonString := string(jsonText)
		sendAnswer200(w, jsonString)
	})
}
