package main

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/footloose/pkg/cluster"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop cluster machines",
	RunE:  stop,
}

var stopOptions struct {
	config string
}

func init() {
	stopCmd.Flags().StringVarP(&stopOptions.config, "config", "c", Footloose, "Cluster configuration file")
	footloose.AddCommand(stopCmd)
}

func stop(cmd *cobra.Command, args []string) error {
	cluster, err := cluster.NewFromFile(configFile(stopOptions.config))
	if err != nil {
		return err
	}
	return cluster.Stop(args)
}
