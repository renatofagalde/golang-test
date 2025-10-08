package dbrepo

import (
	"bootstrap/pkg/data"
	"bootstrap/pkg/repository"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

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
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB
var testRepository repository.DatabaseRepository

func TestMain(m *testing.M) {
	// connect to dokcer. fail if docker not running
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}

	pool = p

	// set up our docker options, specifying the image and so forth
	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "18.0",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"}, //internal port for docker
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port}},
		},
	}

	//get a resource - docker image
	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource %s", err)
	}

	//start the image and wait until it's ready
	if err := pool.Retry(func() error {
		var err error

		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println(err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database %s", err)
	}

	//populate the database with empty tables
	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	testRepository = &PostgresDBRepo{DB: testDB}

	//run tests
	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource %s", err)
	}

	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/users.sql")
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Test_pingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("can't ping database")
	}
}

func TestPostgresDBRepositoryInsertUser(t *testing.T) {
	testUser := data.User{FirstName: "Admin", LastName: "User", Email: "admin@example.com", Password: "secret", IsAdmin: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	id, err := testRepository.InsertUser(testUser)
	if err != nil {
		t.Errorf("insert user returned an error %s", err)
	}

	if id != 1 {
		t.Errorf("insert user returned wrong id, expected 1 but got %d", id)
	}

}

func TestPostgresDBRepositorySelectAllUser(t *testing.T) {
	allUsers, err := testRepository.AllUsers()
	if err != nil {
		t.Errorf("all users reports an error %s", err)
	}

	if len(allUsers) != 1 {
		t.Errorf("all users reports wrong size; expected 1, but got %d", len(allUsers))
	}

	testUser := data.User{FirstName: "test", LastName: "test", Email: "test@example.com", Password: "secret", IsAdmin: 0, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	_, _ = testRepository.InsertUser(testUser)

	allUsers, err = testRepository.AllUsers()
	if err != nil {
		t.Errorf("all users reports an error %s", err)
	}

	if len(allUsers) != 2 {
		t.Errorf("all users reports wrong size; expected 1, but got %d", len(allUsers))
	}

}
func TestPostgresDBRepositorySelectUserByEmail(t *testing.T) {
	var email string = "test@example.com"
	user, err := testRepository.GetUserByEmail(email)
	if err != nil {
		t.Errorf("all users reports an error %s", err)
	}

	if user.Email != "test@example.com" {
		t.Errorf("wrong email returned by GetUserByEmail; expected %s but got %s", email, user.Email)
	}

}
