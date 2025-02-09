package webapp

import (
	"gosearch/pkg/index"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Addr         string
	IndexService *index.Service
	Mux          *mux.Router
}

func New(addr string, ind *index.Service) *Server {
	return &Server{
		Addr:         addr,
		IndexService: ind,
		Mux:          mux.NewRouter(),
	}
}

func (s *Server) Run() error {
	server := http.Server{
		ReadTimeout:  100 * time.Second,
		WriteTimeout: 200 * time.Second,
		Addr:         s.Addr,
		Handler:      s.Mux,
	}

	listener, err := net.Listen("tcp4", s.Addr)
	if err != nil {
		log.Fatal(err)
	}

	s.BindRoutes()
	log.Println("Server running on: " + s.Addr)

	return server.Serve(listener)
}
