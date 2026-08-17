package database

import (
	"database/sql"

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
