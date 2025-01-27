package netsrv

import (
	"bufio"
	"fmt"
	"gosearch/pkg/crawler"
	"gosearch/pkg/index"
	"io"
	"net"
	"strings"
	"time"
)

type Server struct {
	Port         string
	IndexService *index.Service
}

func New(port string, ind *index.Service) *Server {
	return &Server{
		Port:         port,
		IndexService: ind,
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp4", "0.0.0.0:"+s.Port)
	if err != nil {
		return err
	}

	fmt.Printf("Server listen on port %s\n", s.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	defer fmt.Println("Connection closed")
	fmt.Println("Connection opened")

	_ = conn.SetDeadline(time.Now().Add(time.Second * 15))

	r := bufio.NewReader(conn)
	for {
		msg, _, err := r.ReadLine()
		if err != nil {
			return
		}

		words := strings.Fields(string(msg))
		if len(words) < 1 {
			return
		}

		conn.Write([]byte(fmt.Sprintf("Searching for words: %v\n", words)))

		filtered := s.filterWords(words)
		if isEmpty(filtered) {
			conn.Write([]byte("Nothing was found\n"))
		}

		if err := s.writeUrls(filtered, conn); err != nil {
			return
		}

		conn.Write([]byte("\n"))
		_ = conn.SetDeadline(time.Now().Add(time.Second * 15))
	}
}

func (s *Server) filterWords(words []string) map[string][]string {
	filtered := make(map[string][]string, len(words))
	for _, word := range words {
		ids := s.IndexService.Get(word)
		urls := []string{}
		for _, id := range ids {
			if url := crawler.FindPageByID(s.IndexService.Pages, id); url != "" {
				urls = append(urls, crawler.FindPageByID(s.IndexService.Pages, id))
			}
		}
		filtered[word] = urls
	}
	return filtered
}

func (s *Server) writeUrls(foundUrls map[string][]string, w io.Writer) error {
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

func isEmpty(filtered map[string][]string) bool {
	for _, urls := range filtered {
		if len(urls) > 0 {
			return false
		}
	}
	return true
}
