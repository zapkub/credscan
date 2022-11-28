package sqlite

import (
	"context"
	"time"

	"github.com/zapkub/credscan/internal"
)

func (db *DB) InsertJob(ctx context.Context, repositoryName string, repositoryUrl string) (jobID int, err error) {
	now := time.Now()
	jobID, err = insertJob(ctx, db, &internal.Job{
		ID:            jobID,
		RepositoryURL: repositoryUrl,
		Findings: &internal.Findings{
			Findings: make([]*internal.Finding, 0),
		},
		RepositoryName: repositoryName,
		Status:         internal.JobStatusQueued,
		QueuedAt:       &now,
		ScanningAt:     nil,
		FinnishedAt:    nil,
	})
	if err != nil {
		return jobID, err
	}
	return jobID, err
}

func insertJob(ctx context.Context, db *DB, m *internal.Job) (_ int, err error) {
	var jobID int
	if err != nil {
		return 0, err
	}
	err = db.db.QueryRow(ctx, `
		INSERT INTO Job (
			repositoryUrl,
			repositoryName,
			findings,
			status,
			queuedAt
		) VALUES($1,$2,$3,$4,$5)
		RETURNING id
	`, m.RepositoryURL,
		m.RepositoryName,
		m.Findings,
		m.Status,
		m.QueuedAt,
	).Scan(&jobID)
	return jobID, err
}
