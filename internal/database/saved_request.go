package database

import (
	"database/sql"

	"github.com/binay-das/relay-hub/internal/types"
)

func (s Store) SaveRequest(userID int64, payload types.SaveRequestPayload) (types.SavedRequest, error) {
	headersJSON, err := marshalHeaders(payload.Headers)
	if err != nil {
		return types.SavedRequest{}, err
	}

	if payload.CollectionID != nil {
		if _, err := s.GetCollection(userID, *payload.CollectionID); err != nil {
			return types.SavedRequest{}, err
		}
	}

	result, err := s.db.Exec(
		`INSERT INTO saved_requests (user_id, collection_id, name, method, url, headers_json, body)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userID,
		payload.CollectionID,
		payload.Name,
		payload.Method,
		payload.URL,
		nullJSON(headersJSON),
		payload.Body,
	)
	if err != nil {
		return types.SavedRequest{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return types.SavedRequest{}, err
	}

	return s.GetSavedRequest(userID, id)
}

func (s Store) ListSavedRequests(userID int64) ([]types.SavedRequest, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, collection_id, name, method, url, headers_json, body,
		        DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ'),
		        DATE_FORMAT(updated_at, '%Y-%m-%dT%H:%i:%sZ')
		 FROM saved_requests
		 WHERE user_id = ?
		 ORDER BY updated_at DESC, id DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []types.SavedRequest{}
	for rows.Next() {
		request, err := scanSavedRequest(rows)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, rows.Err()
}

func (s Store) GetSavedRequest(userID int64, id int64) (types.SavedRequest, error) {
	row := s.db.QueryRow(
		`SELECT id, user_id, collection_id, name, method, url, headers_json, body,
		        DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ'),
		        DATE_FORMAT(updated_at, '%Y-%m-%dT%H:%i:%sZ')
		 FROM saved_requests
		 WHERE id = ? AND user_id = ?`,
		id,
		userID,
	)

	return scanSavedRequest(row)
}

func (s Store) DeleteSavedRequest(userID int64, id int64) error {
	result, err := s.db.Exec(`DELETE FROM saved_requests WHERE id = ? AND user_id = ?`, id, userID)
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

func scanSavedRequest(row rowScanner) (types.SavedRequest, error) {
	var request types.SavedRequest
	var collectionID sql.NullInt64
	var headersJSON, body sql.NullString

	err := row.Scan(
		&request.ID,
		&request.UserID,
		&collectionID,
		&request.Name,
		&request.Method,
		&request.URL,
		&headersJSON,
		&body,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err != nil {
		return types.SavedRequest{}, err
	}

	if collectionID.Valid {
		request.CollectionID = &collectionID.Int64
	}

	request.Headers = unmarshalHeaders(headersJSON)
	request.Body = body.String
	return request, nil
}
