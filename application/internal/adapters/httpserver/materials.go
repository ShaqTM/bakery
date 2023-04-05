package httpserver

import (
	"bakery/application/internal/domain/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func init() {
	routes = append(routes, route{
		Methods: []string{"GET", "OPTIONS"},
		Path:    "/api/readmaterials",
		Handler: (*Server).readMaterials,
	})
	routes = append(routes, route{
		Methods: []string{"GET", "OPTIONS"},
		Path:    "/api/readmaterial/",
		Handler: (*Server).readMaterial,
	})
	routes = append(routes, route{
		Methods: []string{"POST", "OPTIONS"},
		Path:    "/api/writematerialprice",
		Handler: (*Server).writeMaterialPrice,
	})
	routes = append(routes, route{
		Methods: []string{"POST", "OPTIONS"},
		Path:    "/api/writematerial",
		Handler: (*Server).writeMaterial,
	})

}

func (s *Server) writeMaterial() http.Handler {
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
			sendAnswer400(w, "")
			return
		}
		material := models.Material{}
		err = json.Unmarshal(b, &material)
		if err != nil {
			s.Log.Error("Error unmarshal query")
			s.Log.Error(string(b))
			sendAnswer400(w, err.Error())
			return
		}
		id, err := s.Bakery.WriteMaterial(material)
		if err == nil {
			sendAnswer202(w, strconv.Itoa(id))
		} else {
			sendAnswer400(w, err.Error())
		}

	})
}

func (s *Server) readMaterials() http.Handler {
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
		materials, err := s.Bakery.ReadMaterials(withPrice)
		if err != nil {
			sendAnswer405(w, err.Error())
			return
		}
		jsonArray, err := json.Marshal(materials)
		if err != nil {
			s.Log.Error("Error marshal JSON:", err)
			s.Log.Error(materials)
			sendAnswer400(w, err.Error())
			return
		}
		jsonString := string(jsonArray)
		sendAnswer200(w, jsonString)
	})
}

func (s *Server) readMaterial() http.Handler {
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
		material, err := s.Bakery.ReadMaterial(withPrice, id)
		if err != nil {
			sendAnswer405(w, err.Error())
			return
		}
		jsonText, err := json.Marshal(material)
		if err != nil {
			s.Log.Error("Error marshal JSON:", err)
			s.Log.Error(material)
			sendAnswer400(w, err.Error())
			return
		}
		jsonString := string(jsonText)
		sendAnswer200(w, jsonString)
	})
}

func (s *Server) writeMaterialPrice() http.Handler {
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
		material_price := models.Material_price{}
		err = json.Unmarshal(b, &material_price)
		if err != nil {
			s.Log.Error("Error unmarshal query: ", err)
			s.Log.Error(string(b))
			sendAnswer400(w, err.Error())
			return
		}
		id, err := s.Bakery.WriteMaterialPrice(material_price)
		if err == nil {
			sendAnswer202(w, strconv.Itoa(id))
		} else {
			sendAnswer400(w, err.Error())
		}

	})
}
