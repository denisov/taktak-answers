package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // sqlite3
	"log"
)

var db *sql.DB

func InitDb(dbFile string) {
	log.Println("Opening DB on " + dbFile)
	var err error

	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalln(err)
	}

	err = solutionsInit()
	if err != nil {
		log.Fatalln(err)
	}
}
