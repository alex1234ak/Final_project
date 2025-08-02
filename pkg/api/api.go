package api

import (
	"net/http"

	"github.com/alex1234ak/go_final_project/pkg/db"
)

const Layout = "20060102"

var database *db.Database

func Init(dbInstance *db.Database) {
	database = dbInstance
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", TaskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", doneTaskHandler)
}
