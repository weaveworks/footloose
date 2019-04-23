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
	startCmd.Flags().BoolVarP(&verbosity, "verbose", "v", false, "Verbosity commandline calls")
	startCmd.Flags().StringVarP(&startOptions.config, "config", "c", Footloose, "Cluster configuration file")
	footloose.AddCommand(startCmd)
}

func start(cmd *cobra.Command, args []string) error {
	cluster.GetCommanderInstance().SetVerbosity(verbosity)
	c, err := cluster.NewFromFile(startOptions.config)
	if err != nil {
		return err
	}
	return c.Start(args)
}
