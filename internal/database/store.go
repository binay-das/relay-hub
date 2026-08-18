package database

import (
	"database/sql"
	"encoding/json"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{db: db}
}

type rowScanner interface {
	Scan(dest ...any) error
}

func marshalHeaders(headers map[string]string) (string, error) {
	if headers == nil {
		headers = map[string]string{}
	}

	raw, err := json.Marshal(headers)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

func unmarshalHeaders(raw sql.NullString) map[string]string {
	headers := map[string]string{}
	if !raw.Valid || raw.String == "" {
		return headers
	}

	if err := json.Unmarshal([]byte(raw.String), &headers); err != nil {
		return map[string]string{}
	}

	return headers
}

func nullJSON(raw string) sql.NullString {
	if raw == "" {
		return sql.NullString{}
	}

	return sql.NullString{String: raw, Valid: true}
}
