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
		Path:    "/api/writerecipe",
		Handler: (*Server).writeRecipe,
	})
	routes = append(routes, route{
		Methods: []string{"GET", "OPTIONS"},
		Path:    "/api/readrecipes",
		Handler: (*Server).readRecipes,
	})
	routes = append(routes, route{
		Methods: []string{"GET", "OPTIONS"},
		Path:    "/api/readrecipe/",
		Handler: (*Server).readRecipe,
	})
	routes = append(routes, route{
		Methods: []string{"POST", "OPTIONS"},
		Path:    "/api/writerecipeprice",
		Handler: (*Server).writeRecipePrice,
	})

}

func (s *Server) writeRecipe() http.Handler {
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
			s.Log.Error("Error unmarshal query: ", err.Error())
			s.Log.Error(string(b))
			sendAnswer400(w, err.Error())
			return
		}
		id, err := mdb.UpdateData("recipes", dataMap)
		if err != nil {
			sendAnswer400(w, err.Error())
		}
		err = mdb.UpdateTableData("recipes_content", dataMap["content"].([]interface{}), id)
		if err != nil {
			sendAnswer400(w, err.Error())
		}

		sendAnswer202(w, strconv.Itoa(id))
	})
}

func (s *Server) readRecipes() http.Handler {
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
		recipes, err := mdb.ReadRecipes(withPrice)
		if err != nil {
			sendAnswer405(w, err.Error())
			return
		}
		jsonArray, err := json.Marshal(recipes)
		if err != nil {
			s.Log.Error("Error marshal recipes:", err)
			s.Log.Error(recipes)
			sendAnswer400(w, err.Error())
			return
		}
		jsonString := string(jsonArray)
		sendAnswer200(w, jsonString)
	})
}

func (s *Server) readRecipe() http.Handler {
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
			s.Log.Error("Error reading id:", query["id"][0])
			s.Log.Error(err.Error())
			sendAnswer400(w, err.Error())
			return
		}
		recipe, err := mdb.ReadRecipe(withPrice, id)
		if err != nil {
			sendAnswer405(w, err.Error())
			return
		}
		jsonText, err := json.Marshal(recipe)
		if err != nil {
			s.Log.Error("Error marshal recipe:", err)
			s.Log.Error(recipe)
			sendAnswer400(w, err.Error())
			return
		}
		jsonString := string(jsonText)
		sendAnswer200(w, jsonString)
	})
}

func (s *Server) writeRecipePrice() http.Handler {
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
		dataMap := make(map[string]interface{})
		err = json.Unmarshal(b, &dataMap)
		if err != nil {
			s.Log.Error("Error unmarshal query: ", err)
			s.Log.Error(string(b))
			sendAnswer400(w, err.Error())
			return
		}
		id, err := mdb.UpdateData("recipe_prices", dataMap)
		if err == nil {
			sendAnswer202(w, strconv.Itoa(id))
		} else {
			sendAnswer400(w, err.Error())
		}

	})
}
