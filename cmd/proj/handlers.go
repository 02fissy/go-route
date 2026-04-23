package main

import (
	"fmt"
	"net/http"
	"os"
	
)


func home(w http.ResponseWriter, r *http.Request) {
	render(w, "./ui/html", Sources{Base: "base.tmpl", Page: "pages/home.tmpl"}, nil)
}

func about(w http.ResponseWriter, r *http.Request) {
	render(w, "./ui/html", Sources{Base: "base.tmpl", Page: "pages/about.tmpl"}, nil)
}

func contact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, "./ui/html", Sources{Base: "base.tmpl", Page: "pages/contact.tmpl"}, nil)
		return
	}


	name := r.PostFormValue("name")
	msg := r.PostFormValue("message")

	f, err := os.OpenFile("contacts.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("Name: %s | Message: %s\n", name, msg))
	fmt.Fprint(w, "Your message was saved.")
}