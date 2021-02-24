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

//AddRecipesRoutes Добавляет обработку роутов
func AddRecipesRoutes(router **mux.Router, mdb store.MDB) {
	(*router).Handle("/api/writerecipe", writeRecipe(mdb)).Methods("POST", "OPTIONS")
	(*router).Handle("/api/readrecipes", readRecipes(mdb)).Methods("GET", "OPTIONS")
	(*router).Handle("/api/readrecipe/", readRecipe(mdb)).Methods("GET", "OPTIONS")

	(*router).Handle("/api/writerecipeprice", writeRecipePrice(mdb)).Methods("POST", "OPTIONS")
}

func writeRecipe(mdb store.MDB) http.Handler {
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

func readRecipes(mdb store.MDB) http.Handler {
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
		units, err := mdb.ReadRecipes(withPrice)
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

func readRecipe(mdb store.MDB) http.Handler {
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
		recipe, err := mdb.ReadRecipe(withPrice, id)
		if err != nil {
			sendAnswer405(w, err.Error())
			return
		}
		jsonText, err := json.Marshal(recipe)
		if err != nil {
			fmt.Println("error reading querry:", err)
			sendAnswer400(w, err.Error())
			return
		}
		jsonString := string(jsonText)
		sendAnswer200(w, jsonString)
	})
}

func writeRecipePrice(mdb store.MDB) http.Handler {
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
		id, err := mdb.UpdateData("recipe_prices", dataMap)
		if err == nil {
			sendAnswer202(w, strconv.Itoa(id))
		} else {
			sendAnswer400(w, err.Error())
		}

	})
}
