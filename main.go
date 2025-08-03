package main

import (
	"log"
	"os"

	"github.com/alex1234ak/go_final_project/pkg/db"
	"github.com/alex1234ak/go_final_project/pkg/server"

	_ "modernc.org/sqlite"
)

func main() {
	// Создаём подключение к БД через новую структуру
	database, err := db.NewDatabase(db.Name)
	if err != nil {
		logger := log.New(os.Stdout, "INFO: ", log.LstdFlags)
		logger.Println(err)
		return
	}
	defer database.Close()

	// Новый логгер и сервер
	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags)
	srv := server.NewServer(logger, database)

	// Запуск сервера
	if err := srv.Start(); err != nil {
		logger.Fatal("Error starting server: ", err)
	}
}
