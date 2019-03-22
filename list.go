package main

import (
	"github.com/spf13/cobra"
	"github.com/weaveworks/footloose/pkg/cluster"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all running machines",
	RunE:  list,
}

var listOptions struct {
	format string
	config string
	all    bool
}

func init() {
	listCmd.Flags().StringVarP(&listOptions.config, "config", "c", Footloose, "Cluster configuration file")
	listCmd.Flags().StringVarP(&listOptions.format, "format", "f", "default", "Formatting options")
	listCmd.Flags().BoolVar(&listOptions.all, "all", false, "List all footloose created machines in every cluster.")
	footloose.AddCommand(listCmd)
}

// list will list all machines in a given cluster.
// if --all option is provided it will list every machine created by
// footloose no matter what cluster they are in.
func list(cmd *cobra.Command, args []string) error {
	cluster, err := cluster.NewFromFile(listOptions.config)
	if err != nil {
		return err
	}
	return cluster.List(listOptions.all, listOptions.format)
}
