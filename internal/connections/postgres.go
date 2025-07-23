package connections

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
)

var PGDatabase *sql.DB

func init() {
	connString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	var err error
	PGDatabase, err = sql.Open("pgx", connString)
	if err != nil {
		log.Fatal(err)
	}
	errPing := PGDatabase.Ping()
	if errPing != nil {
		log.Fatal(errPing)
	}
	m, err := migrate.New(
		"file://internal/migrations",
		connString,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("Migrations applied!")

}
