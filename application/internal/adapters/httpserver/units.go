package httpserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// AddUnitsRoutes Добавляет обработку роутов
func init() {
	routes = append(routes, route{
		Methods: []string{"POST", "OPTIONS"},
		Path:    "/api/writeunit",
		Handler: (*Server).writeUnit,
	})
	routes = append(routes, route{
		Methods: []string{"GET", "OPTIONS"},
		Path:    "/api/readunits",
		Handler: (*Server).readUnits,
	})
	routes = append(routes, route{
		Methods: []string{"GET", "OPTIONS"},
		Path:    "/api/readunit/",
		Handler: (*Server).readUnit,
	})
}

func (s *Server) writeUnit() http.Handler {
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
			s.log.Error("Error reading querry:", err)
			sendAnswer400(w, err.Error())
			return
		}
		dataMap := make(map[string]interface{})
		err = json.Unmarshal(b, &dataMap)
		if err != nil {
			s.log.Error("Error unmarshal query: ", err)
			s.log.Error(string(b))
			sendAnswer400(w, err.Error())
			return
		}
		id, err := mdb.UpdateData("units", dataMap)
		if err == nil {
			sendAnswer202(w, strconv.Itoa(id))
		} else {
			sendAnswer400(w, err.Error())
		}

	})
}

func (s *Server) readUnits() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//		if r.Method != "GET" {
		//			sendAnswer405(w, "bad method")
		//			return
		//		}
		if r.Method == "OPTIONS" {
			sendAnswer200(w, "")
			return
		}
		units, err := mdb.ReadUnits()
		if err != nil {
			sendAnswer405(w, err.Error())
			return
		}
		jsonArray, err := json.Marshal(units)
		if err != nil {
			s.Log.Error("Error marshal units:", err)
			s.Log.Error(units)
			sendAnswer400(w, err.Error())
			return
		}
		jsonString := string(jsonArray)
		sendAnswer200(w, jsonString)
	})
}

func (s *Server) readUnit() http.Handler {
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
			s.Log.Error("Error reading id:", query["id"][0])
			s.Log.Error(err.Error())
			sendAnswer400(w, err.Error())
			return
		}
		unit, err := mdb.ReadUnit(id)
		if err != nil {
			sendAnswer405(w, err.Error())
			return
		}
		jsonText, err := json.Marshal(unit)
		if err != nil {
			s.Log.Error("Error marshal unit:", err)
			s.Log.Error(unit)
			sendAnswer400(w, err.Error())
			return
		}
		jsonString := string(jsonText)
		sendAnswer200(w, jsonString)
	})
}
