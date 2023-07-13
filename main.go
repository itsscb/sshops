package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/enable", enableHandler)
	http.HandleFunc("/disable", disableHandler)
	http.HandleFunc("/state", stateHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	log.Println("Server listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
