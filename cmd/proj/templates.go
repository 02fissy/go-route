package main

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"fmt"

	"displaybox.fisayoai.net/internal/models"

)

type templateData struct{
	Quizz models.Quiz
	Quizzes []models.Quiz
	Form any
}

// const tplDir = "./ui/html"
// const tplExt = ".tmpl"

// type tplSrc interface{
// 	Layout() string
// 	Others() []string
// }

// func mustParseSet(tplSet TemplateSet) *template.Template {
// 	tpl, err := parseTemplates(tplDir, tplExt, tplSet)
// 	if err != nil {
// 		panic(fmt.Sprintf("Failed to parse templateSet (%q): %s", tplSet.Others()[0], err.Error()))
// 	}
// 	return tpl
// }

// func mustParseTemplates(base string, others ...string) *template.Template {
// 	return mustParseSet(newSet(base, others...))
// }

// func parseTemplates(baseDir, tplExt string, tplFiles tplSrc) (*template.Template, error) {
// 	files := append([]string{tplFiles.Layout()}, tplFiles.Others()...)
// 	files = normalize(baseDir, tplExt, files...)
// 	// the two lines above can be replaced with a single line (below). That said, the two lines version is better for readability.
// 	// files := normalize(baseDir, tplExt, append([]string{tplFiles.Layout()}, tplFiles.Others()...)...)

// 	ts, err := template.ParseFiles(files...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ts, err = ts.ParseGlob(filepath.Join(baseDir, fmt.Sprint("subs/*", tplExt)))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return ts, nil
// }

// func normalize(baseDir, ext string, files ...string) []string {

// 	for i, f := range files {
// 		if filepath.Ext(f) == "" {
// 			files[i] += tplExt
// 		}
// 		files[i] = filepath.Join(baseDir, files[i])
// 	}

// 	return files
// }

func newTemplateCache() (map[string]*template.Template, error) {
cache := make(map[string]*template.Template)

pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
if err != nil {
	return nil, err
}

partials, err := filepath.Glob("./ui/html/subs/*.tmpl")
    if err != nil {
        return nil, err
    }
for _, page := range pages {
	name := filepath.Base(page)

	files := []string{"./ui/html/base.tmpl"}
        
	files = append(files, partials...)
	
	files = append(files, page)
	ts, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	cache[name] = ts
}

return cache, nil
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data any) error {

	ts, ok := app.templateCache[page]
    if !ok {
        err := fmt.Errorf("the template %s does not exist", page)
        app.logger.Error(err.Error())
        return err
    }

    buf := new(bytes.Buffer)
    err := ts.ExecuteTemplate(buf, "base", data)
    if err != nil {
        return err 
    }

    w.WriteHeader(status)
    buf.WriteTo(w)

    return nil
}

func handleErr(w http.ResponseWriter, err error) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	slog.Error("Error: %s\n", err.Error())
}

