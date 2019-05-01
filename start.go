package main

import (
	"github.com/spf13/cobra"

	"github.com/weaveworks/footloose/pkg/cluster"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start cluster machines",
	RunE:  start,
}

var startOptions struct {
	config string
}

func init() {
	startCmd.Flags().StringVarP(&startOptions.config, "config", "c", Footloose, "Cluster configuration file")
	footloose.AddCommand(startCmd)
}

func start(cmd *cobra.Command, args []string) error {
	cluster, err := cluster.NewFromFile(configFile(startOptions.config))
	if err != nil {
		return err
	}
	return cluster.Start(args)
}
