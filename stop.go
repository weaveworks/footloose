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
	stopCmd.Flags().BoolVarP(&verbosity, "verbose", "v", false, "Verbosity commandline calls")
	stopCmd.Flags().StringVarP(&stopOptions.config, "config", "c", Footloose, "Cluster configuration file")
	footloose.AddCommand(stopCmd)
}

func stop(cmd *cobra.Command, args []string) error {
	cluster.GetCommanderInstance().SetVerbosity(verbosity)
	c, err := cluster.NewFromFile(stopOptions.config)
	if err != nil {
		return err
	}
	return c.Stop(args)
}
