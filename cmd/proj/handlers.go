package main

import (
	"fmt"
	"net/http"
	"strconv"
	"errors"
	
	"displaybox.fisayoai.net/internal/models"
)

var tplHome = mustParseTemplates("base", "pages/home")
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	  quizzes, err := app.quizzes.Latest()
    if err != nil {
        http.Error(w, "Server Error", http.StatusInternalServerError)
        return
    }
    for _, quiz := range quizzes {
        fmt.Fprintf(w, "%+v\n", quiz)
    }
    // files := []strin
	app.logger.Info("Rendering home page")
	render(w, tplHome, nil)
}
var tplView = mustParseTemplates("base", "pages/view")
func (app *application) quizzesView(w http.ResponseWriter, r *http.Request){
	id,err :=strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	quiz, err := app.quizzes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Server Error", http.StatusInternalServerError)
		}
		return
	}
	err = tplView.Execute(w, quiz)
    if err != nil {
        app.logger.Error(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

func (app *application) quizzesCreatePost(w http.ResponseWriter, r *http.Request) {
    // Create some variables holding dummy data. We'll remove these later on
    // during the build.
    skill := "Python"
    quiz := "What is the modern way of declaring a variable in Python?\nA. var x = 10\nB. let x = 10"
    // Pass the data to the QuizModel.Insert() method, receiving the
    // ID of the new record back.
    id, err := app.quizzes.Insert(skill, quiz)
    if err != nil {
        http.Error(w, "Database Error", http.StatusInternalServerError)
        return
    }
    // Redirect the user to the relevant page for the quiz.
    http.Redirect(w, r, fmt.Sprintf("/quizzes/view/%d", id), http.StatusSeeOther)

}
var tplAbout = mustParseTemplates("base", "pages/about")
func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Rendering about page")
	render(w, tplAbout, nil)
}

var tplContact = mustParseTemplates("base", "pages/contact")
func (app *application) contact(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Rendering contact page")
	if r.Method == http.MethodGet {
		render(w, tplContact, nil)
		return
	}


	name := r.PostFormValue("name")
	msg := r.PostFormValue("message")

	if err := saveContactRequest(name, msg); err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Your message was saved.")
}