package sqlite

import "github.com/zapkub/credscan/internal/database"

type DB struct {
	db *database.DB
}

func New(db *database.DB) *DB {
	return &DB{
		db: db,
	}
}
