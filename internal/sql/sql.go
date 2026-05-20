package sql

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"google.golang.org/genai"
	_ "modernc.org/sqlite"

	"github.com/Starscade/Gaia/internal/text"
	"github.com/Starscade/Gaia/internal/tools"
)

func GetLastBody(db_file string) (string, error) {
	db, err := sql.Open("sqlite", db_file)
	tools.ExitOnErr(err)
	defer db.Close()

	var body string
	err = db.QueryRow(text.SQL_GET_LAST_RESPONSE).Scan(&body)
	return body, err
}

func Init(db_file string) {
	db, err := sql.Open("sqlite", db_file)
	tools.ExitOnErr(err)
	defer db.Close()
	_, err = db.Exec(text.SQL_CREATE_TABLE)
	tools.ExitOnErr(err)
}

func InsertMessage(db_file, body string, is_agent, is_topic bool) {
	db, err := sql.Open("sqlite", db_file)
	tools.ExitOnErr(err)
	defer db.Close()

	ts_now := time.Now().Format(time.RFC3339Nano)
	topic_id := uuid.New()

	if is_topic {
		err = db.QueryRow(text.SQL_GET_CURRENT_ID).Scan(&topic_id)
	}

	_, err = db.Exec(text.SQL_INSERT_ROW, topic_id, ts_now, is_agent, body)
	tools.ExitOnErr(err)
}

func SelectMessage(db_file string, chat_history *[]*genai.Content) {
	db, err := sql.Open("sqlite", db_file)
	tools.ExitOnErr(err)
	defer db.Close()

	var topic_id string
	err = db.QueryRow(text.SQL_GET_CURRENT_ID).Scan(&topic_id)
	if err != nil {
		return
	}

	var rows *sql.Rows
	rows, err = db.Query(text.SQL_GET_CONVERSATION, topic_id)
	if err != nil {
		return
	}
	defer rows.Close()

	var role genai.Role

	for rows.Next() {
		var body string
		var is_agent bool
		err := rows.Scan(&is_agent, &body)
		if err != nil {
			return
		}
		role = genai.RoleUser
		if is_agent {
			role = genai.RoleModel
		}
		*chat_history = append(*chat_history, genai.NewContentFromText(body, role))
	}

}

func TruncateTranscript(db_file string) {
	db, err := sql.Open("sqlite", db_file)
	tools.ExitOnErr(err)
	defer db.Close()
	_, err2 := db.Exec(text.SQL_TRUNCATE_TRANSCRIPT)
	tools.ExitOnErr(err2)
}
