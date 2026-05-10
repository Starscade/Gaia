package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func initDb() {
	db, err := sql.Open("sqlite", DEFAULT_DB_FILENAME)
	exitOnErr(err)
	defer db.Close()
	_, err = db.Exec(SQL_CREATE_TABLE)
	exitOnErr(err)
}

func insertMessage(body string, is_agent bool) {
	db, err := sql.Open("sqlite", DEFAULT_DB_FILENAME)
	exitOnErr(err)
	defer db.Close()
	today := "1970-01-01 00:00:00"
	_, err = db.Exec(SQL_INSERT_ROW, today, is_agent, body)
}
