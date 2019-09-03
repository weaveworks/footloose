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

// MachineURI returns the URI identifying a machine in the REST API.
func (a *API) MachineURI(c *cluster.Cluster, m *cluster.Machine) string {
	return fmt.Sprintf("%s/api/clusters/%s/machines/%s", a.BaseURI, c.Name(), m.Hostname())
}

// createMachine creates a machine.
func (a *API) createMachine(w http.ResponseWriter, r *http.Request) {
	var def config.Machine
	if err := json.NewDecoder(r.Body).Decode(&def); err != nil {
		sendError(w, http.StatusBadRequest, errors.Wrap(err, "could not decode body"))
		return
	}
	if def.Name == "" {
		sendError(w, http.StatusBadRequest, errors.New("no machine name provided"))
		return
	}

	vars := mux.Vars(r)
	c, err := a.db.cluster(vars["cluster"])
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	m := c.NewMachine(&def)

	if err := c.CreateMachine(m, 0); err != nil {
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	if err := a.db.addMachine(vars["cluster"], m); err != nil {
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	sendCreated(w, a.MachineURI(c, m))
}

// getMachine returns a machine object
func (a *API) getMachine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	m, err := a.db.machine(vars["cluster"], vars["machine"])
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	formatter, err := cluster.GetFormatter("json")
	if err != nil {
		sendError(w, http.StatusInternalServerError, err)
		return
	}
	if err := formatter.FormatSingle(w, m); err != nil {
		sendError(w, http.StatusInternalServerError, err)
		return
	}
}

// deleteMachine deletes a machine.
func (a *API) deleteMachine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c, err := a.db.cluster(vars["cluster"])
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}
	m, err := a.db.machine(vars["cluster"], vars["machine"])
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	if err := c.DeleteMachine(m, 0); err != nil {
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	_, err = a.db.removeMachine(vars["cluster"], vars["machine"])
	if err != nil {
		sendError(w, http.StatusInternalServerError, err)
		return
	}

	sendOK(w)
}
