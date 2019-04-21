package main

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/footloose/pkg/cluster"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a cluster",
	RunE:  delete,
}

var deleteOptions struct {
	config string
}

func init() {
	deleteCmd.Flags().BoolVarP(&verbosity, "verbose", "v", false, "Verbosity commandline calls")
	deleteCmd.Flags().StringVarP(&deleteOptions.config, "config", "c", Footloose, "Cluster configuration file")
	footloose.AddCommand(deleteCmd)
}

func delete(cmd *cobra.Command, args []string) error {
	cluster.GetCommanderInstance().SetVerbosity(verbosity)
	c, err := cluster.NewFromFile(deleteOptions.config)
	if err != nil {
		return err
	}
	return c.Delete()
}
