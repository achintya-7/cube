package api

import (
	"fmt"
	"log"

	"github.com/achintya-7/cube/worker"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router  *gin.Engine
	address string
	port    int
	worker  *worker.Worker
}

func NewServer(address string, port int, worker *worker.Worker) *Server {
	return &Server{
		router:  gin.Default(),
		address: address,
		port:    port,
		worker:  worker,
	}
}

func (s *Server) setupRoutes() {
	s.router.GET("/tasks", s.getTasks)
	s.router.POST("/tasks", s.startTask)
	s.router.DELETE("/tasks/:id", s.stopTask)
}

func (s *Server) Start() {
	s.setupRoutes()

	log.Println("Starting server on", s.address+":"+fmt.Sprint(s.port))
	if err := s.router.Run(s.address + ":" + fmt.Sprint(s.port)); err != nil {
		panic(err)
	}
}
