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
	Form any
}

const tplDir = "./ui/html"
const tplExt = ".tmpl"

type tplSrc interface{
	Layout() string
	Others() []string
}

func mustParseTemplates(base string, pages ...string) *template.Template {

	files := append([]string{base}, pages...)

	files = normalize(tplDir, tplExt, files...)

	ts, err := template.ParseFiles(files...)
	if err != nil {
		panic(fmt.Sprintf(
			"failed to parse templates (%q): %s",
			pages[0],
			err.Error(),
		))
	}

	ts, err = ts.ParseGlob(
		filepath.Join(tplDir, "subs", "*"+tplExt),
	)
	if err != nil {
		panic(err)
	}


	return ts
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


// func NewTemplateCache() (map[string]*template.Template, error) {
// 	cache := make(map[string]*template.Template)

// 	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
// 	if err != nil {
// 		return nil, err
// 	}

// 	partials, err := filepath.Glob("./ui/html/subs/*.tmpl")
// 		if err != nil {
// 			return nil, err
// 		}
// 	for _, page := range pages {
// 		name := filepath.Base(page)

// 		files := []string{"./ui/html/base.tmpl"}
			
// 		files = append(files, partials...)
		
// 		files = append(files, page)
// 		ts, err := template.ParseFiles(files...)
// 		if err != nil {
// 			return nil, err
// 		}

// 		cache[name] = ts
// 	}

// 	return cache, nil
// }


func newTemplateCacheV2() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	cache["home.tmpl"] = mustParseTemplates("base", "pages/home")
	cache["create.tmpl"] = mustParseTemplates("base", "pages/create")
	cache["view.tmpl"] = mustParseTemplates("base", "pages/view")
	cache["contact.tmpl"] = mustParseTemplates("base", "pages/contact")
	cache["about.tmpl"] = mustParseTemplates("base", "pages/about")


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

