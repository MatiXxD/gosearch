package main

import (
	"flag"
	"fmt"
	"gosearch/internal/sitescan"
	"log"
	"strings"
)

func main() {
	wordsFlag := flag.String("s", "", "Enable word search")
	flag.Parse()

	if *wordsFlag == "" {
		log.Fatal("For search use '-s' flag with words to search.")
	}

	words := strings.Fields(*wordsFlag)
	sites := []string{"https://go.dev", "https://google.com"}

	foundUrls, err := sitescan.ScanSites(sites, words, 2)
	if err != nil {
		log.Fatal(err)
	}

	for word, urls := range foundUrls {
		for _, url := range urls {
			fmt.Printf("%s found in %s\n", word, url)
		}
	}
}
