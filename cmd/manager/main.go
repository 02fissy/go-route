package main

import (
	"fmt"
	"log/slog"
	"os"

	"displaybox.fisayoai.net/internal/templates"
)

func main() {
	fmt.Println("Hello, World!")

	// go run cmd/manager/main.go
	// go build -o manager.exe cmd/manager/main.go

	templateCache, err := templates.NewTemplateCache()
    if err != nil {
        slog.Error(err.Error())
        os.Exit(1)
    }	
}