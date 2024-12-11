package db

import "github.com/jackc/pgx/v5/pgxpool"

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
	db DBTX
}

func NewStore(conn *pgxpool.Pool) Store {
	return &SQLStore{
		Queries: New(conn),
		db:      conn,
	}
}
