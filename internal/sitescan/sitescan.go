package sitescan

import (
	"gosearch/pkg/crawler"
	"gosearch/pkg/crawler/spider"
	"gosearch/pkg/index"
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
