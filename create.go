package main

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/footloose/pkg/cluster"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a cluster",
	RunE:  create,
}

var createOptions struct {
	config string
}

func init() {
	createCmd.Flags().BoolVarP(&verbosity, "verbose", "v", false, "Verbosity commandline calls")
	createCmd.Flags().StringVarP(&createOptions.config, "config", "c", Footloose, "Cluster configuration file")
	footloose.AddCommand(createCmd)
}

func create(cmd *cobra.Command, args []string) error {
	cluster.GetCommanderInstance().SetVerbosity(verbosity)
	c, err := cluster.NewFromFile(createOptions.config)
	if err != nil {
		return err
	}
	return c.Create()
}
