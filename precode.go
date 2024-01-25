package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта

// 1. Обработчик для получения всех задач
func getTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("failed to write responce", err.Error())
		return
	}
}

// 2. Обработчик для отправки задачи на сервер
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// 3. Обработчик для получения задачи по ID
func getTaskForId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp) // не понимаю почему ругается на resp
	if err != nil {
		fmt.Printf("Ошибка записи: %v", err)
	}
	json.NewEncoder(w).Encode(task) // в данной обработке вообще запутался
}

// 4. Обработчик для удаления задачи по ID
func taskDel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Не удалось удалить задачу. Отсутствует запрашиваемая задача.", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	//1. Для получения всех задач:
	r.Get("/tasks", getTasks)
	//2. Для отправки задачи на сервер:
	r.Post("/tasks", postTask)
	//3. для получения задачи по ID:
	r.Get("/tasks/{id}", getTaskForId)
	//4. для удаления задачи по ID:
	r.Delete("/tasks/{id}", taskDel)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
