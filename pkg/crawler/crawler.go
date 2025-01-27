package crawler

type Interface interface {
	Scan(url string, depth int) ([]Document, error)
	// BatchScan(urls []string, depth int, workers int) (<-chan Document, <-chan error)
}

type Document struct {
	ID    int
	URL   string
	Title string
	// Body  string
}

func FindPageByID(pages []Document, id int) string {
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
