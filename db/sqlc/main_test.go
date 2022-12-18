package db

import (
	"database/sql"
	"log"
	"simplebank/db/utils"
	"testing"

	_ "github.com/lib/pq"
)

var MockQueries *Queries
var MockDB *sql.DB
var err error

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	MockDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	MockQueries = New(MockDB)

	m.Run()
}
