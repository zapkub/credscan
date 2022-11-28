package appworker

import (
	"context"
	"errors"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/zapkub/credscan/internal"
	"github.com/zapkub/credscan/internal/appfs"
	"github.com/zapkub/credscan/internal/reposcan"
	"github.com/zapkub/credscan/internal/sqlite"
)

var (
	ErrNoJob      = errors.New("no job")
	ErrNoMoreFile = errors.New("no more file")
)

type WorkerPool struct {
	lock       sync.Mutex
	ActiveJobs map[int]*internal.Job
	db         *sqlite.DB
	runner     *reposcan.Runner
}

func New(db *sqlite.DB, afs *appfs.AppFileSystem) *WorkerPool {
	return &WorkerPool{
		db:         db,
		ActiveJobs: make(map[int]*internal.Job),
		runner:     reposcan.New(afs),
	}
}

func (w *WorkerPool) Start(ctx context.Context) {
	log.Println("start worker")
	for {
		if len(w.ActiveJobs) >= 5 {
			log.Println("job limited")
		} else {
			go w.proceedJob(ctx)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):
			continue
		}
	}
}

func (w *WorkerPool) Terminate() {

}

func (w *WorkerPool) proceedJob(ctx context.Context) {
	queuedJob, err := w.db.GetQueuedJob(ctx)
	if err != nil {
		log.Println("unexpected error", err)
		return
	}

	if queuedJob == nil {
		log.Println("no job")
		return
	}

	log.Println("found queued job")
	log.Printf("job id = %v", queuedJob.ID)

	w.lock.Lock()
	err = w.db.UpdateJobStatusInProgress(ctx, queuedJob.ID)
	if err != nil {
		// TODO write error to job
		log.Println(err)
		w.lock.Unlock()
		return
	}

	w.ActiveJobs[queuedJob.ID] = queuedJob
	w.lock.Unlock()

	repourl, err := url.Parse(queuedJob.RepositoryURL)
	if err != nil {
		// TODO write error to job
		log.Println(err)
		return
	}

	// TODO better to stream write?
	findings, err := w.runner.Scan(ctx, repourl, DefaultRules)
	if err != nil {
		// TODO write error to job
		log.Println(err)
		return
	}

	w.lock.Lock()
	log.Printf("scan completed (%v)", findings)
	err = w.db.UpdateJobResult(ctx, queuedJob.ID, findings)
	if err != nil {
		// TODO write error to job
		log.Println(err)
		w.lock.Unlock()
		return
	}
	w.ActiveJobs[queuedJob.ID] = nil
	w.lock.Unlock()

}

var DefaultRules = []reposcan.RuleRunner{
	reposcan.Lookfor([]reposcan.LookForRule{
		{
			RuleID:   "credscan01",
			Keywords: []string{"public_key"},
		},
		{
			RuleID:   "credscan02",
			Keywords: []string{"private_key"},
		},
	}),
}
