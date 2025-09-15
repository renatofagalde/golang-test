package dbrepo

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "pass"
	dbName   = "users_test"
	port     = "5400"
	dns      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezon=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB

func TestMain(m *testing.M) {
	// connect to dokcer. fail if docker not running
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal("could not connect to docker; is it running? %s", err)
	}

	pool = p

	// set up our docker options, specifying the image and so forth
	opt := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14,5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWOR=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	//get a resource - docker image

	//start the image and wait until it's ready

	//populate the database with empty tables

	//run tests
	code := m.Run()

	os.Exit(code)
}
