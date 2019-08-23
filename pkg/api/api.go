package api

import "github.com/gorilla/mux"

// API represents the footloose REST API.
type API struct {
	BaseURI string
	db      db
}

// New creates a new object able to answer footloose REST API.
func New(baseURI string) *API {
	api := &API{
		BaseURI: baseURI,
	}
	api.db.init()
	return api
}

// Router returns the API request router.
func (a *API) Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/clusters", a.createCluster).Methods("POST")
	router.HandleFunc("/api/clusters/{cluster}", a.deleteCluster).Methods("DELETE")
	router.HandleFunc("/api/clusters/{cluster}/machines", a.createMachine).Methods("POST")
	router.HandleFunc("/api/clusters/{cluster}/machines/{machine}", a.deleteMachine).Methods("DELETE")
	return router
}
