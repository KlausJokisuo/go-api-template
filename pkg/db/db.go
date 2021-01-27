package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	Client *pgxpool.Pool
}

func Get(connStr string) (*DB, error) {
	db, err := get(connStr)
	if err != nil {
		return nil, err
	}

	return &DB{
		Client: db,
	}, nil
}

func (d *DB) Close() {
	d.Client.Close()
}

func get(connStr string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	//
	//if err := db.Ping(); err != nil {
	//	return nil, err
	//}

	return db, nil
}
