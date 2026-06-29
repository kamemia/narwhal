package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	port   int
	router http.Handler
}

func NewServer(port int, router http.Handler) *Server {
	return &Server{port: port, router: router}
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.router,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("narwhal running on :%d", s.port)
	return srv.ListenAndServe()
}
