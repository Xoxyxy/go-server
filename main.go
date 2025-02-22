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

type createTaskRequest struct {
	Text   string `json:"text"`
	IsDone bool   `json:"is_done"`
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	body := createTaskRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "JSON decoding error: "+err.Error(), http.StatusBadRequest)
		return
	}

	task := Task{
		Text:   body.Text,
		IsDone: body.IsDone,
	}

	if err := DB.Create(&task).Error; err != nil {
		http.Error(w, "Error creating task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

type updateTaskRequest struct {
	Text   *string `json:"text"`
	IsDone *bool   `json:"is_done"`
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	task := Task{}

	if err := DB.First(&task, id).Error; err != nil {
		http.Error(w, "Task not found: "+err.Error(), http.StatusNotFound)
		return
	}

	body := updateTaskRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "JSON decoding error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if body.Text != nil {
		task.Text = *body.Text
	}

	if body.IsDone != nil {
		task.IsDone = *body.IsDone
	}

	if err := DB.Save(&task).Error; err != nil {
		http.Error(w, "Error updating task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func main() {
	InitDB()
	DB.AutoMigrate(&Task{})

	router := mux.NewRouter()

	router.HandleFunc("/api/tasks", GetTasksHandler).Methods("GET")
	router.HandleFunc("/api/tasks", CreateTaskHandler).Methods("POST")
	router.HandleFunc("/api/tasks/{id}", UpdateTaskHandler).Methods("PATCH")

	http.ListenAndServe(":8080", router)
}
