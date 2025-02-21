package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []Task

	if err := DB.Find(&tasks).Error; err != nil {
		http.Error(w, "Error getting tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

type requestBody struct {
	Text   string `json:"text"`
	IsDone bool   `json:"is_done"`
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := requestBody{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "JSON decoding error: "+err.Error(), http.StatusBadRequest)
	}

	task := Task{
		Text:   body.Text,
		IsDone: body.IsDone,
	}

	if err := DB.Create(&task).Error; err != nil {
		http.Error(w, "Error creating task:"+err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func main() {
	InitDB()
	DB.AutoMigrate(&Task{})

	router := mux.NewRouter()

	router.HandleFunc("/api/tasks", GetTasksHandler).Methods("GET")
	router.HandleFunc("/api/tasks", CreateTaskHandler).Methods("POST")

	http.ListenAndServe(":8080", router)
}
