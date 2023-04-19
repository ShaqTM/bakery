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
		Method:  "POST",
		Path:    "/api/writerecipe",
		Handler: (*Server).writeRecipe,
	})
	routes = append(routes, route{
		Method:  "GET",
		Path:    "/api/readrecipes",
		Handler: (*Server).readRecipes,
	})
	routes = append(routes, route{
		Method:  "GET",
		Path:    "/api/readrecipe/",
		Handler: (*Server).readRecipe,
	})
	routes = append(routes, route{
		Method:  "POST",
		Path:    "/api/writerecipeprice",
		Handler: (*Server).writeRecipePrice,
	})

}

func (s *Server) writeRecipe() http.HandlerFunc {
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
		recipe := models.Recipe{}
		err = json.Unmarshal(b, &recipe)
		if err != nil {
			s.Log.Error("Error unmarshal query: ", err.Error())
			s.Log.Error(string(b))
			sendAnswer400(w, err.Error())
			return
		}
		id, err := s.Bakery.WriteRecipe(recipe)
		if err != nil {
			sendAnswer400(w, err.Error())
		}
		sendAnswer202(w, strconv.Itoa(id))
	})
}

func (s *Server) readRecipes() http.HandlerFunc {
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
		recipes, err := s.Bakery.ReadRecipes(withPrice)
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

func (s *Server) readRecipe() http.HandlerFunc {
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
		recipe, err := s.Bakery.ReadRecipe(withPrice, id)
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

func (s *Server) writeRecipePrice() http.HandlerFunc {
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
		recipe_price := models.Recipe_price{}
		err = json.Unmarshal(b, &recipe_price)
		if err != nil {
			s.Log.Error("Error unmarshal query: ", err)
			s.Log.Error(string(b))
			sendAnswer400(w, err.Error())
			return
		}
		id, err := s.Bakery.WriteRecipePrice(recipe_price)
		if err == nil {
			sendAnswer202(w, strconv.Itoa(id))
		} else {
			sendAnswer400(w, err.Error())
		}

	})
}
