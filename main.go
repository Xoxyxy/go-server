package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var Task string

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello,", Task)
}

type requestBody struct {
	Task string `json:"task"`
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	body := requestBody{}
	json.NewDecoder(r.Body).Decode(&body)
	Task = body.Task
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/get", GetHandler).Methods("GET")
	router.HandleFunc("/api/post", PostHandler).Methods("POST")

	http.ListenAndServe(":8080", router)
}
