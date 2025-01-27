package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp4", "localhost:8091")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	defer fmt.Println("Connection closed")

	in := bufio.NewReader(os.Stdin)
	srv := bufio.NewReader(conn)
	for {
		fmt.Println("Enter words:")
		msg, _, err := in.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		_, err = conn.Write([]byte(string(msg) + "\n"))
		if err != nil {
			log.Fatal(err)
		}

		for {
			resp, _, err := srv.ReadLine()
			if err != nil {
				log.Fatal(err)
			} else if len(resp) == 0 {
				break
			}
			fmt.Println(string(resp))
		}
		fmt.Println("\n\n\n-----------------------------")
	}
}
