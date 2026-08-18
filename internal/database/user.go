package database

import (
	"github.com/binay-das/relay-hub/internal/types"
)

func (s Store) CreateUser(email string, passwordHash string) (types.User, error) {
	result, err := s.db.Exec(
		`INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)`,
		email,
		email,
		passwordHash,
	)
	if err != nil {
		return types.User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return types.User{}, err
	}

	return s.GetUserByID(id)
}

func (s Store) GetUserByID(id int64) (types.User, error) {
	var user types.User

	err := s.db.QueryRow(
		`SELECT id, COALESCE(email, name), DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ')
		 FROM users
		 WHERE id = ?`,
		id,
	).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}

func (s Store) GetUserPasswordHash(email string) (types.User, string, error) {
	var user types.User
	var passwordHash string

	err := s.db.QueryRow(
		`SELECT id, COALESCE(email, name), password_hash, DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ')
		 FROM users
		 WHERE email = ?`,
		email,
	).Scan(&user.ID, &user.Email, &passwordHash, &user.CreatedAt)
	if err != nil {
		return types.User{}, "", err
	}

	return user, passwordHash, nil
}
