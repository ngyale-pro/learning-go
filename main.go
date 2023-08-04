package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // Using the blank identifier, allow to add the postgre driver library without actually using it explicitly
	"github.com/ngyale-pro/simplebank/api"
	db "github.com/ngyale-pro/simplebank/db/sqlc"
	"github.com/ngyale-pro/simplebank/util"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load configuration.", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to DB.", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("can't initialize server.", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can't start server.", err)
	}

}
