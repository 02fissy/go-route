package main
import (
	"github.com/justinas/alice"
	"net/http"
)
// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
mux := http.NewServeMux()

 dynamic := alice.New(app.sessionManager.LoadAndSave)

fileServer := http.FileServer(http.Dir("./ui/static/"))
mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
mux.Handle("GET /about", dynamic.ThenFunc(app.about))
mux.Handle("GET /contact", dynamic.ThenFunc(app.contact))
mux.Handle("GET /quizzes/create", dynamic.ThenFunc(app.quizzesCreate))
mux.Handle("GET /quizzes/view/{id}", dynamic.ThenFunc(app.quizzesView))
mux.Handle("POST /quizzes/create", dynamic.ThenFunc(app.quizzesCreatePost))

standard := alice.New(app.logRequest, commonHeaders)
return standard.Then(mux)
}