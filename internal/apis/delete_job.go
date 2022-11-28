package apis

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zapkub/credscan/internal/appcontext"
)

type DeleteJobResponse struct {
}

func DeleteJob(w http.ResponseWriter, req *http.Request) {
	var err error
	jobIDStr := req.URL.Query().Get("id")
	var enc = json.NewEncoder(w)
	if jobIDStr == "" {
		w.WriteHeader(400)
		enc.Encode(ErrorBody{
			Error: "invalid request, id is required",
		})
		return
	}

	var jobID int
	if jobID, err = strconv.Atoi(jobIDStr); err != nil {
		w.WriteHeader(400)
		enc.Encode(ErrorBody{
			Error: "unable to parse job id",
		})
		return
	}

	err = deleteJob(req.Context(), jobID)
	if err != nil {
		w.WriteHeader(500)
		enc.Encode(ErrorBody{
			Error: err.Error(),
		})
		return
	}
	w.WriteHeader(200)
	enc.Encode(DeleteJobResponse{})
}

func deleteJob(ctx context.Context, jobID int) error {
	db := appcontext.MustGetDB(ctx)
	return db.DeleteJob(ctx, jobID)
}
