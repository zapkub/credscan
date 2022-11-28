package sqlite

import (
	"context"
	"testing"
)

func Test_GetQueuedJob(t *testing.T) {
	sqlitedb := newdb(t)
	id, err := sqlitedb.InsertJob(context.Background(),
		"TestRepo",
		"https://github.com/zapkub/react-thailand-address-typeahead",
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(id)

	queuedJob, err := sqlitedb.GetQueuedJob(context.Background())
	if err != nil {
		t.Error(err)
	}

	if queuedJob == nil {
		t.Error("expect result to be existed")
	}

}
