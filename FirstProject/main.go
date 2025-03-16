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
	fmt.Fprintln(w, "Task обновлен")

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

	// Обновляем поля задачи, если они переданы в запросе
	if req.Task != "" {
		task.Task = req.Task
	}
	task.IsDone = req.IsDone

	DB.Save(&task)

	// Возвращаем обновленную сущность в виде JSON
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

	DB.Delete(&task)

	// Возвращаем сообщение об успешном удалении
	fmt.Fprintln(w, "Задача успешно удалена")
}
func main() {
	InitDB()
	DB.AutoMigrate(&Task{})

	router := mux.NewRouter()
	router.HandleFunc("/api/tasks", PostHandler).Methods("POST")
	router.HandleFunc("/api/tasks", GetHandler).Methods("GET")
	router.HandleFunc("/api/tasks/{id}", PatchHandler).Methods("PATCH")
	router.HandleFunc("/api/tasks/", DeleteHandler).Methods("DELETE")
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
