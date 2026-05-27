package sql

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"google.golang.org/genai"
	_ "modernc.org/sqlite"

	"github.com/Starscade/Gaia/internal/text"
)

func GetLastBody(database_pointer *sql.DB) (string, error) {
	var body string
	err := database_pointer.QueryRow(text.SqlGetLastResponse).Scan(&body)
	return body, err
}

func Init(db_file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", db_file)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(text.SqlCreateTable)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InsertMessage(database_pointer *sql.DB, body string, is_agent, is_topic bool) error {
	ts_now := time.Now().Format(time.RFC3339Nano)
	topic_id := uuid.New()

	if is_topic {
		err := database_pointer.QueryRow(text.SqlGetCurrentId).Scan(&topic_id)
		if err != nil {
			return err
		}
	}

	_, err := database_pointer.Exec(text.SqlInsertRow, topic_id, ts_now, is_agent, body)
	if err != nil {
		return err
	}
	return nil
}

func SelectMessage(database_pointer *sql.DB, chat_history *[]*genai.Content) ([]*genai.Content, error) {
	var topic_id string
	err := database_pointer.QueryRow(text.SqlGetCurrentId).Scan(&topic_id)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = database_pointer.Query(text.SqlGetConversation, topic_id)
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

func TruncateTopic(database_pointer *sql.DB) error {
	_, err := database_pointer.Exec(text.SqlTruncateTopic)
	if err != nil {
		return err
	}
	return nil
}

func TruncateTranscript(database_pointer *sql.DB) error {
	_, err := database_pointer.Exec(text.SqlTruncateTranscript)
	if err != nil {
		return err
	}
	return nil
}
