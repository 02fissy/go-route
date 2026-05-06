package main
import "net/http"
// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
mux := http.NewServeMux()

fileServer := http.FileServer(http.Dir("./ui/static/"))
mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
mux.HandleFunc("GET /{$}", app.home)
mux.HandleFunc("GET /about", app.about)
mux.HandleFunc("GET /contact", app.contact)
mux.HandleFunc("GET /quizzes/create", app.quizzesCreate)
mux.HandleFunc("GET /quizzes/view/{id}", app.quizzesView)
mux.HandleFunc("POST /quizzes/create", app.quizzesCreatePost)


return app.logRequest(commonHeaders(mux))
}