package database

import (
	"time"
)

func (s Store) CreateSession(userID int64, tokenHash string, expiresAt time.Time) error {
	_, err := s.db.Exec(
		`INSERT INTO auth_sessions (token_hash, user_id, expires_at) VALUES (?, ?, ?)`,
		tokenHash,
		userID,
		expiresAt.UTC(),
	)
	return err
}

func (s Store) UserIDBySession(tokenHash string) (int64, error) {
	var userID int64
	err := s.db.QueryRow(
		`SELECT user_id
		 FROM auth_sessions
		 WHERE token_hash = ? AND expires_at > UTC_TIMESTAMP()`,
		tokenHash,
	).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s Store) DeleteSession(tokenHash string) error {
	_, err := s.db.Exec(`DELETE FROM auth_sessions WHERE token_hash = ?`, tokenHash)
	return err
}
