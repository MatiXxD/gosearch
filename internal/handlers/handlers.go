package handlers

import (
	"encoding/json"
	"gosearch/pkg/index"
	"io"
	"log"
	"net/http"
	"strings"
)

type Handlers struct {
	ind *index.Service
}

func NewHandlers(ind *index.Service) *Handlers {
	return &Handlers{ind}
}

func (h *Handlers) GetIndex(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(h.ind.RevInd); err != nil {
		log.Printf("ERROR: GetIndex: %v\n", err)
		http.Error(w, "Can't encode index", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetDocs(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(h.ind.Pages); err != nil {
		log.Printf("ERROR: GetDocs: %v\n", err)
		http.Error(w, "Can't encode docs", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) FilterByWords(w http.ResponseWriter, r *http.Request) {
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
	data := h.ind.FilterByWords(words)
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

func isEmpty(filtered map[string][]string) bool {
	for _, urls := range filtered {
		if len(urls) > 0 {
			return false
		}
	}
	return true
}
