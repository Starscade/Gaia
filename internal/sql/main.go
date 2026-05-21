package sql

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"google.golang.org/genai"
	_ "modernc.org/sqlite"

	"github.com/Starscade/Gaia/internal/text"
)

func GetLastBody(db_file string) (string, error) {
	db, err := sql.Open("sqlite", db_file)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var body string
	err = db.QueryRow(text.SQL_GET_LAST_RESPONSE).Scan(&body)
	return body, err
}

func Init(db_file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", db_file)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	_, err = db.Exec(text.SQL_CREATE_TABLE)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InsertMessage(db_file, body string, is_agent, is_topic bool) error {
	db, err := sql.Open("sqlite", db_file)
	if err != nil {
		return err
	}
	defer db.Close()

	ts_now := time.Now().Format(time.RFC3339Nano)
	topic_id := uuid.New()

	if is_topic {
		err = db.QueryRow(text.SQL_GET_CURRENT_ID).Scan(&topic_id)
	}

	_, err = db.Exec(text.SQL_INSERT_ROW, topic_id, ts_now, is_agent, body)
	if err != nil {
		return err
	}
	return nil
}

func SelectMessage(db_file string, chat_history *[]*genai.Content) ([]*genai.Content, error) {
	db, err := sql.Open("sqlite", db_file)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var topic_id string
	err = db.QueryRow(text.SQL_GET_CURRENT_ID).Scan(&topic_id)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = db.Query(text.SQL_GET_CONVERSATION, topic_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var role genai.Role

	for rows.Next() {
		var body string
		var is_agent bool
		err := rows.Scan(&is_agent, &body)
		if err != nil {
			return nil, err
		}
		role = genai.RoleUser
		if is_agent {
			role = genai.RoleModel
		}
		*chat_history = append(*chat_history, genai.NewContentFromText(body, role))
	}

	return *chat_history, nil
}

func TruncateTranscript(db_file string) error {
	db, err := sql.Open("sqlite", db_file)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err2 := db.Exec(text.SQL_TRUNCATE_TRANSCRIPT)
	if err2 != nil {
		return err2
	}
	return nil
}
