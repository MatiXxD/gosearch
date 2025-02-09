package index

import (
	"gosearch/pkg/crawler"
	"strings"
)

type Service struct {
	Pages  []crawler.Document `json:"docs"`
	RevInd map[string][]int   `json:"index"`
}

func New() *Service {
	return &Service{
		RevInd: make(map[string][]int),
	}
}

func (s *Service) Add(docs []crawler.Document) {
	s.Pages = docs
	for _, doc := range docs {
		words := strings.Fields(doc.Title)
		for _, word := range words {
			word = strings.ToLower(word)
			if _, ok := s.RevInd[word]; ok {
				if !s.contains(word, doc.ID) {
					s.RevInd[word] = append(s.RevInd[word], doc.ID)
				}
			} else {
				s.RevInd[word] = []int{doc.ID}
			}
		}
	}
}

func (s *Service) Get(name string) []int {
	return s.RevInd[strings.ToLower(name)]
}

func (s *Service) FilterByWords(words []string) map[string][]string {
	filtered := make(map[string][]string, len(words))
	for _, word := range words {
		ids := s.Get(word)
		urls := []string{}
		for _, id := range ids {
			if url := crawler.FindPageByID(s.Pages, id); url != "" {
				urls = append(urls, crawler.FindPageByID(s.Pages, id))
			}
		}
		filtered[word] = urls
	}
	return filtered
}

func (s *Service) contains(word string, ind int) bool {
	for _, id := range s.RevInd[word] {
		if id == ind {
			return true
		}
	}
	return false
}
