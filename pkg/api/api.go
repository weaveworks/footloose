package api

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/weaveworks/footloose/pkg/cluster"
)

// API represents the footloose REST API.
type API struct {
	BaseURI  string
	db       db
	keyStore *cluster.KeyStore
	router   *mux.Router
}

// New creates a new object able to answer footloose REST API.
func New(baseURI string, keyStore *cluster.KeyStore, debug bool) *API {
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	api := &API{
		BaseURI:  baseURI,
		keyStore: keyStore,
		router:   mux.NewRouter(),
	}
	api.db.init()
	return api
}

func httpLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugln(r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}

func (a *API) createDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string][]string{}

	_ = a.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()

		if _, ok := response[path]; ok {
			response[path] = append(response[path], methods[0])
		} else {
			response[path] = methods
		}

		return nil
	})

	payload, err := json.Marshal(response)
	if err != nil {
		log.Fatalln(err)
	}

	sendResponse(w, payload)
}

func (a *API) initRouter() {
	a.router.HandleFunc("/api/keys", a.createPublicKey).Methods("POST")
	a.router.HandleFunc("/api/keys/{key}", a.getPublicKey).Methods("GET")
	a.router.HandleFunc("/api/keys/{key}", a.deletePublicKey).Methods("DELETE")
	a.router.HandleFunc("/api/clusters", a.createCluster).Methods("POST")
	a.router.HandleFunc("/api/clusters/{cluster}", a.deleteCluster).Methods("DELETE")
	a.router.HandleFunc("/api/clusters/{cluster}/machines", a.createMachine).Methods("POST")
	a.router.HandleFunc("/api/clusters/{cluster}/machines/{machine}", a.getMachine).Methods("GET")
	a.router.HandleFunc("/api/clusters/{cluster}/machines/{machine}", a.deleteMachine).Methods("DELETE")

	a.router.HandleFunc("/", a.createDocs).Methods("GET")

	a.router.Use(httpLogger)
}

// Router returns the API request router.
func (a *API) Router() *mux.Router {
	a.initRouter()
	return a.router
}
