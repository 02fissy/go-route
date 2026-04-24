package main

import (
	"fmt"
	"net/http"
)

var tplHome = mustParseTemplates("base", "pages/home")
func home(w http.ResponseWriter, r *http.Request) {
	render(w, tplHome, nil)
}

var tplAbout = mustParseTemplates("base", "pages/about")
func about(w http.ResponseWriter, r *http.Request) {
	render(w, tplAbout, nil)
}

var tplContact = mustParseTemplates("base", "pages/contact")
func contact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, tplContact, nil)
		return
	}


	name := r.PostFormValue("name")
	msg := r.PostFormValue("message")

	if err := saveContactRequest(name, msg); err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Your message was saved.")
}