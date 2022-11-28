package apis

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/zapkub/credscan/internal/appcontext"
)

type PostJobResponseData struct {
	ID int `json:"id"`
}
type PostJobResponse struct {
	Data PostJobResponseData `json:"data"`
}

func PostJob(w http.ResponseWriter, req *http.Request) {

	var repositoryURL = req.FormValue("repositoryUrl")
	var repositoryName = req.FormValue("repositoryName")
	var enc = json.NewEncoder(w)

	if repositoryURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(ErrorBody{
			Error: "repositoryUrl is required",
		})
		return
	}

	if _, err := url.Parse(repositoryURL); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(ErrorBody{
			Error: "repositoryUrl is" + err.Error(),
		})
		return
	}

	if repositoryName == "" {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(ErrorBody{
			Error: "repositoryName is required",
		})
		return
	}

	id, err := createNewJob(req.Context(), repositoryName, repositoryURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(ErrorBody{
			Error: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	enc.Encode(PostJobResponse{
		Data: PostJobResponseData{
			ID: id,
		},
	})
}

func createNewJob(ctx context.Context, repositoryName, repositoryURL string) (id int, err error) {
	db := appcontext.MustGetDB(ctx)
	id, err = db.InsertJob(ctx, repositoryName, repositoryURL)
	return id, err
}
