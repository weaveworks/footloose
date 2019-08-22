package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/weaveworks/footloose/pkg/api"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch a footloose server",
	RunE:  serve,
}

var serveOptions struct {
	listen string
}

func baseURI(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	if host == "" || host == "0.0.0.0" || host == "[::]" {
		host = "localhost"
	}
	return fmt.Sprintf("http://%s:%s", host, port), nil
}

func init() {
	serveCmd.Flags().StringVarP(&serveOptions.listen, "listen", "l", ":2444", "Cluster configuration file")
	footloose.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) error {
	opts := &serveOptions

	baseURI, err := baseURI(opts.listen)
	if err != nil {
		return errors.Wrapf(err, "invalid listen address '%s'", opts.listen)
	}

	api := api.New(baseURI)

	router := mux.NewRouter()
	router.HandleFunc("/api/clusters", api.CreateCluster).Methods("POST")
	router.HandleFunc("/api/clusters/{cluster}", api.DeleteCluster).Methods("DELETE")
	router.HandleFunc("/api/clusters/{cluster}/machines", api.CreateMachine).Methods("POST")
	router.HandleFunc("/api/clusters/{cluster}/machines/{machine}", api.DeleteMachine).Methods("DELETE")

	return http.ListenAndServe(opts.listen, router)
}
