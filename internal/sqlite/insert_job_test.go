package sqlite

import (
	"context"
	"testing"
)

func Test_InsertJob(t *testing.T) {
	sqlitedb := newdb(t)
	id, err := sqlitedb.InsertJob(context.Background(),
		"TestRepo",
		"https://github.com/zapkub/react-thailand-address-typeahead",
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
}
