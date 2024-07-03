package db

import (
	"database/sql"
)

// Store provides all functions to run individual queries as well as transactions
type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return Store{
		db: db,
	}
}

func (s Store) GetDB() *sql.DB {
	return s.db
}
