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

func insertMessage(body string, is_agent bool, is_topic bool) {
	db, err := sql.Open("sqlite", DEFAULT_DB_FILENAME)
	exitOnErr(err)
	defer db.Close()

	ts_now := time.Now().Format(time.RFC3339Nano)
	conversation_id := uuid.New()

	if is_topic {
		err = db.QueryRow(SQL_GET_CURRENT_ID).Scan(&conversation_id)
		exitOnErr(err)
	}

	_, err = db.Exec(SQL_INSERT_ROW, conversation_id, ts_now, is_agent, body)
}

func selectMessage() string {
	db, err := sql.Open("sqlite", DEFAULT_DB_FILENAME)
	exitOnErr(err)
	defer db.Close()

	var conversation_id string
	err = db.QueryRow(SQL_GET_CURRENT_ID).Scan(&conversation_id)
	exitOnErr(err)

	var rows *sql.Rows
	rows, err = db.Query(SQL_GET_CONVERSATION, conversation_id)
	exitOnErr(err)
	defer rows.Close()

	chat_history := ""

	for rows.Next() {
		var body string
		err := rows.Scan(&body)
		exitOnErr(err)
		chat_history = chat_history + "\n\n" + body
	}

	return chat_history
}
