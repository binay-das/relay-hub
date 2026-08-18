package database

import (
	"database/sql"

	"github.com/binay-das/relay-hub/internal/types"
)

func (s Store) AddHistory(userID int64, entry types.RequestHistory) error {
	headersJSON, err := marshalHeaders(entry.Headers)
	if err != nil {
		return err
	}

	responseHeadersJSON, err := marshalHeaders(entry.ResponseHeaders)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		`INSERT INTO request_history
		 (user_id, method, url, headers_json, request_body, status_code, status_text, elapsed_ms,
		  response_headers_json, response_body, body_type, error, message)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userID,
		entry.Method,
		entry.URL,
		nullJSON(headersJSON),
		entry.RequestBody,
		entry.StatusCode,
		entry.StatusText,
		entry.ElapsedMS,
		nullJSON(responseHeadersJSON),
		entry.ResponseBody,
		entry.BodyType,
		entry.Error,
		entry.Message,
	)
	return err
}

func (s Store) ListHistory(userID int64, limit int) ([]types.RequestHistory, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	rows, err := s.db.Query(
		`SELECT id, user_id, method, url, headers_json, request_body, status_code, status_text,
		        elapsed_ms, response_headers_json, response_body, body_type, error,
		        message, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ')
		 FROM request_history
		 WHERE user_id = ?
		 ORDER BY created_at DESC, id DESC
		 LIMIT ?`,
		userID,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := []types.RequestHistory{}
	for rows.Next() {
		var entry types.RequestHistory
		var headersJSON, responseHeadersJSON, requestBody, responseBody, message sql.NullString

		if err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.Method,
			&entry.URL,
			&headersJSON,
			&requestBody,
			&entry.StatusCode,
			&entry.StatusText,
			&entry.ElapsedMS,
			&responseHeadersJSON,
			&responseBody,
			&entry.BodyType,
			&entry.Error,
			&message,
			&entry.CreatedAt,
		); err != nil {
			return nil, err
		}

		entry.Headers = unmarshalHeaders(headersJSON)
		entry.RequestBody = requestBody.String
		entry.ResponseHeaders = unmarshalHeaders(responseHeadersJSON)
		entry.ResponseBody = responseBody.String
		entry.Message = message.String
		history = append(history, entry)
	}

	return history, rows.Err()
}

func (s Store) DeleteHistory(userID int64, id int64) error {
	result, err := s.db.Exec(`DELETE FROM request_history WHERE id = ? AND user_id = ?`, id, userID)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s Store) ClearHistory(userID int64) error {
	_, err := s.db.Exec(`DELETE FROM request_history WHERE user_id = ?`, userID)
	return err
}
