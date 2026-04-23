package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /about", about)
	mux.HandleFunc(" /contact", contact)

	log.Print("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)

	log.Fatal(err)
}