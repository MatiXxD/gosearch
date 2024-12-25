package main

import (
	"flag"
	"gosearch/internal/sitescan"
	"log"
	"os"
	"strings"
)

func main() {
	filename := flag.String("f", "", "Filename")
	word := flag.String("s", "", "Enable word search")
	flag.Parse()

	if *word == "" && *filename == "" {
		log.Fatal("For search use '-s' flag with words to search.")
	}

	var err error
	var isCached bool
	var file *os.File
	if *filename != "" {
		file, isCached, err = sitescan.PrepareFile(*filename, *word)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	var foundUrls map[string][]string
	if !isCached || *filename == "" {
		words := strings.Fields(*word)
		sites := []string{"https://go.dev", "https://google.com"}
		foundUrls, err = sitescan.ScanSites(sites, words, 2)
		if err != nil {
			log.Fatal(err)
		}
	}

	if isCached {
		_ = sitescan.ReadUrls(file, os.Stdout)
	} else {
		_ = sitescan.WriteUrls(foundUrls, os.Stdout)
	}

	if *filename != "" && !isCached {
		_ = sitescan.WriteUrls(foundUrls, file)
	}
}
