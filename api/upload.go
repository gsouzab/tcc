package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GetUploads retorna todos os uploads de experimentos
func GetUploads(w http.ResponseWriter, r *http.Request) {
	// // response := RedisGetAllSensors()

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)
}

// GetUpload retorna um arquivo
func GetUpload(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uploadId := params["id"]

	response, err := os.Open("uploads/" + uploadId)
	
	defer response.Close()
	// if we os.Open returns an error then handle it
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} 
}