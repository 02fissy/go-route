package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"displaybox.fisayoai.net/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
    logger *slog.Logger
	quizzes models.QuizModeler
	templateCache map[string]*template.Template
	sessionManager *scs.SessionManager
}
//using constructor dependency injection
func NewApplication(
    logger *slog.Logger, 
    quizzes *models.QuizModel, 
    templateCache map[string]*template.Template,
	sessionManager *scs.SessionManager,
) *application {
    
    return &application{
        logger:        logger,
        quizzes:       quizzes,
        templateCache: templateCache,
		sessionManager: sessionManager,
    }
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

	
	templateCache, err := newTemplateCacheV2()
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour


	  app := NewApplication(
        logger,
		&models.QuizModel{DB: db},
		templateCache,
		sessionManager,
	  )



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