package database

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zapkub/credscan/internal/appfs"
)

type DB struct {
	db     *sql.DB
	appfs  *appfs.AppFileSystem
	dbname string
}

func (db *DB) Close() error {
	err := db.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteAll() error {
	err := db.Close()
	if err != nil {
		return err
	}
	err = db.appfs.Unlink(db.dbname)
	if err != nil {
		return err
	}

	return err
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return db.db.QueryRowContext(ctx, query, args...)
}

func (db *DB) RunQuery(ctx context.Context, query string, f func(*sql.Rows) error, params ...any) error {
	rows, err := db.db.QueryContext(ctx, query, params...)
	if err != nil {
		return err
	}
	_, err = processRows(rows, f)
	return err
}

func processRows(rows *sql.Rows, f func(*sql.Rows) error) (int, error) {
	defer rows.Close()
	n := 0
	for rows.Next() {
		n++
		if err := f(rows); err != nil {
			return n, err
		}
	}
	return n, rows.Err()
}

func (db *DB) Exec(ctx context.Context, query string, args ...any) error {
	_, err := db.db.ExecContext(ctx, query, args...)
	return err
}

func OpenSqlite(appfs *appfs.AppFileSystem) (_ *DB, err error) {

	var db DB
	db.appfs = appfs
	db.dbname = "/credscan.db"
	sqldb, err := sql.Open("sqlite3", appfs.Pathname(db.dbname))
	if err != nil {
		return nil, err
	}
	db.db = sqldb

	return &db, nil
}
