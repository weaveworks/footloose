package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/weaveworks/footloose/pkg/cluster"
	"github.com/weaveworks/footloose/pkg/config"
)

// ClusterURI returns the URI identifying a cluster in the REST API.
func (a *API) ClusterURI(c *cluster.Cluster) string {
	return fmt.Sprintf("%s/api/clusters/%s", a.BaseURI, c.Name())
}

// createCluster creates a cluster.
func (a *API) createCluster(w http.ResponseWriter, r *http.Request) {
	var def config.Cluster
	if err := json.NewDecoder(r.Body).Decode(&def); err != nil {
		sendError(w, http.StatusBadRequest, errors.Wrap(err, "could not decode body"))
		return
	}
	if def.Name == "" {
		sendError(w, http.StatusBadRequest, errors.New("no cluster name provided"))
		return
	}

	cluster, err := cluster.New(config.Config{Cluster: def})
	if err != nil {
		sendError(w, http.StatusInternalServerError, err)
		return
	}
	cluster.SetKeyStore(a.keyStore)

	if err := a.db.addCluster(def.Name, cluster); err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	if err := cluster.Create(); err != nil {
		_, _ = a.db.removeCluster(def.Name)
		sendError(w, http.StatusInternalServerError, err)
		return
	}
	sendCreated(w, a.ClusterURI((cluster)))
}

// deleteCluster deletes a cluster.
func (a *API) deleteCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c, err := a.db.cluster(vars["cluster"])
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	// Starts by deleting the machines associated with the cluster.
	machines, err := a.db.machines(vars["cluster"])
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}
	for _, m := range machines {
		if err := c.DeleteMachine(m, 0); err != nil {
			sendError(w, http.StatusInternalServerError, err)
			return
		}
		if _, err := a.db.removeMachine(vars["cluster"], m.Hostname()); err != nil {
			sendError(w, http.StatusInternalServerError, err)
			return
		}
	}

	// Delete cluster.
	if err := c.Delete(); err != nil {
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	_, err = a.db.removeCluster(vars["cluster"])
	if err != nil {
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	sendOK(w)
}
