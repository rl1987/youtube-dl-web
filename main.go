package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("public_html"))
	http.Handle("/", fs)

	log.Println("Listening for HTTP requests...")
	http.ListenAndServe("0.0.0.0:8000", nil)
}
