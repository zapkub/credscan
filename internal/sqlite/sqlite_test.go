package sqlite

import (
	"os"
	"path"
	"testing"

	"github.com/zapkub/credscan/internal/appfs"
	"github.com/zapkub/credscan/internal/database"
)

func newdb(t *testing.T) *DB {
	wd, _ := os.Getwd()
	var appfs = appfs.NewAppFileSystem(
		0755,
		appfs.WorkingDir(path.Join(wd, "../../test/work_dir/.test_data", t.Name())),
	)
	db, err := database.OpenSqlite(appfs)
	if err != nil {
		panic(err)
	}
	err = db.DeleteAll()
	if err != nil {
		panic(err)
	}

	db, err = database.OpenSqlite(appfs)
	if err != nil {
		panic(err)
	}
	sqlitedb := New(db)
	err = sqlitedb.MigrateJobTable()
	if err != nil {
		panic(err)
	}
	return sqlitedb
}
