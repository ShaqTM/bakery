package httpserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func init() {
	routes = append(routes, route{
		Methods: []string{"POST", "OPTIONS"},
		Path:    "/api/writeorder",
		Handler: (*Server).writeOrder,
	})
	routes = append(routes, route{
		Methods: []string{"GET", "OPTIONS"},
		Path:    "/api/readorders",
		Handler: (*Server).readOrders,
	})
	routes = append(routes, route{
		Methods: []string{"GET", "OPTIONS"},
		Path:    "/api/readorder/",
		Handler: (*Server).readOrder,
	})

}

func (s *Server) writeOrder() http.Handler {
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
			s.Log.Error("Error reading querry:", err)
			sendAnswer400(w, err.Error())
			return
		}
		dataMap := make(map[string]interface{})
		err = json.Unmarshal(b, &dataMap)
		if err != nil {
			s.Log.Error("Error unmarshal query: ", err)
			s.Log.Error(string(b))
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

func (s *Server) readOrders() http.Handler {
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
			s.Log.Error("Error marshal orders:", err)
			s.Log.Error(orders)
			sendAnswer400(w, err.Error())
			return
		}
		jsonString := string(jsonArray)
		sendAnswer200(w, jsonString)
	})
}

func (s *Server) readOrder() http.Handler {
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
