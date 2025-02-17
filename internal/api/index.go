package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func (a *API) GetIndex(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(a.IndexService.RevInd); err != nil {
		log.Printf("ERROR: GetIndex: %v\n", err)
		http.Error(w, "Can't encode index", http.StatusInternalServerError)
		return
	}
}

func (a *API) NewIndex(w http.ResponseWriter, r *http.Request) {
	a.IndexService.RecreateIndex()
}
