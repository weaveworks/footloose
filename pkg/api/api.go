package api

import (
	"github.com/gorilla/mux"
	"github.com/weaveworks/footloose/pkg/cluster"
)

// API represents the footloose REST API.
type API struct {
	BaseURI  string
	db       db
	keyStore *cluster.KeyStore
}

// New creates a new object able to answer footloose REST API.
func New(baseURI string, keyStore *cluster.KeyStore) *API {
	api := &API{
		BaseURI:  baseURI,
		keyStore: keyStore,
	}
	api.db.init()
	return api
}

// Router returns the API request router.
func (a *API) Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/keys", a.createPublicKey).Methods("POST")
	router.HandleFunc("/api/keys/{key}", a.getPublicKey).Methods("GET")
	router.HandleFunc("/api/keys/{key}", a.deletePublicKey).Methods("DELETE")
	router.HandleFunc("/api/clusters", a.createCluster).Methods("POST")
	router.HandleFunc("/api/clusters/{cluster}", a.deleteCluster).Methods("DELETE")
	router.HandleFunc("/api/clusters/{cluster}/machines", a.createMachine).Methods("POST")
	router.HandleFunc("/api/clusters/{cluster}/machines/{machine}", a.getMachine).Methods("GET")
	router.HandleFunc("/api/clusters/{cluster}/machines/{machine}", a.deleteMachine).Methods("DELETE")
	return router
}
