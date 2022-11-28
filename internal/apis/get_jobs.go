package apis

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/zapkub/credscan/internal"
	"github.com/zapkub/credscan/internal/appcontext"
)

type GetJobsResponse struct {
	Data []*internal.Job `json:"data"`
}

func GetJobs(w http.ResponseWriter, req *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()
	var jobs []*internal.Job
	jobs, err = queryJobList(req.Context())
	var enc = json.NewEncoder(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(&ErrorBody{
			Error: err.Error(),
		})
		return
	}
	w.WriteHeader(200)
	err = enc.Encode(GetJobsResponse{
		Data: jobs,
	})
}

func queryJobList(ctx context.Context) ([]*internal.Job, error) {
	db := appcontext.MustGetDB(ctx)
	return db.QueryJobs(ctx)
}
