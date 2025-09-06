package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

var pathToTemplates = "./templates/"

func (app *application) Home(response http.ResponseWriter, request *http.Request) {
	//tempalte date
	var td = make(map[string]any)

	if app.Session.Exists(request.Context(), "test") {
		msg := app.Session.Get(request.Context(), "test")
		td["test"] = msg
	} else {
		app.Session.Put(request.Context(), "test", "Hit this page at "+time.Now().UTC().String())
	}

	_ = app.render(response, request, "home.page.gohtml", &TemplateData{Data: td})
}

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (app *application) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {

	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplates, t),
		path.Join(pathToTemplates, "base.layout.gohtml"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	data.IP = app.ipFromContext(r.Context())

	err = parsedTemplate.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	form := NewForm(r.PostForm)
	form.Required("email", "password")
	if !form.Valid() {
		fmt.Fprint(w, "failed validation")
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.DB.GetUserByEmail(email)

	log.Println(email, password)

	fmt.Fprint(w, email)
}
