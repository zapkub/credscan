package appcontext

import (
	"context"

	"github.com/zapkub/credscan/internal"
)

type DB interface {
	InsertJob(ctx context.Context, repositoryName string, repositoryUrl string) (jobID int, err error)
	QueryJobs(ctx context.Context) ([]*internal.Job, error)
	DeleteJob(ctx context.Context, jobID int) error
}

type ContextDBKey struct{}

func WithDB(ctx context.Context, db DB) context.Context {
	return context.WithValue(ctx, ContextDBKey{}, db)
}

func MustGetDB(ctx context.Context) DB {
	if db, ok := ctx.Value(ContextDBKey{}).(DB); ok {
		return db
	}
	panic("unable to get DB from context")
}
