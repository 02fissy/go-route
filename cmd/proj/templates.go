package main

import (
	"path/filepath"
	"html/template"
	"net/http"
)

type tempSrc interface{
	Layout() string
	Others() []string
}

type Sources struct{
	Base string
	Page string 
}
func (t Sources) Layout() string {
	return t.Base
}

func (t Sources) Others() []string {
	return []string{t.Page}
}


func parseTemplates(baseDir string, tplFiles tempSrc) (*template.Template, error) {
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
func render(w http.ResponseWriter, baseDir string, tpl tempSrc, data any) {
	ts, err := parseTemplates(baseDir, tpl)
	if err != nil {
		http.Error(w, "Template Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Execution Error: "+err.Error(), http.StatusInternalServerError)
	}
}