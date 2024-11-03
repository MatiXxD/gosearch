package sitescan

import (
	"gosearch/pkg/crawler"
	"gosearch/pkg/crawler/spider"
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

	return filterWords(allPages, words), nil
}

func filterWords(pages []crawler.Document, words []string) map[string][]string {
	filtered := make(map[string][]string, len(words))
	for _, word := range words {
		slice := []string{}
		for _, page := range pages {
			if strings.Contains(strings.ToLower(page.Title), strings.ToLower(word)) {
				slice = append(slice, page.URL)
			}
		}
		filtered[word] = slice
	}
	return filtered
}
