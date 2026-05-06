package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"

	"displaybox.fisayoai.net/internal/models"

)

type templateData struct{
	Quizz models.Quiz
	Quizzes []models.Quiz
}

const tplDir = "./ui/html"
const tplExt = ".tmpl"

type tplSrc interface{
	Layout() string
	Others() []string
}

func mustParseSet(tplSet TemplateSet) *template.Template {
	tpl, err := parseTemplates(tplDir, tplExt, tplSet)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse templateSet (%q): %s", tplSet.Others()[0], err.Error()))
	}
	return tpl
}

func mustParseTemplates(base string, others ...string) *template.Template {
	return mustParseSet(newSet(base, others...))
}

func parseTemplates(baseDir, tplExt string, tplFiles tplSrc) (*template.Template, error) {
	files := append([]string{tplFiles.Layout()}, tplFiles.Others()...)
	files = normalize(baseDir, tplExt, files...)
	// the two lines above can be replaced with a single line (below). That said, the two lines version is better for readability.
	// files := normalize(baseDir, tplExt, append([]string{tplFiles.Layout()}, tplFiles.Others()...)...)

	ts, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	ts, err = ts.ParseGlob(filepath.Join(baseDir, fmt.Sprint("subs/*", tplExt)))
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func normalize(baseDir, ext string, files ...string) []string {

	for i, f := range files {
		if filepath.Ext(f) == "" {
			files[i] += tplExt
		}
		files[i] = filepath.Join(baseDir, files[i])
	}

	return files
}

func render(w http.ResponseWriter, tpl *template.Template, data any) error {

	buff := new(bytes.Buffer)
	if err := tpl.Execute(buff, data); err != nil {
		handleErr(w, err)
		return err
	}
	
	_, err := buff.WriteTo(w)
	if err != nil {
		handleErr(w, err)
		return err 
	}

	return nil
}

func handleErr(w http.ResponseWriter, err error) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	slog.Error("Error: %s\n", err.Error())
}

