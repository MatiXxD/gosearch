package main

import (
	"gosearch/internal/sitescan"
	"gosearch/internal/webapp"
	"gosearch/pkg/index"
	"log"
)

var urls []string = []string{"https://www.google.com/", "https://go.dev/"}

func main() {
	pages, err := sitescan.GetAllPages(urls, 2)
	if err != nil {
		log.Fatal(err)
	}

	ind := index.New()
	ind.AddMulti(pages)
	log.Println("Index created")

	s := webapp.New(":8091", ind)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
