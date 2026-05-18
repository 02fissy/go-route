package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"displaybox.fisayoai.net/internal/database"
	"displaybox.fisayoai.net/internal/models"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

type application struct {
    logger *slog.Logger
	quizzes models.QuizModeler
	templateCache map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {

	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", "displaybox.db", "SQlite database")
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
	sessionManager.Store = sqlite3store.New(db)
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

func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(database.MigrationsFS)

	if err:= goose.SetDialect("sqlite"); err != nil {
		return err
	}
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	return nil
}

func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("sqlite", dsn)
    if err != nil {
        return nil, err
    }
	err = runMigrations(db)
	if err != nil {
		db.Close()
		return nil, err
	}
	db.SetMaxOpenConns(1)
    err = db.Ping()
    if err != nil {
        db.Close()
        return nil, err
    }

	// schema, err := os.ReadFile("internal/database/quizzes.sql")
	// if err != nil {
	// 	db.Close()
	// 	return nil, err
	// }
	// _, err = db.Exec(string(schema))
	// if err != nil {
	// 	db.Close()
	// 	return nil, err
	// }

    return db, nil
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
