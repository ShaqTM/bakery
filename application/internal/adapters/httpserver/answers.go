package httpserver

import (
	"fmt"
	"net/http"
	"strconv"
)

func sendAnswer400(w http.ResponseWriter, answer string) {
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(answer)))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(400)
	fmt.Fprintf(w, answer)
}
func sendAnswer404(w http.ResponseWriter, answer string) {
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(answer)))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(404)
	fmt.Fprintf(w, answer)
}
func sendAnswer405(w http.ResponseWriter, answer string) {
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(answer)))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(405)
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
