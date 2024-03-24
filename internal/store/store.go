package store

import (
	"database/sql"
	db "github.com/tredoc/go-balances/db/sqlc"
)

type Store struct {
	DB *sql.DB
	*db.Queries
}

func New(conn *sql.DB) *Store {
	return &Store{
		DB:      conn,
		Queries: db.New(conn),
	}
}
