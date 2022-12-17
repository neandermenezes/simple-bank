package db

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/postgres?sslmode=disable"
)

var MockQueries *Queries
var MockDB *sql.DB
var err error

func TestMain(m *testing.M) {
	MockDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	MockQueries = New(MockDB)

	m.Run()
}
