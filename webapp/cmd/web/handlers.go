package main

import (
	"bootstrap/pkg/data"
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
func (app *application) Profile(response http.ResponseWriter, request *http.Request) {

	_ = app.render(response, request, "profile.page.gohtml", &TemplateData{})
}

type TemplateData struct {
	IP   string
	Data map[string]any
}


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
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	form := NewForm(r.PostForm)
	form.Required("email", "password")
	if !form.Valid() {
		app.Session.Put(r.Context(), "error", "Invalid login credential")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.DB.GetUserByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid login")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.Println(email, password, user.ID)

	_ = app.Session.RenewToken(r.Context())
	app.Session.Put(r.Context(), "flash", "Successfully logged in!")
	http.Redirect(w, r, "/u/p", http.StatusSeeOther)
	return
}
