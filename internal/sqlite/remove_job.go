package sqlite

import "context"

func (db *DB) DeleteJob(ctx context.Context, jobID int) (err error) {
	return db.db.Exec(ctx, `
		DELETE FROM Job WHERE id = $1
	`, jobID)
}
