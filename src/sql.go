package main

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"google.golang.org/genai"
	_ "modernc.org/sqlite"
)

func initDb() {
	db, err := sql.Open("sqlite", DB_FILE)
	exitOnErr(err)
	defer db.Close()
	_, err = db.Exec(SQL_CREATE_TABLE)
	exitOnErr(err)
}

func insertMessage(body string, is_agent bool, is_topic bool) {
	db, err := sql.Open("sqlite", DB_FILE)
	exitOnErr(err)
	defer db.Close()

	ts_now := time.Now().Format(time.RFC3339Nano)
	topic_id := uuid.New()

	if is_topic {
		err = db.QueryRow(SQL_GET_CURRENT_ID).Scan(&topic_id)
		exitOnErr(err)
	}

	_, err = db.Exec(SQL_INSERT_ROW, topic_id, ts_now, is_agent, body)
}

func selectMessage(chat_history *[]*genai.Content) {
	db, err := sql.Open("sqlite", DB_FILE)
	exitOnErr(err)
	defer db.Close()

	var topic_id string
	err = db.QueryRow(SQL_GET_CURRENT_ID).Scan(&topic_id)
	exitOnErr(err)

	var rows *sql.Rows
	rows, err = db.Query(SQL_GET_CONVERSATION, topic_id)
	exitOnErr(err)
	defer rows.Close()

	var role genai.Role

	for rows.Next() {
		var body string
		var is_agent bool
		err := rows.Scan(&is_agent, &body)
		exitOnErr(err)
		role = genai.RoleUser
		if is_agent {
			role = genai.RoleModel
		}
		*chat_history = append(*chat_history, genai.NewContentFromText(body, role))
	}

}
