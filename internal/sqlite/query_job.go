package sqlite

import (
	"context"
	"database/sql"

	"github.com/zapkub/credscan/internal"
)

func (db *DB) QueryJobs(ctx context.Context) ([]*internal.Job, error) {
	return queryJobs(ctx, db)
}

func queryJobs(ctx context.Context, db *DB) ([]*internal.Job, error) {
	var jobs = make([]*internal.Job, 0)
	err := db.db.RunQuery(
		ctx, `
		SELECT id, repositoryName, repositoryUrl, findings, status, queuedAt, scanningAt, finnishedAt FROM Job
	`, func(r *sql.Rows) error {
			var job internal.Job
			var err = r.Scan(&job.ID, &job.RepositoryName, &job.RepositoryURL, &job.Findings, &job.Status, &job.QueuedAt, &job.ScanningAt, &job.FinnishedAt)
			if err != nil {
				return err
			}
			jobs = append(jobs, &job)
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return jobs, err
}

func (db *DB) GetQueuedJob(ctx context.Context) (*internal.Job, error) {
	var job *internal.Job
	err := db.db.RunQuery(ctx, `
		SELECT id, repositoryUrl, status FROM Job WHERE status = $1 LIMIT 1
	`, func(r *sql.Rows) error {
		job = &internal.Job{}
		r.Scan(&job.ID, &job.RepositoryURL, &job.Status)
		return nil
	}, internal.JobStatusQueued)
	if err != nil {
		return nil, err
	}
	return job, nil
}
