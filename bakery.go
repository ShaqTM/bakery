package main

import (
	dbservice "bakery/dbService"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	var db *sql.DB
	mdb := dbservice.MDB{Pdb: &db}
	//	fmt.Println(mdb.Pdb)
	(&mdb).InitDatabase()
	//	fmt.Println(mdb.Pdb)
	router := mux.NewRouter()
	router.Handle("/api/writeunit", writeUnit(mdb)).Methods("POST", "OPTIONS")
	router.Handle("/api/readunits", readUnits(mdb)).Methods("GET", "OPTIONS")
	router.Handle("/api/readunit/", readUnit(mdb)).Methods("GET", "OPTIONS")
	//	http.Handle("/writeMaterial", UpdateData(mdb))
	//	http.Handle("/readMaterial", handleChart(mdb))
	//	http.Handle("/getDataArray", handlegetDataArray(mdb))
	//	http.Handle("/getDeviceList", handleDeviceList(mdb))
	//	http.Handle("/getParameters", handleGetParameters(mdb))

	http.ListenAndServe(":5000", router)
	for {

	}
}

func writeUnit(mdb dbservice.MDB) http.Handler {
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
		id, err := mdb.UpdateData("units", dataMap)
		if err == nil {
			sendAnswer202(w, strconv.Itoa(id))
		} else {
			sendAnswer400(w, err.Error())
		}

	})
}

func readUnits(mdb dbservice.MDB) http.Handler {
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
			fmt.Println("error reading querry:", err)
			sendAnswer400(w, err.Error())
			return
		}
		jsonString := string(jsonArray)
		sendAnswer200(w, jsonString)
	})
}

func readUnit(mdb dbservice.MDB) http.Handler {
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
		unit, err := mdb.ReadUnit(id)
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

func sendAnswer400(w http.ResponseWriter, answer string) {
	w.WriteHeader(400)
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(answer)))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, answer)
}
func sendAnswer404(w http.ResponseWriter, answer string) {
	w.WriteHeader(404)
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(answer)))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, answer)
}
func sendAnswer405(w http.ResponseWriter, answer string) {
	w.WriteHeader(405)
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(answer)))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, answer)
}
func sendAnswer200(w http.ResponseWriter, answer string) {
	//	w.WriteHeader(200)
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(answer)))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	fmt.Fprintf(w, answer)

}

func sendAnswer202(w http.ResponseWriter, answer string) {
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(answer)))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(202)
	fmt.Fprintf(w, answer)

}

//func handlegetLastData(mdb mdatabase.MDB) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		device_name := r.URL.Query().Get("device")
//		datetime := r.URL.Query().Get("datetime")
//		fmt.Println(datetime)
//		data := mdb.Get_last_data(device_name, datetime)
//		fmt.Println(data)
//		fmt.Fprintf(w, data)
//	})
//}

//func handlegetDataArray(mdb mdatabase.MDB) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		device_name := r.URL.Query().Get("device")
//		datetime1 := r.URL.Query().Get("datetime1")
//		datetime2 := r.URL.Query().Get("datetime2")
//		data := mdb.Get_data_array(device_name, datetime1, datetime2)
//
//		fmt.Fprintf(w, data)
//	})
//}
//func handleChart(mdb mdatabase.MDB) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		options := ""
//		device_list := mdb.Get_devices()
//		for _, device_name := range device_list {
//			options = options + fmt.Sprintf(optionHTML, device_name, device_name)
//		}
//		fmt.Println(options)
//		fmt.Fprintf(w, chartHTML, options)
//	})
//}

//func handleDeviceList(mdb mdatabase.MDB) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		device_list := mdb.Get_devices()
//		data, _ := json.Marshal(device_list)
//		fmt.Fprintf(w, string(data))
//	})
//}
//func handleGetParameters(mdb mdatabase.MDB) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		device_name := r.URL.Query().Get("device")
//		datetime1 := r.URL.Query().Get("datetime1")
//		datetime2 := r.URL.Query().Get("datetime2")
//		data := mdb.Get_parameters(device_name, datetime1, datetime2)
//		fmt.Fprintf(w, data)
//	})
//}
