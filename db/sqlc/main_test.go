package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" // Using the blank identifier, allow to add the postgre driver library without actually using it explicitly
	"github.com/ngyale-pro/simplebank/util"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../../.")
	if err != nil {
		log.Fatal("Can't load configuration.", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can't connect to DB")
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
