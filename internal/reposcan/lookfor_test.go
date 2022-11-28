package reposcan

import (
	"context"
	"strings"
	"testing"
	"time"
)

func Test_Lookfor_ShouldFoundKeyword(t *testing.T) {

	var inputbuf = strings.NewReader("public_key")
	var ctx = context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	issue, err := Lookfor([]LookForRule{{
		RuleID:   "gs200",
		Keywords: []string{"public_key"},
	}})(
		ctx,
		"/some/where",
		inputbuf,
	)

	if err != nil {
		t.Error(err)
	}

	if issue[0].RuleID != "gs200" {
		t.Errorf("Expect for lookfor issue")
	}

}

func Test_Lookfor_ShouldNotFoundKeyword(t *testing.T) {

	var inputbuf = strings.NewReader("public_key")
	var ctx = context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	issues, err := Lookfor([]LookForRule{{
		RuleID:   "gs200",
		Keywords: []string{"private_key"},
	}})(
		ctx,
		"/path/file",
		inputbuf,
	)
	if err != nil {
		t.Error(err)
	}

	if len(issues) != 0 {
		t.Error("expect no finding")
	}

}
