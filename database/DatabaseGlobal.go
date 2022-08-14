package database

import (
	"database/sql"
	"github.com/MysGal/Mon3tr/utils"
	_ "github.com/mattn/go-sqlite3"
)

var GlobalDatabase *sql.DB

func InitDatabase() {
	database, err := sql.Open("sqlite3", "./data/database/database.db?cache=shared&mode=memory")

	if err != nil {
		utils.GlobalLogger.Fatal(err)
	}
	database.SetMaxOpenConns(1)
	GlobalDatabase = database
}
