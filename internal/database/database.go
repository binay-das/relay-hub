package database

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectToDB(dbDsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dbDsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(120) NOT NULL UNIQUE,
			email VARCHAR(190) NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS auth_sessions (
			token_hash CHAR(64) PRIMARY KEY,
			user_id BIGINT NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_auth_sessions_user_id (user_id),
			INDEX idx_auth_sessions_expires_at (expires_at)
		)`,
		`CREATE TABLE IF NOT EXISTS collections (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			user_id BIGINT NOT NULL,
			name VARCHAR(120) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_collections_user_id (user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS saved_requests (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			user_id BIGINT NOT NULL,
			collection_id BIGINT NULL,
			name VARCHAR(160) NOT NULL,
			method VARCHAR(12) NOT NULL,
			url TEXT NOT NULL,
			headers_json JSON NULL,
			body MEDIUMTEXT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_saved_requests_user_id (user_id),
			CONSTRAINT fk_saved_requests_collection
				FOREIGN KEY (collection_id) REFERENCES collections(id)
				ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS request_history (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			user_id BIGINT NOT NULL,
			method VARCHAR(12) NOT NULL,
			url TEXT NOT NULL,
			headers_json JSON NULL,
			request_body MEDIUMTEXT NULL,
			status_code INT NOT NULL DEFAULT 0,
			status_text VARCHAR(80) NOT NULL DEFAULT '',
			elapsed_ms DOUBLE NOT NULL DEFAULT 0,
			response_headers_json JSON NULL,
			response_body MEDIUMTEXT NULL,
			body_type VARCHAR(20) NOT NULL DEFAULT 'text',
			error BOOLEAN NOT NULL DEFAULT FALSE,
			message TEXT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_request_history_user_id (user_id),
			INDEX idx_request_history_created_at (created_at)
		)`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}

	if err := ensureColumn(db, "users", "email", `ALTER TABLE users ADD COLUMN email VARCHAR(190) NULL UNIQUE AFTER name`); err != nil {
		return err
	}
	if err := ensureColumn(db, "users", "password_hash", `ALTER TABLE users ADD COLUMN password_hash VARCHAR(255) NOT NULL DEFAULT '' AFTER email`); err != nil {
		return err
	}

	legacyUserID, err := ensureUser(db, "legacy")
	if err != nil {
		return err
	}

	if err := ensureUserIDColumn(db, "collections", legacyUserID); err != nil {
		return err
	}
	if err := ensureUserIDColumn(db, "saved_requests", legacyUserID); err != nil {
		return err
	}
	if err := ensureUserIDColumn(db, "request_history", legacyUserID); err != nil {
		return err
	}

	return nil
}

func ensureUser(db *sql.DB, name string) (int64, error) {
	email := name + "@relay.local"
	_, err := db.Exec(
		`INSERT IGNORE INTO users (name, email, password_hash) VALUES (?, ?, '')`,
		name,
		email,
	)
	if err != nil {
		return 0, err
	}

	var id int64
	err = db.QueryRow(`SELECT id FROM users WHERE name = ?`, name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func ensureColumn(db *sql.DB, tableName string, columnName string, alterStatement string) error {
	var found string
	err := db.QueryRow(
		`SELECT column_name FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = ? AND column_name = ?`,
		tableName,
		columnName,
	).Scan(&found)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = db.Exec(alterStatement)
	return err
}

func ensureUserIDColumn(db *sql.DB, tableName string, legacyUserID int64) error {
	if err := ensureColumn(db, tableName, "user_id", `ALTER TABLE `+tableName+` ADD COLUMN user_id BIGINT NULL AFTER id`); err != nil {
		return err
	}

	if _, err := db.Exec(`UPDATE `+tableName+` SET user_id = ? WHERE user_id IS NULL`, legacyUserID); err != nil {
		return err
	}

	return nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
