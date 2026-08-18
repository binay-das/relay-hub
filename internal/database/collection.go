package database

import (
	"github.com/binay-das/relay-hub/internal/types"
)

func (s Store) CreateCollection(userID int64, name string) (types.Collection, error) {
	result, err := s.db.Exec(`INSERT INTO collections (user_id, name) VALUES (?, ?)`, userID, name)
	if err != nil {
		return types.Collection{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return types.Collection{}, err
	}

	return s.GetCollection(userID, id)
}

func (s Store) GetCollection(userID int64, id int64) (types.Collection, error) {
	var collection types.Collection

	err := s.db.QueryRow(
		`SELECT id, user_id, name,
		        DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ'),
		        DATE_FORMAT(updated_at, '%Y-%m-%dT%H:%i:%sZ')
		 FROM collections
		 WHERE id = ? AND user_id = ?`,
		id,
		userID,
	).Scan(&collection.ID, &collection.UserID, &collection.Name, &collection.CreatedAt, &collection.UpdatedAt)
	if err != nil {
		return types.Collection{}, err
	}

	return collection, nil
}

func (s Store) ListCollections(userID int64) ([]types.Collection, error) {
	rows, err := s.db.Query(
		`SELECT id, user_id, name,
		        DATE_FORMAT(created_at, '%Y-%m-%dT%H:%i:%sZ'),
		        DATE_FORMAT(updated_at, '%Y-%m-%dT%H:%i:%sZ')
		 FROM collections
		 WHERE user_id = ?
		 ORDER BY name ASC, id DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collections := []types.Collection{}
	for rows.Next() {
		var collection types.Collection
		if err := rows.Scan(
			&collection.ID,
			&collection.UserID,
			&collection.Name,
			&collection.CreatedAt,
			&collection.UpdatedAt,
		); err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	return collections, rows.Err()
}
