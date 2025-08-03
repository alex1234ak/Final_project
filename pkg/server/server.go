package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/alex1234ak/go_final_project/pkg/api"
	"github.com/alex1234ak/go_final_project/pkg/db"
	"github.com/alex1234ak/go_final_project/tests"
)

type Server struct {
	logger   *log.Logger
	port     string
	database *db.Database
}

func NewServer(logger *log.Logger, database *db.Database) *Server {
	port := strconv.Itoa(tests.Port)
	return &Server{
		logger:   logger,
		port:     port,
		database: database,
	}
}

func (s *Server) Start() error {
	api.Init(s.database)

	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)
	s.logger.Println("The server is running on http://localhost:" + s.port)
	return http.ListenAndServe(":"+s.port, nil)
}
