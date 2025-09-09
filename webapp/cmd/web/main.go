package main

import (
	"bootstrap/pkg/data"
	"bootstrap/pkg/repository"
	"bootstrap/pkg/repository/dbrepo"
	"encoding/gob"
	"flag"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	DSN     string
	DB      repository.DatabaseRepository
	Session *scs.SessionManager
}

func main() {

	gob.Register(data.User{})

	app := application{}

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5490 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	app.DB.AllUsers()

	//get a session manager
	app.Session = getSession()

	//print out a message
	log.Println("statirng server on port 8080")

	//start the server
	err = http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}

}
