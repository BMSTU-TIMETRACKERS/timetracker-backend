package time_tracker

import (
	"log"
	"net/http"
)

type Server struct {
	HttpServer *http.Server
}

func (s *Server) Start() error {
	log.Println("service starting at: ", s.HttpServer.Addr)
	return s.HttpServer.ListenAndServe()
}
