package api

import (
	"encoding/json"
	"gosearch/pkg/crawler"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func (a *API) GetDocs(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(a.IndexService.Pages); err != nil {
		log.Printf("ERROR: GetDocs: %v\n", err)
		http.Error(w, "Can't encode docs", http.StatusInternalServerError)
		return
	}
}

func (a *API) FilterByWords(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		http.Error(w, "Wrong content type", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Can't read request body", http.StatusBadRequest)
		return
	}

	words := strings.Fields(string(body))
	data := a.IndexService.FilterByWords(words)
	if isEmpty(data) {
		http.Error(w, "Nothing was found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("ERROR: FilterByWords: %v\n", err)
		http.Error(w, "Can't encode filtered data", http.StatusInternalServerError)
		return
	}
}

func (a *API) CreateDocument(w http.ResponseWriter, r *http.Request) {
	var doc crawler.Document
	err := json.NewDecoder(r.Body).Decode(&doc)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Can't decode data", http.StatusUnprocessableEntity)
		return
	}

	indDoc := a.IndexService.AddDoc(doc)
	if err := json.NewEncoder(w).Encode(indDoc); err != nil {
		http.Error(w, "Can't encode data", http.StatusInternalServerError)
		return
	}
}

func (a *API) GetDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Wrong document id", http.StatusBadRequest)
		return
	}

	doc, ok := a.IndexService.GetDoc(id)
	if !ok {
		http.Error(w, "Can't find doc", http.StatusNotFound)
	}

	if err := json.NewEncoder(w).Encode(doc); err != nil {
		http.Error(w, "Can't encode data", http.StatusInternalServerError)
		return
	}
}

func (a *API) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Wrong document id", http.StatusBadRequest)
		return
	}

	if err := a.IndexService.DeleteDoc(id); err != nil {
		http.Error(w, "Can't find doc", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Wrong document id", http.StatusBadRequest)
		return
	}

	var doc crawler.Document
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		http.Error(w, "Can't decode data", http.StatusUnprocessableEntity)
		return
	}

	if err := a.IndexService.UpdateDoc(id, doc); err != nil {
		http.Error(w, "Can't find doc", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func isEmpty(filtered map[string][]string) bool {
	for _, urls := range filtered {
		if len(urls) > 0 {
			return false
		}
	}
	return true
}
