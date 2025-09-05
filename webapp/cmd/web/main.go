package main

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	Session *scs.SessionManager
}

func main() {

	app := application{}

	//get a session manager
	app.Session = getSession()

	//print out a message
	log.Println("statirng server on port 8080")

	//start the server
	err := http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}

}
