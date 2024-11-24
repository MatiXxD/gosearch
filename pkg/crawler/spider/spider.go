package spider

import (
	"gosearch/pkg/crawler"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Scan(url string, depth int) ([]crawler.Document, error) {
	pages := make(map[string]string)
	err := parse(url, url, depth, pages)
	if err != nil {
		return nil, err
	}

	var data []crawler.Document
	for url, title := range pages {
		item := crawler.Document{
			URL:   url,
			Title: title,
		}
		data = append(data, item)
	}

	return data, nil
}

func parse(url, baseurl string, depth int, data map[string]string) error {
	if depth == 0 {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	pageRoot, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	data[url] = pageTitle(pageRoot)
	if depth == 1 {
		return nil
	}

	links := pageLinks(nil, pageRoot)
	for _, link := range links {
		link = strings.TrimSuffix(link, "/")
		if strings.HasPrefix(link, "/") && len(link) > 1 {
			link = baseurl + link
		}

		if data[link] != "" {
			continue
		}

		if strings.HasPrefix(link, baseurl) {
			err := parse(link, baseurl, depth-1, data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func pageTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
		return n.FirstChild.Data
	}

	var title string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title = pageTitle(c)
		if title != "" {
			break
		}
	}

	return title
}

func pageLinks(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				if !contains(links, a.Val) {
					links = append(links, a.Val)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = pageLinks(links, c)
	}

	return links
}

func contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
