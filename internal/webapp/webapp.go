package webapp

import (
	"gosearch/internal/api"
	"gosearch/pkg/index"
	"log"
	"net"
	"net/http"
	"time"
)

type Server struct {
	Addr string
	API  *api.API
}

func New(addr string, ind *index.Service) *Server {
	return &Server{
		Addr: addr,
		API:  api.New(ind),
	}
}

func (s *Server) Run() error {
	server := http.Server{
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 200 * time.Second,
		Addr:         s.Addr,
		Handler:      s.API.Mux,
	}

	listener, err := net.Listen("tcp4", s.Addr)
	if err != nil {
		log.Fatal(err)
	}

	s.API.BindRoutes()
	log.Println("Server running on: " + s.Addr)

	return server.Serve(listener)
}
