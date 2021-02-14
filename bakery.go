package main

import (
	dbservice "bakery/dbService"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	mdb := dbservice.MDB{}
	mdb.InitDatabase()

	http.Handle("/api/writeunit", writeunit(mdb))
	//	http.Handle("/writeMaterial", UpdateData(mdb))
	//	http.Handle("/readMaterial", handleChart(mdb))
	//	http.Handle("/getDataArray", handlegetDataArray(mdb))
	//	http.Handle("/getDeviceList", handleDeviceList(mdb))
	//	http.Handle("/getParameters", handleGetParameters(mdb))

	http.ListenAndServe(":9000", nil)
	for {

	}
}

func writeunit(mdb dbservice.MDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var value map[string]interface{}

		mdb.UpdateData("units")
		options := ""
		device_list := mdb.Get_devices()
		for _, device_name := range device_list {
			options = options + fmt.Sprintf(optionHTML, device_name, device_name)
		}
		fmt.Println(options)
		fmt.Fprintf(w, rootHTML, options)
	})
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
