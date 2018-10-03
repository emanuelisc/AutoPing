package main

import (
	"encoding/json"
	"log"
	"net/http"
	// "os"

	"github.com/rs/cors"
)

func main() {

	router := NewRouter()

	if err := http.ListenAndServe(":8080", cors.AllowAll().Handler(router)); err != nil {
		log.Fatal(err)
	}
}

func responseError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func responseJSONCode(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
