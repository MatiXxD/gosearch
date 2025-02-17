package api

import (
	"gosearch/pkg/index"

	"github.com/gorilla/mux"
)

type API struct {
	IndexService *index.Service
	Mux          *mux.Router
}

func New(ind *index.Service) *API {
	return &API{
		IndexService: ind,
		Mux:          mux.NewRouter(),
	}
}
