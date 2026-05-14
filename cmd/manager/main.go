package main

import (
	"fmt"
	"log/slog"
	"os"
	"net/http"

	"displaybox.fisayoai.net/internal/templates"
)

type application struct{
	templateCache map[string]*template.Template
}

func main() {
	fmt.Println("Hello, World!")

	// go run cmd/manager/main.go
	// go build -o manager.exe cmd/manager/main.go

	templateCache, err := templates.NewTemplateCache()
    if err != nil {
        slog.Error(err.Error())
        os.Exit(1)
    }	

	app := &application{
        templateCache: templateCache,
    }

    slog.Info("Starting server on :4000")
    err = http.ListenAndServe(":4000", nil)
    if err != nil {
        os.Exit(1)
    }
}

