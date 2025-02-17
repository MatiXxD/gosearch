package sitescan

import (
	"gosearch/pkg/crawler"
	"gosearch/pkg/crawler/spider"
)

func GetAllPages(urls []string, depth int) ([]crawler.Document, error) {
	scanner := spider.New()

	var allPages []crawler.Document
	for _, url := range urls {
		scanned, err := scanner.Scan(url, depth)
		if err != nil {
			return nil, err
		}
		allPages = append(allPages, scanned...)
	}

	return allPages, nil
}
