package sqlite

import (
	"context"
	"log"
	"time"

	"github.com/zapkub/credscan/internal"
)

func (db *DB) UpdateJobStatusInProgress(ctx context.Context, jobID int) (err error) {
	log.Printf("update job id %v to status %v", jobID, internal.JobStatusInProgress)
	return db.db.Exec(ctx, `
		UPDATE Job	
		SET status = $1, scanningAt = $2
		WHERE id = $3
	`, internal.JobStatusInProgress, time.Now(), jobID)
}

func (db *DB) UpdateJobResult(ctx context.Context, jobID int, results []*internal.Finding) (err error) {
	return updateJobResult(ctx, db, jobID, results)
}

func updateJobResult(ctx context.Context, db *DB, jobID int, results []*internal.Finding) (err error) {
	var findings internal.Findings
	findings.Findings = results
	return db.db.Exec(ctx, `
		UPDATE Job	
		SET status = $1, findings = $2, finnishedAt = $3
		WHERE id = $4
	`, internal.JobStatusSuccess, findings, time.Now(), jobID)
}
