package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/weaveworks/footloose/pkg/config"
)

func (a *API) keyURI(name string) string {
	return fmt.Sprintf("%s/keys/%s", a.BaseURI, name)
}

func (a *API) createPublicKey(w http.ResponseWriter, r *http.Request) {
	var def config.PublicKey
	if err := json.NewDecoder(r.Body).Decode(&def); err != nil {
		sendError(w, http.StatusBadRequest, errors.Wrap(err, "could not decode body"))
		return
	}
	if def.Name == "" {
		sendError(w, http.StatusBadRequest, errors.New("no key name provided"))
		return
	}

	if err := a.keyStore.Store(def.Name, def.Key); err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	sendCreated(w, a.keyURI(def.Name))
}

func (a *API) getPublicKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data, err := a.keyStore.Get(vars["key"])
	if err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}

	key := config.PublicKey{
		Name: vars["key"],
		Key:  string(data),
	}
	if err := json.NewEncoder(w).Encode(&key); err != nil {
		sendError(w, http.StatusInternalServerError, err)
		return
	}
}

func (a *API) deletePublicKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if err := a.keyStore.Remove(vars["key"]); err != nil {
		sendError(w, http.StatusBadRequest, err)
		return
	}
	sendOK(w)
}
