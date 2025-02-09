package webapp

import (
	"gosearch/internal/handlers"
	"net/http"
)

func (s *Server) BindRoutes() {
	h := handlers.NewHandlers(s.IndexService)
	s.Mux.HandleFunc("/index", h.GetIndex).Methods(http.MethodGet)
	s.Mux.HandleFunc("/docs", h.GetDocs).Methods(http.MethodGet)
	s.Mux.HandleFunc("/filter", h.FilterByWords).Methods(http.MethodPost)
}
