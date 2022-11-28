package main

import (
	"context"
	"os"

	"github.com/zapkub/credscan/internal/appfs"
	"github.com/zapkub/credscan/internal/appworker"
	"github.com/zapkub/credscan/internal/database"
	"github.com/zapkub/credscan/internal/sqlite"
)

func main() {
	var err error

	var afs *appfs.AppFileSystem

	if workdir := os.Getenv("CREDSCAN_WORKDIR"); workdir != "" {
		afs = appfs.NewAppFileSystem(0755, appfs.WorkingDir(workdir))
	} else {
		afs = appfs.NewAppFileSystem(0755)
	}

	db, err := database.OpenSqlite(afs)
	if err != nil {
		panic(err)
	}
	sqlite := sqlite.New(db)

	var aw = appworker.New(sqlite, afs)
	ctx := context.Background()
	aw.Start(ctx)

}
