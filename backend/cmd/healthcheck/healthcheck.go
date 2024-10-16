package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide the server address (host:port) to test as a positional argument.")
	}
	serverAddr := os.Args[1]
	if serverAddr == "" {
		log.Fatal("Please provide a server address (host:port) as a positional argument")
	}

	if _, err := http.Get(serverAddr); err != nil {
		os.Exit(1)
	}
}
