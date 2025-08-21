package main

import (
	"html/template"
	"net/http"
)

func (app *application) Home(response http.ResponseWriter, request *http.Request) {
	_ = app.render(response, request, "home.page.gohtml", &TemplateData{})
}

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (app *application) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {

	parsedTemplate, err := template.ParseFiles("./templates/" + t)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	err = parsedTemplate.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
