package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/weaveworks/footloose/pkg/api"
	"github.com/weaveworks/footloose/pkg/cluster"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch a footloose server",
	RunE:  serve,
}

var serveOptions struct {
	listen       string
	keyStorePath string
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
	serveCmd.Flags().StringVar(&serveOptions.keyStorePath, "keystore-path", defaultKeyStorePath, "Path of the public keys store")
	footloose.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) error {
	opts := &serveOptions

	baseURI, err := baseURI(opts.listen)
	if err != nil {
		return errors.Wrapf(err, "invalid listen address '%s'", opts.listen)
	}

	keyStore := cluster.NewKeyStore(opts.keyStorePath)
	if err := keyStore.Init(); err != nil {
		return errors.Wrapf(err, "could not init keystore")
	}

	api := api.New(baseURI, keyStore)
	router := api.Router()

	return http.ListenAndServe(opts.listen, router)
}
