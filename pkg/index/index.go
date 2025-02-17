package index

import (
	"gosearch/pkg/crawler"
	"strings"
)

type Service struct {
	Pages    []crawler.Document `json:"docs"`
	RevInd   map[string][]int   `json:"index"`
	curIndex int
}

func New() *Service {
	return &Service{
		Pages:    []crawler.Document{},
		RevInd:   make(map[string][]int),
		curIndex: 0,
	}
}

func (s *Service) AddDoc(doc crawler.Document) crawler.Document {
	doc.ID = s.curIndex
	s.curIndex++

	s.Pages = append(s.Pages, doc)
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

	return doc
}

func (s *Service) GetDoc(id int) (crawler.Document, bool) {
	for _, p := range s.Pages {
		if p.ID == id {
			return p, true
		}
	}
	return crawler.Document{}, false
}

func (s *Service) AddMulti(docs []crawler.Document) {
	for _, doc := range docs {
		s.AddDoc(doc)
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
