package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
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

	ts_now := time.Now().Format(time.RFC3339)
	new_uuid := uuid.New()

	_, err = db.Exec(SQL_INSERT_ROW, new_uuid, ts_now, is_agent, body)
}

func selectMessage(conversation_id string) string {
	db, err := sql.Open("sqlite", DEFAULT_DB_FILENAME)
	exitOnErr(err)
	defer db.Close()

	var rows *sql.Rows
	rows, err = db.Query(SQL_GET_CONVERSATIONS)
	exitOnErr(err)
	defer rows.Close()

	chat_history := ""

	for rows.Next() {
		var body string
		err := rows.Scan(&body)
		exitOnErr(err)
		chat_history = chat_history + body
	}

	return chat_history
}
