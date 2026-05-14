package main

import (
	"fmt"
	"net/http"
	"strconv"
	"errors"
	
	"displaybox.fisayoai.net/internal/models"
	"displaybox.fisayoai.net/internal/validator"
)
type quizCreateForm struct {
	Skill string 
	Quiz string 
	validator.Validator
}

//var tplHome = mustParseTemplates("base", "pages/home")
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	quizzes, err := app.quizzes.Latest()
    if err != nil {
        http.Error(w, "Server Error", http.StatusInternalServerError)
        return
    }
    data := templateData{
		Quizzes: quizzes,
	}
	err = app.render(w, r, http.StatusOK, "home.tmpl", data)
    if err != nil {
        app.logger.Error(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
//var tplCreate = mustParseTemplates("base", "pages/create")
func (app *application) quizzesCreate(w http.ResponseWriter, r *http.Request){
	err := app.render(w, r, http.StatusOK, "create.tmpl", nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
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
	data := templateData{
		Quizz: quiz,
	}
	err = app.render(w, r, http.StatusOK, "view.tmpl", data)
    if err != nil {
        app.logger.Error(err.Error())
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

func (app *application) quizzesCreatePost(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
	if  err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	form := quizCreateForm{
		Skill: r.PostForm.Get("skill"),
		Quiz: r.PostForm.Get("quiz"),
	}
	form.CheckField(validator.NotBlank(form.Skill), "skill", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Quiz), "quiz", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Skill, 100), "skill", "This field cannot be more than 100 characters long")
	form.CheckField(validator.MaxChars(form.Quiz, 1000), "quiz", "This field cannot be more than 1000 characters long")

	

	if !form.Valid() {
		data := templateData{
			Form: form,
		}
		err = app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		if err != nil {
			app.logger.Error(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	id, err := app.quizzes.Insert(form.Skill, form.Quiz)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Quiz successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/quizzes/view/%d", id), http.StatusSeeOther)
}
func (app *application) about(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Rendering about page")
	app.render(w, r, http.StatusOK, "about.tmpl", nil)
}


func (app *application) contact(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Rendering contact page")
	if r.Method == http.MethodGet {
		app.render(w, r, http.StatusOK, "contact.tmpl", nil)
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