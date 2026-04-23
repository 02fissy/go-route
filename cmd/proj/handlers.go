package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type tplSrcr interface {
	Layout() string
	Others() []string
}

type TplSrc string

func (c TplSrc) Layout() string {
	p := strings.Split(string(c), "|")
	return p[0]
}

func (c TplSrc) Others() []string {
	p := strings.Split(string(c), "|")
	if len(p) > 1 {
		return p[1:]
	}
	return []string{}
}

func parseTemplates(baseDir string, tplFiles tplSrcr) (*template.Template, error) {
	files := []string{filepath.Join(baseDir, tplFiles.Layout())}

	for _, path := range tplFiles.Others() {
		files = append(files, filepath.Join(baseDir, path))
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	ts, err = ts.ParseGlob(filepath.Join(baseDir, "subs/*.tmpl"))
	if err != nil {
		return nil, err
	}

	return ts, nil
}

type TemplateProvider interface {
	Render(w http.ResponseWriter, src tplSrcr, data any)
}

type FileTemplateLoader struct {
	Base string
}

func (f FileTemplateLoader) Render(w http.ResponseWriter, src tplSrcr, data any) {
	ts, err := parseTemplates(f.Base, src)
	if err != nil {
		http.Error(w, "Template Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Execution Error", http.StatusInternalServerError)
	}
}

type DataStore interface {
	Save(s string) error
}

type application struct {
	templates TemplateProvider
	store     DataStore
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.templates.Render(w, TplSrc("base.tmpl|pages/home.tmpl"), nil)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.templates.Render(w, TplSrc("base.tmpl|pages/about.tmpl"), nil)
}

func (app *application) contact(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.templates.Render(w, TplSrc("base.tmpl|pages/contact.tmpl"), nil)
		return
	}

	name := r.PostFormValue("name")
	msg := r.PostFormValue("message")
	fullMsg := fmt.Sprintf("Name: %s, Message: %s", name, msg)

	if err := app.store.Save(fullMsg); err != nil {
		http.Error(w, "Save failed", 500)
		return
	}
	fmt.Fprint(w, "Message Saved!")
}


type FileStore struct{}

func (f FileStore) Save(s string) error {
	file, err := os.OpenFile("contacts.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(s + "\n")
	return err
}
