package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var GlobalDatabase *sql.DB

func InitDatabase() {
	database, err := sql.Open("sqlite3", "file:./data/database/database.db?cache=shared&mode=memory")
	if err != nil {
		log.Fatal(err)
	}
	database.SetMaxOpenConns(1)
	GlobalDatabase = database
}
