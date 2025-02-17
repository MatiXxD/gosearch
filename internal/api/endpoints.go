package api

import (
	"net/http"
)

func (a *API) BindRoutes() {
	a.Mux.HandleFunc("/api/v1/index", a.GetIndex).Methods(http.MethodGet)

	a.Mux.HandleFunc("/api/v1/docs", a.GetDocs).Methods(http.MethodGet)
	a.Mux.HandleFunc("/api/v1/docs/{id}", a.GetDocument).Methods(http.MethodGet)
	a.Mux.HandleFunc("/api/v1/docs", a.CreateDocument).Methods(http.MethodPost)
	a.Mux.HandleFunc("/api/v1/docs/filter", a.FilterByWords).Methods(http.MethodPost)
}
