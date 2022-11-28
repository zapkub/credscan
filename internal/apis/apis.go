package apis

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zapkub/credscan/internal/appcontext"
)

type AppAPI struct {
	m  *mux.Router
	db appcontext.DB
}

func (r *AppAPI) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.m.ServeHTTP(w, req)
}

func NewAppAPI(db appcontext.DB) *AppAPI {

	var appapi AppAPI
	appapi.m = mux.NewRouter()
	appapi.db = db

	appapi.m.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			r = r.WithContext(appcontext.WithDB(ctx, db))
			h.ServeHTTP(w, r)
		})
	})

	appapi.m.Methods("GET").Path("/jobs").HandlerFunc(GetJobs)
	appapi.m.Methods("POST").Path("/job").HandlerFunc(PostJob)
	appapi.m.Methods("DELETE").Path("/job").HandlerFunc(DeleteJob)

	return &appapi
}

type ErrorBody struct {
	Error string `json:"error"`
}
