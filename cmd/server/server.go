package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/zapkub/credscan/internal/apis"
	"github.com/zapkub/credscan/internal/appfs"
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
	err = db.DeleteAll()
	if err != nil {
		panic(err)
	}

	db, err = database.OpenSqlite(afs)
	if err != nil {
		panic(err)
	}
	sqlite := sqlite.New(db)
	err = sqlite.MigrateJobTable()
	if err != nil {
		panic(err)
	}
	var appAPI = apis.NewAppAPI(sqlite)
	log.Println("application start at :3000")
	go http.ListenAndServe("0.0.0.0:3000", appAPI)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigchan

	// TODO gracefully shutdown
	log.Println("exiting....")
	os.Exit(0)

}
