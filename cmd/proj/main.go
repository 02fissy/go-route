package main

import (
	"log"
	"net/http"
)

func main() {

	app := &application{
		templates: FileTemplateLoader{Base: "./ui/html"},
		store:     FileStore{},
	}
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /about", app.about)
	mux.HandleFunc("/contact", app.contact)

	log.Print("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)

	log.Fatal(err)
}