package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type TaskRequest struct {
	Task   string `json:"task"`
	IsDone bool   `json:"isDone"`
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	DB.Find(&tasks)
	json.NewEncoder(w).Encode(tasks)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка JSON", http.StatusBadRequest)
		return
	}

	task := Task{Task: req.Task, IsDone: req.IsDone}
	DB.Create(&task)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func PatchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка JSON", http.StatusBadRequest)
		return
	}

	var task Task
	if result := DB.First(&task, id); result.Error != nil {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	if req.Task != "" {
		task.Task = req.Task
	}
	task.IsDone = req.IsDone

	DB.Save(&task)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var task Task
	if result := DB.First(&task, id); result.Error != nil {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	if result := DB.Delete(&task); result.Error != nil {
		http.Error(w, "Ошибка при удалении", http.StatusInternalServerError)
		return
	}
	//Возвращает задачу и была ошибка в нем почему в начале не удалял и просто задача шла
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Задача успешно удалена")
}

func main() {
	InitDB()
	DB.AutoMigrate(&Task{})

	router := mux.NewRouter()
	router.HandleFunc("/api/tasks", PostHandler).Methods("POST")
	router.HandleFunc("/api/tasks", GetHandler).Methods("GET")
	router.HandleFunc("/api/tasks/{id}", PatchHandler).Methods("PATCH")
	router.HandleFunc("/api/tasks/{id}", DeleteHandler).Methods("DELETE")
	// Также и здесь забыл поставить id в delete и просто не было ответа от json
	http.ListenAndServe(":8080", router)
}
