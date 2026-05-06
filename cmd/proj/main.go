package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"displaybox.fisayoai.net/internal/models"
	_ "github.com/go-sql-driver/mysql"
)
type application struct {
    logger *slog.Logger
	quizzes *models.QuizModel
}
func main() {

	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", "displaybox_web:0204@/displaybox?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger:= slog.New(slog.NewTextHandler(os.Stdout, nil))

	 // from the command-line flag.
    db, err := openDB(*dsn)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }
	defer db.Close()

	  app := &application{
        logger: logger,
		quizzes: &models.QuizModel{DB: db},
	  }



	logger.Info("Starting server on ", "addr", *addr)
	err = http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
}
func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    err = db.Ping()
    if err != nil {
        db.Close()
        return nil, err
    }
    return db, nil
}