package index

import (
	"gosearch/pkg/crawler"
	"strings"
)

type Service struct {
	Pages  []crawler.Document
	revInd map[string][]int
}

func New() *Service {
	return &Service{
		revInd: make(map[string][]int),
	}
}

func (s *Service) Add(docs []crawler.Document) {
	s.Pages = docs
	for _, doc := range docs {
		words := strings.Fields(doc.Title)
		for _, word := range words {
			word = strings.ToLower(word)
			if _, ok := s.revInd[word]; ok {
				if !s.contains(word, doc.ID) {
					s.revInd[word] = append(s.revInd[word], doc.ID)
				}
			} else {
				s.revInd[word] = []int{doc.ID}
			}
		}
	}
}

func (s *Service) Get(name string) []int {
	return s.revInd[strings.ToLower(name)]
}

func (s *Service) contains(word string, ind int) bool {
	for _, id := range s.revInd[word] {
		if id == ind {
			return true
		}
	}
	return false
}
