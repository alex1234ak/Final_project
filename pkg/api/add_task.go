package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alex1234ak/go_final_project/pkg/db"
)

// ok
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var err error
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONerror(w, fmt.Sprintf("incorrect data in JSON: %s", err), http.StatusBadRequest)
		return
	}
	if task.Title == "" {
		writeJSONerror(w, "title required", http.StatusBadRequest)
		return
	}
	if err := checkDate(&task); err != nil {
		writeJSONerror(w, fmt.Sprintf("Check date error: %s", err), http.StatusBadRequest)
		return
	}
	id, err := database.AddTask(&task)
	if err != nil {
		writeJSONerror(w, fmt.Sprintf("DB add task error: + %s", err), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]any{"id": id})
}

// OK
func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(value)
}

func writeJSONerror(w http.ResponseWriter, err string, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}

// OK
func checkDate(task *db.Task) error {

	now := time.Now().UTC()
	now = now.Truncate(24 * time.Hour)

	// Если дата не указана, устанавливаем сегодняшнюю
	if task.Date == "" {
		task.Date = now.UTC().Format(Layout)
	}
	// парсим дату
	date, err := time.Parse(Layout, task.Date)
	if err != nil {
		return err
	}

	var next string

	// Если правило повторения указано, вычисляем следующую дату
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// Проверка, что дата больше чем сейчас
	if afterNow(now, date) {
		if task.Repeat == "" {
			task.Date = now.Format(Layout)
		} else {
			task.Date = next
		}
	}

	return nil
}
