package main

import (
	"bootstrap/pkg/data"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
	IP    string
	Data  map[string]any
	Error string
	Flash string
	User  data.User
}

func (app *application) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) error {

	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplates, t),
		path.Join(pathToTemplates, "base.layout.gohtml"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	td.IP = app.ipFromContext(r.Context())

	td.Error = app.Session.PopString(r.Context(), "error")
	td.Flash = app.Session.PopString(r.Context(), "flash")

	err = parsedTemplate.Execute(w, td)
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

	if !app.authenticate(r, user, password) {
		app.Session.Put(r.Context(), "error", "Invalid login")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	_ = app.Session.RenewToken(r.Context())
	app.Session.Put(r.Context(), "flash", "Successfully logged in!")
	http.Redirect(w, r, "/u/p", http.StatusSeeOther)
	//return
}

func (app *application) authenticate(request *http.Request, user *data.User, password string) bool {

	if valid, err := user.PasswordMatches(password); err != nil || !valid {
		return false
	}

	app.Session.Put(request.Context(), "user", user)

	return true
}

func (app *application) UploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	files, err := app.UploadFile(r, "./static/img/")
	if err != nil || files == nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := app.Session.Get(r.Context(), "user").(data.User)

	var i = data.UserImage{
		UserID:   user.ID,
		FileName: files[0].OriginalFileName,
	}

	_, err = app.DB.InsertUserImage(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, err := app.DB.GetUser(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	app.Session.Put(r.Context(), "user", updatedUser)

	http.Redirect(w, r, "/u/p", http.StatusSeeOther)

}

type UploadedFile struct {
	OriginalFileName string
	FileSize         int64
}

func (app *application) UploadFile(r *http.Request, uploadDir string) ([]*UploadedFile, error) {
	var uploadedFiles []*UploadedFile

	// 5 MB
	if err := r.ParseMultipartForm(5 * 1024 * 1024); err != nil {
		return nil, fmt.Errorf("the upload file is too big: %w", err)
	}

	// Garante que o diretório existe (caminho relativo ao CWD do processo!)
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return nil, fmt.Errorf("cannot ensure upload dir: %w", err)
	}

	for _, fHeaders := range r.MultipartForm.File {
		for _, hdr := range fHeaders {
			var uploadedFile UploadedFile

			infile, err := hdr.Open()
			if err != nil {
				return nil, fmt.Errorf("cannot open uploaded part: %w", err)
			}
			defer infile.Close()

			// evita traversal: só o nome base
			uploadedFile.OriginalFileName = filepath.Base(hdr.Filename)
			destPath := filepath.Join(uploadDir, uploadedFile.OriginalFileName)

			out, err := os.Create(destPath)
			if err != nil {
				return nil, fmt.Errorf("cannot create destination file: %w", err)
			}
			defer out.Close()

			n, err := io.Copy(out, infile)
			if err != nil {
				return nil, fmt.Errorf("cannot write file: %w", err)
			}
			uploadedFile.FileSize = n

			uploadedFiles = append(uploadedFiles, &uploadedFile)
		}
	}

	return uploadedFiles, nil
}
