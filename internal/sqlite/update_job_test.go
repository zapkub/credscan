package sqlite

import (
	"context"
	"testing"

	"github.com/zapkub/credscan/internal"
)

func Test_UpdateJobResult(t *testing.T) {
	sqlitedb := newdb(t)

	id, err := sqlitedb.InsertJob(context.Background(),
		"TestRepo",
		"https://github.com/zapkub/react-thailand-address-typeahead",
	)
	if err != nil {
		t.Error(err)
	}

	err = sqlitedb.UpdateJobResult(
		context.Background(),
		id,
		[]*internal.Finding{},
	)
	if err != nil {
		t.Error(err)
	}
}

func Test_UpdateJobStatusInProgress(t *testing.T) {
	sqlitedb := newdb(t)

	id, err := sqlitedb.InsertJob(context.Background(),
		"TestRepo",
		"https://github.com/zapkub/react-thailand-address-typeahead",
	)
	if err != nil {
		t.Error(err)
	}

	err = sqlitedb.UpdateJobStatusInProgress(
		context.Background(),
		id,
	)
	if err != nil {
		t.Error(err)
	}

	job, err := sqlitedb.GetQueuedJob(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if job != nil {
		t.Errorf("should not have queued job but %v status is %v", job.ID, job.Status)
	}

}
