package sitescan

import (
	"bufio"
	"fmt"
	"gosearch/pkg/crawler"
	"gosearch/pkg/crawler/spider"
	"gosearch/pkg/index"
	"io"
	"os"
	"strings"
)

func ScanSites(urls, words []string, depth int) (map[string][]string, error) {
	scanner := spider.New()

	var allPages []crawler.Document
	for _, url := range urls {
		scanned, err := scanner.Scan(url, depth)
		if err != nil {
			return nil, err
		}
		allPages = append(allPages, scanned...)
	}

	for i := range allPages {
		allPages[i].ID = i
	}

	return filterWords(allPages, words), nil
}

func PrepareFile(filename string, word string) (*os.File, bool, error) {
	isCached := true
	for _, w := range strings.Fields(word) {
		filename += "_" + w
	}
	filename += ".txt"

	var err error
	var file *os.File
	if !fileExists(filename) {
		file, err = os.Create(filename)
		isCached = false
		if err != nil {
			return nil, false, err
		}
	} else {
		file, err = os.OpenFile(filename, os.O_CREATE, 0666)
		if err != nil {
			return nil, false, err
		}
	}

	return file, isCached, nil
}

func WriteUrls(foundUrls map[string][]string, w io.Writer) error {
	for word, urls := range foundUrls {
		for _, url := range urls {
			s := fmt.Sprintf("%s found in %s\n", word, url)
			_, err := w.Write([]byte(s))
			if err != nil {
				return fmt.Errorf("Error while writing: %v\n", err)
			}
		}
	}
	return nil
}

func ReadUrls(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		_, err := w.Write([]byte(scanner.Text() + "\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func filterWords(pages []crawler.Document, words []string) map[string][]string {
	ind := index.New()
	ind.Add(pages)

	filtered := make(map[string][]string, len(words))
	for _, word := range words {
		ids := ind.Get(word)
		urls := []string{}
		for _, id := range ids {
			if url := binarySearch(pages, id); url != "" {
				urls = append(urls, binarySearch(pages, id))
			}
		}
		filtered[word] = urls
	}
	return filtered
}

func binarySearch(pages []crawler.Document, id int) string {
	i, j := 0, len(pages)-1
	for i <= j {
		mid := (i + j) / 2
		if pages[mid].ID == id {
			return pages[mid].URL
		} else if pages[mid].ID < id {
			i = mid + 1
		} else {
			j = mid - 1
		}
	}
	return ""
}
